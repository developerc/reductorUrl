package memory

import (
	"errors"
	"log"
	"reflect"

	"github.com/developerc/reductorUrl/internal/config"
)

type ShortURLAttr struct {
	Settings config.ServerSettings
	Cntr     int
	MapURL   map[int]string
}

func (s Service) GetCounter() int {
	//fmt.Printf("repo: %v , type %T\n", s.repo, s.repo)
	val := reflect.ValueOf(s.repo)
	//reflect.Value(s.repo)
	//cntr := val.FieldByName("Cntr")
	//fmt.Println("val.Kind(): ", val.Kind())
	//fmt.Printf("Pointer to %v : %v\n", val.Elem().Type(), val.Elem())
	//fmt.Println("cntr: ", val.Elem().FieldByName("Cntr"))
	//fmt.Println("CanSet: ", val.Elem().FieldByName("Cntr").CanSet())
	intCntr := val.Elem().FieldByName("Cntr").Int()

	//fmt.Println("intCntr: ", intCntr)
	return int(intCntr)
}

func (s Service) IncrCounter() {
	val := reflect.ValueOf(s.repo)
	field := val.Elem().FieldByName("Cntr")
	intCntr := field.Int()
	intCntr++
	if field.CanSet() {
		field.SetInt(intCntr)
	}
}

func (s Service) GetAdresBase() string {
	val := reflect.ValueOf(s.repo)
	settings := val.Elem().FieldByName("Settings")
	adresBase := settings.FieldByName("AdresBase")
	return adresBase.String()
}

func (s Service) GetLongURL(i int) (string, error) {
	val := reflect.ValueOf(s.repo)
	mapURL := val.Elem().FieldByName("MapURL")
	if mapURL.Kind() != reflect.Map {
		return "", errors.New("error, this is not a Map")
	}
	return mapURL.MapIndex(reflect.ValueOf(i)).String(), nil
	//fmt.Println(str)
}

func (s Service) AddLongURL(i int, link string) {
	val := reflect.ValueOf(s.repo)
	mapURL := val.Elem().FieldByName("MapURL")
	if mapURL.Kind() != reflect.Map {
		log.Println("Error, this is not a Map")
		return
	}
	if !mapURL.CanSet() {
		log.Println("Error, can not set a Map")
		return
	}
	mapURL.SetMapIndex(reflect.ValueOf(i), reflect.ValueOf(link))
}

func (s Service) GetAdresRun() string {
	val := reflect.ValueOf(s.repo)
	settings := val.Elem().FieldByName("Settings")
	adresBase := settings.FieldByName("AdresRun")
	return adresBase.String()
	//return s.repo.(*ShortURLAttr).Settings.AdresRun
}

func (s Service) GetLogLevel() string {
	val := reflect.ValueOf(s.repo)
	settings := val.Elem().FieldByName("Settings")
	adresBase := settings.FieldByName("LogLevel")
	return adresBase.String()
	//return s.repo.(*ShortURLAttr).Settings.LogLevel
}
