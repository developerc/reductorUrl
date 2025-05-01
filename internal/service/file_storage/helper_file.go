package filestorage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/general"
	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/service"
	"go.uber.org/zap"
)

// StorageFile структура для слоя File Storage
type StorageFile struct {
}

// NewServiceMemory конструктор сервиса
func NewServiceFile(ctx context.Context, settings *config.ServerSettings) (*service.Service, error) {
	var err error
	general.NewCntrAtom()

	shu := new(service.ShortURLAttr)
	shu.Settings = *settings
	shu.MapURL = make(map[int]service.MapURLVal)

	service := service.Service{Shu: shu}
	service.Logger, err = logger.Initialize(shu.Settings.LogLevel)
	service.InitSecure()

	InitFileStorage(&service)
	service.Storage = &StorageFile{}
	return &service, err
}

// AddLinkIface добавляет в хранилище длинный URL, возвращает короткий
func (sm *StorageFile) AddLinkIface(ctx context.Context, link, usr string, s *service.Service) (string, error) {
	s.Mu.Lock()
	if err := addToFileStorage(s.Shu.Cntr, link, usr, s); err != nil {
		return "", err
	}
	s.Shu.MapUser[usr] = true
	s.Mu.Unlock()
	s.AddLongURL(s.Shu.Cntr, link, usr)
	return s.Shu.Settings.AdresBase + "/" + strconv.Itoa(s.Shu.Cntr), nil
}

// addToFileStorage добавляет длинный URL в файловое хранилище
func addToFileStorage(cntr int, link, usr string, s *service.Service) error {
	if cntr < 0 {
		return errors.New("not valid counter")
	}
	event := Event{UUID: uint(cntr), OriginalURL: link, Usr: usr, IsDeleted: "false"}
	producer, err := NewProducer(s.Shu.Settings.FileStorage)
	if err != nil {
		return err
	}
	if err := producer.WriteEvent(&event); err != nil {
		//log.Println(err)
		s.Logger.Info("addToFileStorage", zap.String("error", err.Error()))
	}
	return nil
}

// InitFileStorage инициализирует файловое хранилище
func InitFileStorage(s *service.Service) {
	s.Shu.MapUser = make(map[string]bool)
	if err := getFileSettings(s); err != nil {
		s.Logger.Info("addToFileStorage", zap.String("error", err.Error()))
	}
}

// getFileSettings получает настройки из файлового хранилища
func getFileSettings(s *service.Service) error {
	if _, err := NewConsumer(s.Shu.Settings.FileStorage); err != nil {
		return err
	}
	consumer, err := NewConsumer(s.Shu.Settings.FileStorage)
	if err != nil {
		return err
	}
	events, err := consumer.ListEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		if event.UUID > math.MaxInt32 {
			event.UUID = math.MaxInt32
		}

		s.Shu.MapURL[int(event.UUID)] = service.MapURLVal{OriginalURL: event.OriginalURL, Usr: event.Usr, IsDeleted: event.IsDeleted}
		s.Shu.MapUser[s.Shu.MapURL[int(event.UUID)].Usr] = true
	}
	s.Shu.Cntr = len(events)

	if _, err := NewProducer(s.Shu.Settings.FileStorage); err != nil {
		return err
	}
	return nil
}

// GetLongLinkIface получает длинный URL по ID
func (sm *StorageFile) GetLongLinkIface(ctx context.Context, id string, s *service.Service) (string, bool, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		return "", false, err
	}
	return s.Shu.MapURL[i].OriginalURL, false, nil
}

// PingIface заглушка для пинг БД
func (sm *StorageFile) PingIface(s *service.Service) error {
	return nil
}

// HandleBatchJSONIface добавляет в хранилище несколько длинных URL
func (sm *StorageFile) HandleBatchJSONIface(ctx context.Context, buf bytes.Buffer, usr string, s *service.Service) ([]byte, error) {
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
func (sm *StorageFile) FetchURLsIface(ctx context.Context, cookieValue string, s *service.Service) ([]byte, error) {
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
func (sm *StorageFile) DelURLsIface(cookieValue string, buf bytes.Buffer, s *service.Service) error {
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
	changeFileStorage(s)

	general.CntrAtomVar.DecrCntr()
	general.CntrAtomVar.SentNotif()
	return nil
}

// CloseDBIface закрывает соединение с БД
func (sm *StorageFile) CloseDBIface(s *service.Service) error {
	return nil
}

// GetStatsSvcIface получает статистику по количеству сокращённых URL в сервисе и количество пользователей в сервисе
func (sm *StorageFile) GetStatsSvcIface(ctx context.Context, ip net.IP, s *service.Service) ([]byte, error) {
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

// changeFileStorage заменяет файл файлового хранилища на обновленный
func changeFileStorage(s *service.Service) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	os.Remove(s.Shu.Settings.FileStorage)
	for uuid, mapURL := range s.Shu.MapURL {
		event := Event{UUID: uint(uuid), OriginalURL: mapURL.OriginalURL, Usr: mapURL.Usr, IsDeleted: mapURL.IsDeleted}
		producer, err := NewProducer(s.Shu.Settings.FileStorage)
		if err != nil {
			return err
		}
		if err := producer.WriteEvent(&event); err != nil {
			log.Println(err)
		}
	}
	return nil
}
