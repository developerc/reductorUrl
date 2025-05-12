package service

import (
	"bytes"
	"context"
	"database/sql"
	"net"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/general"
	"github.com/gorilla/securecookie"
	"go.uber.org/zap"
)

type repository interface {
	AddLink(ctx context.Context, link string, usr string) (string, error)
	Ping() error
	GetLongLink(ctx context.Context, id string) (string, bool, error)
	HandleBatchJSON(ctx context.Context, buf bytes.Buffer, usr string) ([]byte, error)
	AsURLExists(err error) bool
	FetchURLs(ctx context.Context, cookieValue string) ([]byte, error)
	HandleCookie(cookieValue string) (*http.Cookie, string, error)
	DelURLs(cookieValue string, buf bytes.Buffer) error
	CloseDB() error
	GetStatsSvc(ctx context.Context, ip net.IP) ([]byte, error)
}

// Storager интерфейс для функций обработчиков
type Storager interface {
	AddLinkIface(ctx context.Context, link, usr string, s *Service) (string, error)
	PingIface(s *Service) error
	GetLongLinkIface(ctx context.Context, id string, s *Service) (string, bool, error)
	HandleBatchJSONIface(ctx context.Context, buf bytes.Buffer, usr string, s *Service) ([]byte, error)
	FetchURLsIface(ctx context.Context, cookieValue string, s *Service) ([]byte, error)
	DelURLsIface(cookieValue string, buf bytes.Buffer, s *Service) error
	CloseDBIface(s *Service) error
	GetStatsSvcIface(ctx context.Context, ip net.IP, s *Service) ([]byte, error)
}

// MapURLVal структура для значения map MapURL
type MapURLVal struct {
	OriginalURL string
	Usr         string
	IsDeleted   string
}

// ShortURLAttr структура аттрибутов коротких URL
type ShortURLAttr struct {
	MapURL   map[int64]MapURLVal
	MapUser  map[string]bool
	DB       *sql.DB
	Settings config.ServerSettings
	Cntr     int64
}

// Service структура сервиса приложения
type Service struct {
	//repo    repository
	Logger  *zap.Logger
	Secure  *securecookie.SecureCookie
	Shu     *ShortURLAttr
	Mu      sync.RWMutex
	Storage Storager
}

// User структура пользователя
type User struct {
	Name string
}

// AsURLExists делает проверку существования длинного URL
func (s *Service) AsURLExists(err error) bool {
	var errorURLExists general.ErrorURLExists
	return errorURLExists.AsURLExists(err)
}

// InitSecure создает обработчик куки
func (s *Service) InitSecure() {
	var hashKey = []byte("very-secret-qwer")
	var blockKey = []byte("a-lot-secret-qwe")
	s.Secure = securecookie.New(hashKey, blockKey)
}

// AddLink добавляет в хранилище длинный URL, возвращает короткий
func (s *Service) AddLink(ctx context.Context, link, usr string) (string, error) {
	//s.Shu.Cntr++
	atomic.AddInt64(&s.Shu.Cntr, 1)
	return s.Storage.AddLinkIface(ctx, link, usr, s)
}

// Ping проверяет живучесть БД
func (s *Service) Ping() error {
	return s.Storage.PingIface(s)
}

// GetLongLink получает длинный URL по ID
func (s *Service) GetLongLink(ctx context.Context, id string) (string, bool, error) {
	return s.Storage.GetLongLinkIface(ctx, id, s)
}

// HandleBatchJSON добавляет в хранилище несколько длинных URL
func (s *Service) HandleBatchJSON(ctx context.Context, buf bytes.Buffer, usr string) ([]byte, error) {
	return s.Storage.HandleBatchJSONIface(ctx, buf, usr, s)
}

// FetchURLs получает URL-ы определенного пользователя
func (s *Service) FetchURLs(ctx context.Context, cookieValue string) ([]byte, error) {
	return s.Storage.FetchURLsIface(ctx, cookieValue, s)
}

// HandleCookie метод для работы с куками
func (s *Service) HandleCookie(cookieValue string) (*http.Cookie, string, error) {
	var usr string
	var cookie *http.Cookie
	u := &User{
		Name: usr,
	}

	if cookieValue == "" {
		//usr = "user" + strconv.Itoa(s.Shu.Cntr)
		usr = "user" + strconv.FormatInt(s.Shu.Cntr, 10)
		u.Name = usr
		if encoded, err := s.Secure.Encode("user", u); err == nil {
			cookie = &http.Cookie{
				Name:  "user",
				Value: encoded,
			}
			return cookie, usr, nil
		} else {
			return nil, "", err
		}
	}
	if err := s.Secure.Decode("user", cookieValue, u); err != nil {
		return nil, "", err
	}
	s.Mu.RLock()
	_, ok := s.Shu.MapUser[u.Name]
	s.Mu.RUnlock()
	if ok {
		return nil, u.Name, nil
	} else {
		//usr = "user" + strconv.Itoa(s.Shu.Cntr)
		usr = "user" + strconv.FormatInt(s.Shu.Cntr, 10)
		u.Name = usr
		if encoded, err := s.Secure.Encode("user", u); err == nil {
			cookie = &http.Cookie{
				Name:  "user",
				Value: encoded,
			}
			s.Mu.Lock()
			s.Shu.MapUser[usr] = true
			s.Mu.Unlock()
			return cookie, usr, nil
		} else {
			return nil, "", err
		}
	}
}

// DelURLs делает отметку об удалении коротких URL-ы определенного пользователя
func (s *Service) DelURLs(cookieValue string, buf bytes.Buffer) error {
	return s.Storage.DelURLsIface(cookieValue, buf, s)
}

// CloseDB закрывает соединение с БД
func (s *Service) CloseDB() error {
	return s.Storage.CloseDBIface(s)
}

// GetStatsSvc получает статистику по количеству сокращённых URL в сервисе и количество пользователей в сервисе
func (s *Service) GetStatsSvc(ctx context.Context, ip net.IP) ([]byte, error) {
	return s.Storage.GetStatsSvcIface(ctx, ip, s)
}

// AddLongURL добавляет длинный URL в Map
func (s *Service) AddLongURL(i int64, link, usr string) {
	mapURLVal := MapURLVal{OriginalURL: link, Usr: usr, IsDeleted: "false"}
	s.Shu.MapURL[i] = mapURLVal
}
