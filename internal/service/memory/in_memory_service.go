package memory

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/logger"
	dbstorage "github.com/developerc/reductorUrl/internal/service/db_storage"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type repository interface {
	AddLink(link string) (string, error)
	Ping() error
	GetLongLink(id string) (string, error)
	HandleBatchJSON(buf bytes.Buffer) ([]byte, error)
	AsURLExists(err error) bool
	GetShu() *ShortURLAttr
	GetCripto() (string, error)
	FetchURLs() ([]byte, error)
	IsRegisteredCookie(cookie string) bool
}

type Service struct {
	repo   repository
	logger *zap.Logger
}

// IsRegisteredCookie implements server.svc.
func (s *Service) IsRegisteredCookie(cookie string) bool {
	return false
}

// GetCripto implements server.svc.
/*func (s *Service) GetCripto() ([]byte, error) {
	panic("unimplemented")
}*/

func (s *Service) AsURLExists(err error) bool {
	var errorURLExists ErrorURLExists
	return errorURLExists.AsURLExists(err)
}

type ErrorURLExists struct {
	s string
}

func (e *ErrorURLExists) Error() string {
	return e.s
}

func (e *ErrorURLExists) AsURLExists(err error) bool {
	return errors.As(err, &e)
}

func (s *Service) AddLink(link string) (string, error) {
	var shURL string
	var err error

	s.IncrCounter()
	switch s.GetShortURLAttr().Settings.TypeStorage {
	case config.MemoryStorage:
		s.AddLongURL(s.GetCounter(), link)
		return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
	case config.FileStorage:
		if err := s.GetShortURLAttr().addToFileStorage(s.GetCounter(), link); err != nil {
			return "", err
		}
		s.AddLongURL(s.GetCounter(), link)
		return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
	case config.DBStorage:
		shURL, err = dbstorage.InsertRecord(s.GetShortURLAttr().DB, link)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == "must_be_different" {
				shortURL, err := s.GetShortByOriginalURL(link)
				if err != nil {
					return "", err
				}
				return shortURL, &ErrorURLExists{"this original URL exists"}
			}
			return "", err
		}
	}
	return s.GetAdresBase() + "/" + shURL, nil
}

func (s *Service) GetShortByOriginalURL(originalURL string) (string, error) {
	shortURL, err := dbstorage.GetShortByOriginalURL(s.GetShortURLAttr().DB, originalURL)
	return s.GetAdresBase() + "/" + shortURL, err
}

func (s *Service) GetLongLink(id string) (string, error) {
	var longURL string
	i, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	switch s.GetShortURLAttr().Settings.TypeStorage {
	case config.MemoryStorage:
		longURL, err = s.GetLongURL(i)
		if err != nil {
			return "", err
		}
	case config.FileStorage:
		longURL, err = s.GetLongURL(i)
		if err != nil {
			return "", err
		}
	case config.DBStorage:
		longURL, err = dbstorage.GetLongByUUID(s.GetShortURLAttr().DB, i)
		if err != nil {
			return "", err
		}
	}
	return longURL, nil
}

func (s *Service) HandleBatchJSON(buf bytes.Buffer) ([]byte, error) {
	arrLongURL, err := listLongURL(buf)
	if err != nil {
		return nil, err
	}
	if len(arrLongURL) == 0 {
		return nil, errors.New("error: length array is zero")
	}

	jsonBytes, err := s.handleArrLongURL(arrLongURL)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func NewInMemoryService() (*Service, error) {
	var err error

	shu := new(ShortURLAttr)
	shu.Settings = *config.NewServerSettings()
	shu.MapURL = make(map[int]string)

	switch shu.Settings.TypeStorage {
	case config.FileStorage:
		if err := getFileSettings(shu); err != nil {
			log.Println(err)
		}
	case config.DBStorage:
		dsn := shu.Settings.DBStorage
		shu.DB, err = sql.Open("pgx", dsn)
		if err != nil {
			return nil, err
		}
		if err := dbstorage.CreateTable(shu.DB); err != nil {
			log.Println(err)
		}
		shu.MapCookie, err = CreateMapCookie(shu)
		if err != nil {
			return nil, err
		}
	}

	service := Service{repo: shu}
	service.logger, err = logger.Initialize(service.GetLogLevel())
	return &service, err
}

func (shu *ShortURLAttr) AddLink(link string) (string, error) {
	return "", nil
}

func (shu *ShortURLAttr) addToFileStorage(cntr int, link string) error {
	if cntr < 0 {
		return errors.New("not valid counter")
	}
	event := filestorage.Event{UUID: uint(cntr), ShortURL: strconv.Itoa(cntr), OriginalURL: link}
	producer, err := filestorage.NewProducer(shu.Settings.FileStorage)
	if err != nil {
		return err
	}
	if err := producer.WriteEvent(&event); err != nil {
		log.Println(err)
	}
	return nil
}

func (shu *ShortURLAttr) IsRegisteredCookie(cookie string) bool {
	if _, ok := shu.MapCookie[cookie]; ok {
		return true
	}
	return false
}

func (shu *ShortURLAttr) Ping() error {
	return nil
}

func (shu *ShortURLAttr) GetLongLink(id string) (string, error) {
	return "", nil
}

func (shu *ShortURLAttr) HandleBatchJSON(buf bytes.Buffer) ([]byte, error) {
	return nil, nil
}

func (shu *ShortURLAttr) AsURLExists(err error) bool {
	return true
}

func (shu *ShortURLAttr) GetShu() *ShortURLAttr {
	return shu
}

func (shu *ShortURLAttr) GetCripto() (string, error) {
	return "", nil
}

func (shu *ShortURLAttr) FetchURLs() ([]byte, error) {
	return nil, nil
}
