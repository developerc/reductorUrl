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
	val := reflect.ValueOf(s.repo)
	intCntr := val.Elem().FieldByName("Cntr").Int()

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

func (s Service) GetDSN() (string, error) {
	val := reflect.ValueOf(s.repo)
	settings := val.Elem().FieldByName("Settings")
	dsn := settings.FieldByName("DBStorage").String()
	if len(dsn) == 0 {
		return "", errors.New("get wrong DSN")
	}
	return dsn, nil
}

func (s Service) GetLongURL(i int) (string, error) {
	val := reflect.ValueOf(s.repo)
	mapURL := val.Elem().FieldByName("MapURL")
	if mapURL.Kind() != reflect.Map {
		return "", errors.New("error, this is not a Map")
	}
	return mapURL.MapIndex(reflect.ValueOf(i)).String(), nil
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
}

func (s Service) GetLogLevel() string {
	val := reflect.ValueOf(s.repo)
	settings := val.Elem().FieldByName("Settings")
	adresBase := settings.FieldByName("LogLevel")
	return adresBase.String()
}
