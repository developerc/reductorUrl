// memory пакет для размещения сервисных методов приложения.
package memory

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/securecookie"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/general"
	"github.com/developerc/reductorUrl/internal/logger"
	dbstorage "github.com/developerc/reductorUrl/internal/service/db_storage"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
)

type repository interface {
	AddLink(ctx context.Context, link string, usr string) (string, error)
	Ping() error
	GetLongLink(ctx context.Context, id string) (string, bool, error)
	HandleBatchJSON(ctx context.Context, buf bytes.Buffer, usr string) ([]byte, error)
	AsURLExists(err error) bool
	FetchURLs(ctx context.Context, cookieValue string) ([]byte, error)
	HandleCookie(cookieValue string) (*http.Cookie, string, error)
	DelURLs(ctx context.Context, cookieValue string, buf bytes.Buffer) error
	CloseDB() error
	GetStatsSvc(ctx context.Context, ip net.IP) ([]byte, error)
}

// Service структура сервиса приложения
type Service struct {
	repo   repository
	logger *zap.Logger
	secure *securecookie.SecureCookie
	shu    *ShortURLAttr
	mu     sync.RWMutex
}

// AsURLExists делает проверку существования длинного URL
func (s *Service) AsURLExists(err error) bool {
	var errorURLExists ErrorURLExists
	//fmt.Println(errorURLExists.AsURLExists(err))
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
func (s *Service) AddLink(ctx context.Context, link, usr string) (string, error) {
	var shURL string
	var err error

	s.IncrCounter()
	switch s.shu.Settings.TypeStorage {
	case config.MemoryStorage:
		s.AddLongURL(s.GetCounter(), link, usr)
		s.mu.Lock()
		s.shu.MapUser[usr] = true
		s.mu.Unlock()
		return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
	case config.FileStorage:
		s.mu.Lock()
		if err = s.shu.addToFileStorage(s.GetCounter(), link, usr); err != nil {
			return "", err
		}
		s.shu.MapUser[usr] = true
		s.mu.Unlock()
		s.AddLongURL(s.GetCounter(), link, usr)
		return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
	case config.DBStorage:
		shURL, err = dbstorage.InsertRecord(ctx, s.shu.DB, link, usr)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == "must_be_different" {
				shortURL, err2 := s.GetShortByOriginalURL(ctx, link)
				if err2 != nil {
					return "", err
				}
				return shortURL, &ErrorURLExists{"this original URL exists"}
			}
			return "", err
		}
		s.mu.Lock()
		s.shu.MapUser[usr] = true
		s.mu.Unlock()
	}
	return s.GetAdresBase() + "/" + shURL, nil
}

// GetShortByOriginalURL получает короткий URL по значению длинного
func (s *Service) GetShortByOriginalURL(ctx context.Context, originalURL string) (string, error) {
	shortURL, err := dbstorage.GetShortByOriginalURL(ctx, s.shu.DB, originalURL)
	return s.GetAdresBase() + "/" + shortURL, err
}

// GetLongLink получает длинный URL по ID
func (s *Service) GetLongLink(ctx context.Context, id string) (longURL string, isDeleted bool, err error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	switch s.shu.Settings.TypeStorage {
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
		longURL, isDeleted, err = dbstorage.GetLongByUUID(ctx, s.shu.DB, i)
		if err != nil {
			return
		}
	}
	return
}

// HandleBatchJSON добавляет в хранилище несколько длинных URL
func (s *Service) HandleBatchJSON(ctx context.Context, buf bytes.Buffer, usr string) ([]byte, error) {
	arrLongURL, err := listLongURL(buf)
	if err != nil {
		return nil, err
	}
	if len(arrLongURL) == 0 {
		return nil, errors.New("error: length array is zero")
	}

	jsonBytes, err := s.handleArrLongURL(ctx, arrLongURL, usr)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}

// CloseDB закрывает соединение с БД
func (s *Service) CloseDB() error {
	if s.shu.Settings.TypeStorage != config.DBStorage {
		return nil
	}

	return s.shu.DB.Close()
}

// NewInMemoryService конструктор сервиса
func NewInMemoryService(ctx context.Context) (*Service, error) {
	var err error
	general.NewCntrAtom()

	shu := new(ShortURLAttr)
	shu.Settings = *config.NewServerSettings()
	shu.MapURL = make(map[int]MapURLVal)

	switch shu.Settings.TypeStorage {
	case config.MemoryStorage:
		shu.MapUser = make(map[string]bool)
	case config.FileStorage:
		shu.MapUser = make(map[string]bool)
		if err = getFileSettings(shu); err != nil {
			log.Println(err)
		}
	case config.DBStorage:
		dsn := shu.Settings.DBStorage
		shu.DB, err = sql.Open("pgx", dsn)
		if err != nil {
			return nil, err
		}
		if err = dbstorage.CreateTable(ctx, shu.DB); err != nil {
			log.Println(err)
		}
		shu.MapUser, err = CreateMapUser(ctx, shu)
		if err != nil {
			return nil, err
		}

	}

	service := Service{shu: shu}
	service.logger, err = logger.Initialize(service.GetLogLevel())
	service.InitSecure()
	return &service, err
}

// InitSecure создает обработчик куки
func (s *Service) InitSecure() {
	var hashKey = []byte("very-secret-qwer")
	var blockKey = []byte("a-lot-secret-qwe")
	s.secure = securecookie.New(hashKey, blockKey)
}

// addToFileStorage добавляет длинный URL в файловое хранилище
func (shu *ShortURLAttr) addToFileStorage(cntr int, link, usr string) error {
	if cntr < 0 {
		return errors.New("not valid counter")
	}
	event := filestorage.Event{UUID: uint(cntr), OriginalURL: link, Usr: usr, IsDeleted: "false"}
	producer, err := filestorage.NewProducer(shu.Settings.FileStorage)
	if err != nil {
		return err
	}
	if err := producer.WriteEvent(&event); err != nil {
		log.Println(err)
	}
	return nil
}

// func (shu *ShortURLAttr) changeFileStorage() error {
func (s *Service) changeFileStorage() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	os.Remove(s.shu.Settings.FileStorage)
	for uuid, mapURL := range s.shu.MapURL {
		event := filestorage.Event{UUID: uint(uuid), OriginalURL: mapURL.OriginalURL, Usr: mapURL.Usr, IsDeleted: mapURL.IsDeleted}
		producer, err := filestorage.NewProducer(s.shu.Settings.FileStorage)
		if err != nil {
			return err
		}
		if err := producer.WriteEvent(&event); err != nil {
			log.Println(err)
		}
	}
	return nil
}

// GetStatsSvc получает статистику по количеству сокращённых URL в сервисе и количество пользователей в сервисе
func (s *Service) GetStatsSvc(ctx context.Context, ip net.IP) ([]byte, error) {
	//fmt.Println("ip: ", ip)
	//fmt.Println(s.shu.Settings.TrustedSubnet)
	if s.shu.Settings.TrustedSubnet == "" {
		return nil, nil
	}
	var users int
	urls := len(s.shu.MapURL)
	_, ipNet, err := net.ParseCIDR(s.shu.Settings.TrustedSubnet)
	if err != nil {
		return nil, err
	}
	ipCheck := net.ParseIP(ip.String())
	if ipNet.Contains(ipCheck) {
		//return []byte("proba"), nil
		var jsonBytes []byte
		if s.shu.Settings.TypeStorage != config.DBStorage {
			usersMap := make(map[string]bool)

			for _, val := range s.shu.MapURL {
				//fmt.Println(val)
				usersMap[val.Usr] = true
			}
			//fmt.Println(usersMap)
			users = len(usersMap)
			arrGetStats := general.ArrGetStats{URLs: urls, Users: users}
			jsonBytes, err = json.Marshal(arrGetStats)
			if err != nil {
				return nil, err
			}
			return jsonBytes, nil
		} else {
			return dbstorage.GetStatsDB(ctx, s.shu.DB)
		}
	}
	return nil, nil
}
