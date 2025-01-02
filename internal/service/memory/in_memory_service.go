package memory

import (
	"bytes"
	"errors"
	"log"
	"strconv"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/logger"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
	"go.uber.org/zap"
)

type repository interface {
	AddLink(link string) (string, error)
}

type Service struct {
	repo   repository
	logger *zap.Logger
}

func (s *Service) AddLink(link string) (string, error) {
	var shURL string
	var err error
	const memoryStorage string = "MemoryStorage"
	const fileStorage string = "FileStorage"
	const dbStorage string = "DBStorage"

	s.IncrCounter()
	switch s.GetShortURLAttr().Settings.TypeStorage {
	case memoryStorage:
		s.AddLongURL(s.GetCounter(), link)
		return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
	case fileStorage:
		if err := s.GetShortURLAttr().addToFileStorage(s.GetCounter(), link); err != nil {
			return "", err
		}
		s.AddLongURL(s.GetCounter(), link)
		return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
	case dbStorage:
		shURL, err = insertRecord(s.GetShortURLAttr(), link)
		if err != nil {
			return "", err
		}
	}
	return s.GetAdresBase() + "/" + shURL, nil
}

func (s *Service) GetShortByOriginalURL(originalURL string) (string, error) {
	shortURL, err := getShortByOriginalURL(s.GetShortURLAttr(), originalURL)
	return s.GetAdresBase() + "/" + shortURL, err
}

func (s *Service) GetLongLink(id string) (string, error) {
	var longURL string
	i, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	switch s.GetShortURLAttr().Settings.TypeStorage {
	case "MemoryStorage":
		longURL, err = s.GetLongURL(i)
		if err != nil {
			return "", err
		}
	case "FileStorage":
		longURL, err = s.GetLongURL(i)
		if err != nil {
			return "", err
		}
	case "DBStorage":
		longURL, err = getLongByUUID(s.GetShortURLAttr(), i)
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
	case "FileStorage":
		if err := getFileSettings(shu); err != nil {
			log.Println(err)
		}
	case "DBStorage":
		if err := createTable(shu); err != nil {
			log.Println(err)
		}
	}

	service := Service{repo: shu}
	service.logger, err = logger.Initialize(service.GetLogLevel())
	return &service, err
}

func (shu *ShortURLAttr) AddLink(link string) (string, error) {
	return "proba", nil
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
