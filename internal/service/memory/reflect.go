package memory

// GetCounter возвращает значение счетчика для формирования идентификатора сокращенного URL
func (s *Service) GetCounter() int {
	return s.repo.GetShu().Cntr
}

func (s *Service) IncrCounter() {
	s.repo.GetShu().Cntr++
}

func (s *Service) GetAdresBase() string {
	return s.repo.GetShu().Settings.AdresBase
}

func (s *Service) GetDSN() (string, error) {
	return s.repo.GetShu().Settings.DBStorage, nil
}

func (s *Service) GetLongURL(i int) (string, error) {
	return s.repo.GetShu().MapURL[i], nil
}

func (s *Service) AddLongURL(i int, link string) {
	s.repo.GetShu().MapURL[i] = link
}

func (s *Service) GetAdresRun() string {
	return s.repo.GetShu().Settings.AdresRun
}

func (s *Service) GetLogLevel() string {
	return s.repo.GetShu().Settings.LogLevel
}

func (s *Service) GetShortURLAttr() *ShortURLAttr {
	return s.repo.GetShu()
}
