package memory

import (
	//"reductorUrl/internal/service"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/service/filestorage"
)

type repository interface {
	AddLink(link string) (string, error)
	//GetLongLink(id string) (string, error)
}

type Service struct {
	repo repository
	//shu  *ShortURLAttr
}

type ShortURLAttr struct {
	Settings config.ServerSettings
	Cntr     int
	MapURL   map[int]string
}

// AddLink implements server.svc.
func (s Service) AddLink(link string) (string, error) {
	//fmt.Println("from service")
	s.repo.(*ShortURLAttr).Cntr++
	if err := addToFileStorage(s.repo.(*ShortURLAttr).Cntr, link); err != nil {
		return "", err
	}
	s.repo.(*ShortURLAttr).MapURL[s.repo.(*ShortURLAttr).Cntr] = link
	return s.repo.(*ShortURLAttr).Settings.AdresBase + "/" + strconv.Itoa(s.repo.(*ShortURLAttr).Cntr), nil
}

func (s Service) GetLongLink(id string) (string, error) {
	//log.Println("map: ", s.repo.(*ShortURLAttr))
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

func NewInMemoryService() Service {
	shu := ShortURLAttr{}
	shu.Settings = *config.NewServerSettings()
	shu.MapURL = make(map[int]string)
	//заполняем map значениями из файла file_storage.txt
	filestorage.NewConsumer(shu.Settings.FileStorage)
	events, err := filestorage.GetConsumer().GetEvents()
	if err != nil {
		fmt.Println("error!")
	}
	for _, event := range events {
		shu.MapURL[int(event.UUID)] = event.OriginalURL
	}
	shu.Cntr = len(events)

	filestorage.NewProducer(shu.Settings.FileStorage)
	return Service{repo: &shu}
}

func (shu *ShortURLAttr) AddLink(link string) (string, error) {
	fmt.Println("from shu")
	return "proba", nil
}

func addToFileStorage(cntr int, link string) error {
	event := filestorage.Event{UUID: uint(cntr), ShortURL: strconv.Itoa(cntr), OriginalURL: link}
	if err := filestorage.GetProducer().WriteEvent(&event); err != nil {
		return err
	}
	return nil
}
