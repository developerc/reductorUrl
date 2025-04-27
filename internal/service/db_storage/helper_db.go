package dbstorage

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/general"
	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/service"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

// ErrorURLExists структура типизированной ошибки существования длинного URL
/*type ErrorURLExists struct {
	s string
}

// Error возвращает строку со значением ошибки существования длинного URL
func (e *ErrorURLExists) Error() string {
	return e.s
}

// AsURLExists проверяет существование длинного URL
func (e *ErrorURLExists) AsURLExists(err error) bool {
	return errors.As(err, &e)
}*/

// StorageDB
type StorageDB struct {
}

// NewServiceDB конструктор сервиса
func NewServiceDB(ctx context.Context, settings *config.ServerSettings) (*service.Service, error) {
	var err error
	general.NewCntrAtom()

	shu := new(service.ShortURLAttr)
	shu.Settings = *settings
	shu.MapURL = make(map[int]service.MapURLVal)

	service := service.Service{Shu: shu}
	service.Logger, err = logger.Initialize(shu.Settings.LogLevel)
	service.InitSecure()

	InitDbStorage(ctx, &service)
	service.Storage = &StorageDB{}

	return &service, err
}

// AddLinkIface
func (sm *StorageDB) AddLinkIface(ctx context.Context, link, usr string, s *service.Service) (string, error) {
	shURL, err := InsertRecord(ctx, s.Shu.DB, link, usr)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == "must_be_different" {
			shortURL, err2 := getShortByOriginalURL(ctx, link, s)
			if err2 != nil {
				return "", err
			}
			return shortURL, &general.ErrorURLExists{Str: "this original URL exists"}
		}
		return "", err
	}
	s.Mu.Lock()
	s.Shu.MapUser[usr] = true
	s.Mu.Unlock()
	return s.Shu.Settings.AdresBase + "/" + shURL, nil
}

// getShortByOriginalURL
func getShortByOriginalURL(ctx context.Context, originalURL string, s *service.Service) (string, error) {
	shortURL, err := GetShortByOriginalURL(ctx, s.Shu.DB, originalURL)
	return s.Shu.Settings.AdresBase + "/" + shortURL, err
}

// InitDbStorage
func InitDbStorage(ctx context.Context, s *service.Service) error {
	var err error
	dsn := s.Shu.Settings.DBStorage
	s.Shu.DB, err = sql.Open("pgx", dsn)
	if err != nil {
		s.Logger.Info("InitDbStorage", zap.String("error", err.Error()))
		return err
	}
	if err = CreateTable(ctx, s.Shu.DB); err != nil {
		s.Logger.Info("InitDbStorage", zap.String("error", err.Error()))
	}
	s.Shu.MapUser, err = createMapUser(ctx, s)
	if err != nil {
		return err
	}
	return nil
}

// CreateMapUser создает Map пользователей
func createMapUser(ctx context.Context, s *service.Service) (map[string]bool, error) {
	mapUser, err := CreateMapUser(ctx, s.Shu.DB)
	if err != nil {
		return nil, err
	}
	s.Shu.Cntr = len(mapUser)
	return mapUser, nil
}

// GetLongLinkIface
func (sm *StorageDB) GetLongLinkIface(ctx context.Context, id string, s *service.Service) (string, bool, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return "", false, err
	}
	longURL, isDeleted, err := GetLongByUUID(ctx, s.Shu.DB, i)
	if err != nil {
		return "", false, err
	}
	return longURL, isDeleted, nil
}

// PingIface
func (sm *StorageDB) PingIface(s *service.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.Shu.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

// HandleBatchJSONIface
func (sm *StorageDB) HandleBatchJSONIface(ctx context.Context, buf bytes.Buffer, usr string, s *service.Service) ([]byte, error) {
	arrLongURL := make([]general.ArrLongURL, 0)
	if err := json.Unmarshal(buf.Bytes(), &arrLongURL); err != nil {
		return nil, err
	}
	if len(arrLongURL) == 0 {
		return nil, errors.New("error: length array is zero")
	}
	if err := InsertBatch2(ctx, arrLongURL, s.Shu.DB, usr); err != nil {
		return nil, err
	}

	arrShortURL := make([]general.ArrShortURL, 0)
	for _, longURL := range arrLongURL {
		short, err := GetShortByOriginalURL(ctx, s.Shu.DB, longURL.OriginalURL)
		if err != nil {
			return nil, err
		}
		shortURL := general.ArrShortURL{CorellationID: longURL.CorellationID, ShortURL: s.Shu.Settings.AdresBase + "/" + short}
		arrShortURL = append(arrShortURL, shortURL)
	}
	jsonBytes, err := json.Marshal(arrShortURL)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// FetchURLsIface
func (sm *StorageDB) FetchURLsIface(ctx context.Context, cookieValue string, s *service.Service) ([]byte, error) {
	u := &general.User{}
	if err := s.Secure.Decode("user", cookieValue, u); err != nil {
		return nil, err
	}

	s.Mu.RLock()
	_, ok := s.Shu.MapUser[u.Name]
	s.Mu.RUnlock()
	if !ok {
		return nil, http.ErrNoCookie
	}
	arrRepoURL, err := ListRepoURLs(ctx, s.Shu.DB, s.Shu.Settings.AdresBase, u.Name)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(arrRepoURL)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// DelURLsIface
func (sm *StorageDB) DelURLsIface(cookieValue string, buf bytes.Buffer, s *service.Service) error {
	general.CntrAtomVar.IncrCntr()
	u := &general.User{}
	if err := s.Secure.Decode("user", cookieValue, u); err != nil {
		return err
	}

	s.Mu.RLock()
	_, ok := s.Shu.MapUser[u.Name]
	s.Mu.RUnlock()
	if !ok {
		return http.ErrNoCookie
	}

	arrShortURL := make([]string, 0)
	if err := json.Unmarshal(buf.Bytes(), &arrShortURL); err != nil {
		return err
	}

	if err := SetDelBatch(arrShortURL, s.Shu.DB, u.Name); err != nil {
		return err
	}

	general.CntrAtomVar.DecrCntr()
	general.CntrAtomVar.SentNotif()
	return nil
}

// CloseDBIface
func (sm *StorageDB) CloseDBIface(s *service.Service) error {
	return s.Shu.DB.Close()
}

// GetStatsSvcIface
func (sm *StorageDB) GetStatsSvcIface(ctx context.Context, ip net.IP, s *service.Service) ([]byte, error) {
	if s.Shu.Settings.TrustedSubnet == "" {
		return nil, nil
	}

	_, ipNet, err := net.ParseCIDR(s.Shu.Settings.TrustedSubnet)
	if err != nil {
		return nil, err
	}
	ipCheck := net.ParseIP(ip.String())
	if ipNet.Contains(ipCheck) {
		return GetStatsDB(ctx, s.Shu.DB)
	}

	return nil, nil
}
