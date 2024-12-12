package memory

import (
	//"reductorUrl/internal/service"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/developerc/reductorUrl/internal/config"
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
	//s.repo.Cntr++
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
	var shu ShortURLAttr = ShortURLAttr{}
	shu.Settings = *config.NewServerSettings()
	shu.MapURL = make(map[int]string)
	return Service{repo: &shu}
}

func (shu *ShortURLAttr) AddLink(link string) (string, error) {
	fmt.Println("from shu")
	return "proba", nil
	/*shu.Cntr++
	shu.MapURL[shu.Cntr] = link
	return shu.Settings.AdresBase + "/" + strconv.Itoa(shu.Cntr), nil*/
}

/*func (shu *ShortURLAttr) GetLink(cntr int) (string, error) {
	return shu.MapURL[cntr], nil
}*/
