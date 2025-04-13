// memory пакет для размещения сервисных методов приложения.
package memory

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/securecookie"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/logger"
	dbstorage "github.com/developerc/reductorUrl/internal/service/db_storage"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
)

type repository interface {
	AddLink(link string, usr string) (string, error)
	Ping() error
	GetLongLink(id string) (string, bool, error)
	HandleBatchJSON(buf bytes.Buffer, usr string) ([]byte, error)
	AsURLExists(err error) bool
	GetShu() *ShortURLAttr
	FetchURLs(cookieValue string) ([]byte, error)
	HandleCookie(cookieValue string) (*http.Cookie, string, error)
	DelURLs(cookieValue string, buf bytes.Buffer) (bool, error)
	CloseDb() error
}

// Service структура сервиса приложения
type Service struct {
	repo   repository
	logger *zap.Logger
	secure *securecookie.SecureCookie
	//chDbClose chan struct{}
}

// AsURLExists делает проверку существования длинного URL
func (s *Service) AsURLExists(err error) bool {
	var errorURLExists ErrorURLExists
	return errorURLExists.AsURLExists(err)
}

// ErrorURLExists структура типизированной ошибки существования длинного URL
type ErrorURLExists struct {
	s string
}

// Error возвращает строку со значением ошибки существования длинного URL
func (e *ErrorURLExists) Error() string {
	return e.s
}

// AsURLExists проверяет существование длинного URL
func (e *ErrorURLExists) AsURLExists(err error) bool {
	return errors.As(err, &e)
}

// AddLink добавляет в хранилище длинный URL, возвращает короткий
func (s *Service) AddLink(link, usr string) (string, error) {
	var shURL string
	var err error

	s.IncrCounter()
	switch s.repo.GetShu().Settings.TypeStorage {
	case config.MemoryStorage:
		s.AddLongURL(s.GetCounter(), link)
		return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
	case config.FileStorage:
		if err = s.repo.GetShu().addToFileStorage(s.GetCounter(), link); err != nil {
			return "", err
		}
		s.AddLongURL(s.GetCounter(), link)
		return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
	case config.DBStorage:
		shURL, err = dbstorage.InsertRecord(s.repo.GetShu().DB, link, usr)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == "must_be_different" {
				shortURL, err2 := s.GetShortByOriginalURL(link)
				if err2 != nil {
					return "", err
				}
				return shortURL, &ErrorURLExists{"this original URL exists"}
			}
			return "", err
		}
		s.repo.GetShu().MapUser[usr] = true
	}
	return s.GetAdresBase() + "/" + shURL, nil
}

// GetShortByOriginalURL получает короткий URL по значению длинного
func (s *Service) GetShortByOriginalURL(originalURL string) (string, error) {
	shortURL, err := dbstorage.GetShortByOriginalURL(s.repo.GetShu().DB, originalURL)
	return s.GetAdresBase() + "/" + shortURL, err
}

// GetLongLink получает длинный URL по ID
func (s *Service) GetLongLink(id string) (longURL string, isDeleted bool, err error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	switch s.repo.GetShu().Settings.TypeStorage {
	case config.MemoryStorage:
		longURL, err = s.GetLongURL(i)
		if err != nil {
			return
		}
	case config.FileStorage:
		longURL, err = s.GetLongURL(i)
		if err != nil {
			return
		}
	case config.DBStorage:
		longURL, isDeleted, err = dbstorage.GetLongByUUID(s.repo.GetShu().DB, i)
		if err != nil {
			return
		}
	}
	return
}

// HandleBatchJSON добавляет в хранилище несколько длинных URL
func (s *Service) HandleBatchJSON(buf bytes.Buffer, usr string) ([]byte, error) {
	arrLongURL, err := listLongURL(buf)
	if err != nil {
		return nil, err
	}
	if len(arrLongURL) == 0 {
		return nil, errors.New("error: length array is zero")
	}

	jsonBytes, err := s.handleArrLongURL(arrLongURL, usr)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (s *Service) CloseDb() error {
	//fmt.Println("CloseChan")
	//close(s.chDbClose)
	//s.repo.GetShu().DB.
	if err := s.repo.GetShu().DB.Close(); err != nil {
		//s.logger.Info("Close DB", zap.String("error", err.Error()))
		return err
	}
	return nil
}

// NewInMemoryService конструктор сервиса
func NewInMemoryService() (*Service, error) {
	var err error

	shu := new(ShortURLAttr)
	shu.Settings = *config.NewServerSettings()
	shu.MapURL = make(map[int]string)

	switch shu.Settings.TypeStorage {
	case config.FileStorage:
		if err = getFileSettings(shu); err != nil {
			log.Println(err)
		}
	case config.DBStorage:
		dsn := shu.Settings.DBStorage
		shu.DB, err = sql.Open("pgx", dsn)
		if err != nil {
			return nil, err
		}
		if err = dbstorage.CreateTable(shu.DB); err != nil {
			log.Println(err)
		}
		shu.MapUser, err = CreateMapUser(shu)
		if err != nil {
			return nil, err
		}

	}

	service := Service{repo: shu}
	service.logger, err = logger.Initialize(service.GetLogLevel())
	service.InitSecure()
	/*go func() {
		service.chDbClose = make(chan struct{})
		<-service.chDbClose
		fmt.Println("service.chDbClose closed")
		//shu.DB.Close()
	}()*/
	return &service, err
}

// InitSecure создает обработчик куки
func (s *Service) InitSecure() {
	var hashKey = []byte("very-secret-qwer")
	var blockKey = []byte("a-lot-secret-qwe")
	s.secure = securecookie.New(hashKey, blockKey)
}

// AddLink заглушка для ShortURLAttr
func (shu *ShortURLAttr) AddLink(link, usr string) (string, error) {
	return "", nil
}

// addToFileStorage добавляет длинный URL в файловое хранилище
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

// Ping заглушка для ShortURLAttr
func (shu *ShortURLAttr) Ping() error {
	return nil
}

// GetLongLink заглушка для ShortURLAttr
func (shu *ShortURLAttr) GetLongLink(id string) (longURL string, isDeleted bool, err error) {
	return "", false, nil
}

// HandleBatchJSON заглушка для ShortURLAttr
func (shu *ShortURLAttr) HandleBatchJSON(buf bytes.Buffer, usr string) ([]byte, error) {
	return nil, nil
}

// AsURLExists заглушка для ShortURLAttr
func (shu *ShortURLAttr) AsURLExists(err error) bool {
	return true
}

// GetShu заглушка для ShortURLAttr
func (shu *ShortURLAttr) GetShu() *ShortURLAttr {
	return shu
}

// GetCripto заглушка для ShortURLAttr
func (shu *ShortURLAttr) GetCripto() (string, error) {
	return "", nil
}

// FetchURLs заглушка для ShortURLAttr
func (shu *ShortURLAttr) FetchURLs(cookieValue string) ([]byte, error) {
	return nil, nil
}

// GetCounter заглушка для ShortURLAttr
func (shu *ShortURLAttr) GetCounter() int {
	return 0
}

// HandleCookie заглушка для ShortURLAttr
func (shu *ShortURLAttr) HandleCookie(cookieValue string) (*http.Cookie, string, error) {
	return nil, "", nil
}

// DelURLs заглушка для ShortURLAttr
func (shu *ShortURLAttr) DelURLs(cookieValue string, buf bytes.Buffer) (bool, error) {
	return false, nil
}

// CloseDb заглушка для ShortURLAttr
func (shu *ShortURLAttr) CloseDb() error {
	return nil
}
