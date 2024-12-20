package memory

import (
	"errors"
	"log"
	"math"
	"strconv"

	"github.com/developerc/reductorUrl/internal/config"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
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

var service Service

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
		log.Println(err)
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

func NewInMemoryService() *Service {
	if service.repo != nil {
		return &service
	}
	shu := ShortURLAttr{}
	shu.Settings = *config.NewServerSettings()
	shu.MapURL = make(map[int]string)

	if err := filestorage.NewConsumer(shu.Settings.FileStorage); err != nil {
		log.Println(err)
	}
	events, err := filestorage.GetConsumer().ListEvents()
	if err != nil {
		log.Println(err)
	}
	for _, event := range events {
		if event.UUID > math.MaxInt32 {
			event.UUID = math.MaxInt32
		}
		shu.MapURL[int(event.UUID)] = event.OriginalURL
	}
	shu.Cntr = len(events)

	if err := filestorage.NewProducer(shu.Settings.FileStorage); err != nil {
		log.Println(err)
	}
	service = Service{repo: &shu}
	return &service
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
