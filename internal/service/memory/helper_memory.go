package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/general"
	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/service"
)

// StorageMem структура для слоя Memory Storage
type StorageMem struct {
}

// NewServiceMemory конструктор сервиса
func NewServiceMemory(ctx context.Context, settings *config.ServerSettings) (*service.Service, error) {
	var err error
	general.NewCntrAtom()

	shu := new(service.ShortURLAttr)
	shu.Settings = *settings
	shu.MapURL = make(map[int]service.MapURLVal)

	service := service.Service{Shu: shu}
	service.Logger, err = logger.Initialize(shu.Settings.LogLevel)
	service.InitSecure()

	service.Shu.MapUser = make(map[string]bool)
	service.Storage = &StorageMem{}
	return &service, err
}

// AddLinkIface добавляет в хранилище длинный URL, возвращает короткий
func (sm *StorageMem) AddLinkIface(ctx context.Context, link, usr string, s *service.Service) (string, error) {
	s.AddLongURL(s.Shu.Cntr, link, usr)
	s.Mu.Lock()
	s.Shu.MapUser[usr] = true
	s.Mu.Unlock()
	return s.Shu.Settings.AdresBase + "/" + strconv.Itoa(s.Shu.Cntr), nil
}

// GetLongLinkIface получает длинный URL по ID
func (sm *StorageMem) GetLongLinkIface(ctx context.Context, id string, s *service.Service) (string, bool, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return "", false, err
	}
	return s.Shu.MapURL[i].OriginalURL, false, nil
}

// PingIface заглушка для пинга БД
func (sm *StorageMem) PingIface(s *service.Service) error {
	return nil
}

// HandleBatchJSONIface добавляет в хранилище несколько длинных URL
func (sm *StorageMem) HandleBatchJSONIface(ctx context.Context, buf bytes.Buffer, usr string, s *service.Service) ([]byte, error) {
	arrLongURL := make([]general.ArrLongURL, 0)
	if err := json.Unmarshal(buf.Bytes(), &arrLongURL); err != nil {
		return nil, err
	}
	if len(arrLongURL) == 0 {
		return nil, errors.New("error: length array is zero")
	}

	arrShortURL := make([]general.ArrShortURL, 0)
	for _, longURL := range arrLongURL {
		URL, err := s.AddLink(ctx, longURL.OriginalURL, usr)
		if err != nil {
			return nil, err
		}
		shortURL := general.ArrShortURL{CorellationID: longURL.CorellationID, ShortURL: URL}
		arrShortURL = append(arrShortURL, shortURL)
	}
	jsonBytes, err := json.Marshal(arrShortURL)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// FetchURLsIface получает URL-ы определенного пользователя
func (sm *StorageMem) FetchURLsIface(ctx context.Context, cookieValue string, s *service.Service) ([]byte, error) {
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
	var err error
	arrRepoURL := make([]general.ArrRepoURL, 0)
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for uuid, val := range s.Shu.MapURL {
		if val.Usr == u.Name {
			repoURL := general.ArrRepoURL{}
			repoURL.ShortURL = s.Shu.Settings.AdresBase + "/" + strconv.Itoa(uuid)
			repoURL.OriginalURL = val.OriginalURL
			arrRepoURL = append(arrRepoURL, repoURL)
		}
	}

	jsonBytes, err := json.Marshal(arrRepoURL)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// DelURLsIface делает отметку об удалении коротких URL-ы определенного пользователя
func (sm *StorageMem) DelURLsIface(cookieValue string, buf bytes.Buffer, s *service.Service) error {
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

	for _, shortURL := range arrShortURL {
		intShortURL, err := strconv.Atoi(shortURL)
		if err != nil {
			return err
		}
		if val, ok := s.Shu.MapURL[intShortURL]; ok {
			if val.Usr == u.Name {
				s.Mu.Lock()
				mapURLVal := s.Shu.MapURL[intShortURL]
				mapURLVal.IsDeleted = "true"
				s.Shu.MapURL[intShortURL] = mapURLVal
				s.Mu.Unlock()
			}
		}
	}

	general.CntrAtomVar.DecrCntr()
	general.CntrAtomVar.SentNotif()
	return nil
}

// CloseDBIface закрывает соединение с БД
func (sm *StorageMem) CloseDBIface(s *service.Service) error {
	return nil
}

// GetStatsSvcIface получает статистику по количеству сокращённых URL в сервисе и количество пользователей в сервисе
func (sm *StorageMem) GetStatsSvcIface(ctx context.Context, ip net.IP, s *service.Service) ([]byte, error) {
	if s.Shu.Settings.TrustedSubnet == "" {
		return nil, nil
	}
	var users int
	urls := len(s.Shu.MapURL)
	_, ipNet, err := net.ParseCIDR(s.Shu.Settings.TrustedSubnet)
	if err != nil {
		return nil, err
	}
	ipCheck := net.ParseIP(ip.String())
	if ipNet.Contains(ipCheck) {
		usersMap := make(map[string]bool)
		for _, val := range s.Shu.MapURL {
			usersMap[val.Usr] = true
		}
		users = len(usersMap)
		arrGetStats := general.ArrGetStats{URLs: urls, Users: users}
		jsonBytes, err := json.Marshal(arrGetStats)
		if err != nil {
			return nil, err
		}
		return jsonBytes, nil
	}

	return nil, nil
}
