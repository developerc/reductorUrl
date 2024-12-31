package memory

import (
	"bytes"
	"errors"
	"log"

	//"math"
	"strconv"

	"github.com/developerc/reductorUrl/internal/config"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
	//db "github.com/developerc/reductorUrl/internal/service/db_storage"
)

type repository interface {
	AddLink(link string) (string, error)
}

type Service struct {
	repo repository
}

var service Service

//var shu *ShortURLAttr

func (s *Service) AddLink(link string) (string, error) {
	var shURL string
	var err error
	s.IncrCounter()
	switch s.repo.(*ShortURLAttr).Settings.TypeStorage {
	case "MemoryStorage":
		{
			s.AddLongURL(s.GetCounter(), link)
			return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
		}

	case "FileStorage":
		{
			if err = s.repo.(*ShortURLAttr).addToFileStorage(s.GetCounter(), link); err != nil {
				return "", err
			}
			s.AddLongURL(s.GetCounter(), link)
			return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
		}
	case "DBStorage":
		{
			//log.Println("link: ", link)
			shURL, err = insertRecord(s.repo.(*ShortURLAttr), link)
			if err != nil {
				return "", err
			}
			// для DBStorage
			//s.CreateTable()
		}
	}

	//s.AddLongURL(s.GetCounter(), link)
	//return s.GetAdresBase() + "/" + strconv.Itoa(s.GetCounter()), nil
	return s.GetAdresBase() + "/" + shURL, nil
}

func (s Service) GetShortByOriginalURL(originalURL string) (string, error) {
	shortURL, err := getShortByOriginalURL(s.repo.(*ShortURLAttr), originalURL)
	return s.GetAdresBase() + "/" + shortURL, err
}

func (s Service) GetLongLink(id string) (string, error) {
	var longURL string
	i, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		return "", err
	}
	switch s.repo.(*ShortURLAttr).Settings.TypeStorage {
	case "MemoryStorage":
		{
			longURL, err = s.GetLongURL(i)
			if err != nil {
				return "", err
			}
			//return longURL, nil
		}
	case "FileStorage":
		{
			longURL, err = s.GetLongURL(i)
			if err != nil {
				return "", err
			}
			//return longURL, nil
		}
	case "DBStorage":
		{
			longURL, err = getLongByUUID(s.repo.(*ShortURLAttr), i)
			if err != nil {
				return "", err
			}
		}
	}
	/*longURL, err := s.GetLongURL(i)
	if err != nil {
		return "", err
	}*/
	return longURL, nil
}

func (s Service) HandleBatchJSON(buf bytes.Buffer) ([]byte, error) {
	arrLongURL, err := listLongURL(buf)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//fmt.Println(arrLongURL)

	//проверка на нулевую длину arrLongURL
	if len(arrLongURL) == 0 {
		return nil, errors.New("error: length array is zero")
	}

	jsonBytes, err := s.handleArrLongURL(arrLongURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return jsonBytes, nil
}

func NewInMemoryService() *Service {
	if service.repo != nil {
		return &service
	}
	//shu = ShortURLAttr{}
	shu := new(ShortURLAttr)
	shu.Settings = *config.NewServerSettings()
	shu.MapURL = make(map[int]string)

	switch shu.Settings.TypeStorage {
	case "FileStorage":
		getFileSettings(shu)
	case "DBStorage":
		createTable(shu)
		//db.CreateTable()
	}

	service = Service{repo: shu}
	//service := new(Service)
	//service.repo = shu
	return &service
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
	producer.WriteEvent(&event)
	return nil
}
