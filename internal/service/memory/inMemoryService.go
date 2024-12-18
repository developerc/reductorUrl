package memory

import (
	"errors"
	"math"
	"strconv"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/service/filestorage"
	"go.uber.org/zap"
)

type repository interface {
	AddLink(link string) (string, error)
}

type Service struct {
	repo repository
}

type ShortURLAttr struct {
	Settings config.ServerSettings
	Cntr     int
	MapURL   map[int]string
}

func (s Service) AddLink(link string) (string, error) {
	s.repo.(*ShortURLAttr).Cntr++
	if err := addToFileStorage(s.repo.(*ShortURLAttr).Cntr, link); err != nil {
		return "", err
	}
	s.repo.(*ShortURLAttr).MapURL[s.repo.(*ShortURLAttr).Cntr] = link
	return s.repo.(*ShortURLAttr).Settings.AdresBase + "/" + strconv.Itoa(s.repo.(*ShortURLAttr).Cntr), nil
}

func (s Service) GetLongLink(id string) (string, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		logger.Log.Info("GetLongLink", zap.String("Atoi", "error"))
		return "", err
	}
	longURL, ok := s.repo.(*ShortURLAttr).MapURL[i]
	if !ok {
		return "", errors.New("wrong id")
	}
	return longURL, nil
}

func (s Service) GetAdresRun() string {
	return s.repo.(*ShortURLAttr).Settings.AdresRun
}

func (s Service) GetLogLevel() string {
	return s.repo.(*ShortURLAttr).Settings.LogLevel
}

func NewInMemoryService() Service {
	shu := ShortURLAttr{}
	shu.Settings = *config.NewServerSettings()
	shu.MapURL = make(map[int]string)

	if err := filestorage.NewConsumer(shu.Settings.FileStorage); err != nil {
		logger.Log.Info("NewInMemoryService", zap.String("GetEvents", err.Error()))
	}
	events, err := filestorage.GetConsumer().GetEvents()
	if err != nil {
		logger.Log.Info("NewInMemoryService", zap.String("GetEvents", err.Error()))
	}
	for _, event := range events {
		if event.UUID > math.MaxInt32 {
			event.UUID = math.MaxInt32
			logger.Log.Info("NewInMemoryService", zap.String("GetEvents", "too big event.UUID"))
		}
		shu.MapURL[int(event.UUID)] = event.OriginalURL
	}
	shu.Cntr = len(events)

	if err := filestorage.NewProducer(shu.Settings.FileStorage); err != nil {
		logger.Log.Info("NewInMemoryService", zap.String("GetEvents", err.Error()))
	}
	return Service{repo: &shu}
}

func (shu *ShortURLAttr) AddLink(link string) (string, error) {
	return "proba", nil
}

func addToFileStorage(cntr int, link string) error {
	if cntr < 0 {
		return errors.New("not valid counter")
	}
	event := filestorage.Event{UUID: uint(cntr), ShortURL: strconv.Itoa(cntr), OriginalURL: link}
	if err := filestorage.GetProducer().WriteEvent(&event); err != nil {
		return err
	}
	return nil
}
