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

var service Service
var shu ShortURLAttr

func (s Service) AddLink(link string) (string, error) {
	s.IncrCounter()
	if err := addToFileStorage(s.GetCounter(), link); err != nil {
		return "", err
	}
	s.AddLongURL(s.GetCounter(), link)
	return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
}

func (s Service) GetLongLink(id string) (string, error) {
	i, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return "", err
	}
	longURL, err := s.GetLongURL(i)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func NewInMemoryService() *Service {
	if service.repo != nil {
		return &service
	}
	shu = ShortURLAttr{}
	shu.Settings = *config.NewServerSettings()
	shu.MapURL = make(map[int]string)

	if _, err := filestorage.NewConsumer(shu.Settings.FileStorage); err != nil {
		log.Println(err)
	}
	consumer, err := filestorage.NewConsumer(shu.Settings.FileStorage)
	if err != nil {
		log.Println(err)
	}
	events, err := consumer.ListEvents()
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

	if _, err := filestorage.NewProducer(shu.Settings.FileStorage); err != nil {
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
	producer, err := filestorage.NewProducer(shu.Settings.FileStorage)
	if err != nil {
		return err
	}
	producer.WriteEvent(&event)
	return nil
}
