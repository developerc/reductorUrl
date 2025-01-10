package memory

import "reflect"

func (s *Service) GetCounter() int {
	shu := s.GetShortURLAttr()
	return shu.Cntr
}

func (s *Service) IncrCounter() {
	shu := s.GetShortURLAttr()
	shu.Cntr++
}

func (s *Service) GetAdresBase() string {
	shu := s.GetShortURLAttr()
	return shu.Settings.AdresBase
}

func (s *Service) GetDSN() (string, error) {
	shu := s.GetShortURLAttr()
	return shu.Settings.DBStorage, nil
}

func (s *Service) GetLongURL(i int) (string, error) {
	shu := s.GetShortURLAttr()
	return shu.MapURL[i], nil
}

func (s *Service) AddLongURL(i int, link string) {
	shu := s.GetShortURLAttr()
	shu.MapURL[i] = link
}

func (s *Service) GetAdresRun() string {
	shu := s.GetShortURLAttr()
	return shu.Settings.AdresRun
}

func (s *Service) GetLogLevel() string {
	shu := s.GetShortURLAttr()
	return shu.Settings.LogLevel
}

func (s *Service) GetShortURLAttr() *ShortURLAttr {
	val := reflect.ValueOf(s.repo)
	return (*ShortURLAttr)(val.UnsafePointer())
}
