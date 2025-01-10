package memory

import "reflect"

func (s *Service) GetCounter() int {
	shu := s.GetShortURLAttr()
	return shu.Cntr
	/*val := reflect.ValueOf(s.repo)
	intCntr := val.Elem().FieldByName("Cntr").Int()

	return int(intCntr)*/
}

func (s *Service) IncrCounter() {
	shu := s.GetShortURLAttr()
	shu.Cntr++
	/*val := reflect.ValueOf(s.repo)
	field := val.Elem().FieldByName("Cntr")
	intCntr := field.Int()
	intCntr++
	if field.CanSet() {
		field.SetInt(intCntr)
	}*/
}

func (s *Service) GetAdresBase() string {
	shu := s.GetShortURLAttr()
	return shu.Settings.AdresBase
	/*val := reflect.ValueOf(s.repo)
	settings := val.Elem().FieldByName("Settings")
	adresBase := settings.FieldByName("AdresBase")
	return adresBase.String()*/
}

func (s *Service) GetDSN() (string, error) {
	shu := s.GetShortURLAttr()
	return shu.Settings.DBStorage, nil
	/*val := reflect.ValueOf(s.repo)
	settings := val.Elem().FieldByName("Settings")
	dsn := settings.FieldByName("DBStorage").String()
	if dsn == "" {
		return "", errors.New("get wrong DSN")
	}
	return dsn, nil*/
}

func (s *Service) GetLongURL(i int) (string, error) {
	shu := s.GetShortURLAttr()
	return shu.MapURL[i], nil
	/*val := reflect.ValueOf(s.repo)
	mapURL := val.Elem().FieldByName("MapURL")
	if mapURL.Kind() != reflect.Map {
		return "", errors.New("error, this is not a Map")
	}
	return mapURL.MapIndex(reflect.ValueOf(i)).String(), nil*/
}

func (s *Service) AddLongURL(i int, link string) {
	shu := s.GetShortURLAttr()
	shu.MapURL[i] = link
	/*val := reflect.ValueOf(s.repo)
	mapURL := val.Elem().FieldByName("MapURL")
	if mapURL.Kind() != reflect.Map {
		log.Println("Error, this is not a Map")
		return
	}
	if !mapURL.CanSet() {
		log.Println("Error, can not set a Map")
		return
	}
	mapURL.SetMapIndex(reflect.ValueOf(i), reflect.ValueOf(link))*/
}

func (s *Service) GetAdresRun() string {
	shu := s.GetShortURLAttr()
	return shu.Settings.AdresRun
	/*val := reflect.ValueOf(s.repo)
	settings := val.Elem().FieldByName("Settings")
	adresBase := settings.FieldByName("AdresRun")
	return adresBase.String()*/
}

func (s *Service) GetLogLevel() string {
	shu := s.GetShortURLAttr()
	return shu.Settings.LogLevel
	/*val := reflect.ValueOf(s.repo)
	settings := val.Elem().FieldByName("Settings")
	adresBase := settings.FieldByName("LogLevel")
	return adresBase.String()*/
}

func (s *Service) GetShortURLAttr() *ShortURLAttr {
	/*shu := s.GetShortURLAttr()
	return shu*/
	val := reflect.ValueOf(s.repo)
	return (*ShortURLAttr)(val.UnsafePointer())
}
