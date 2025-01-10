package memory

func (s *Service) GetCounter() int {
	shu := s.repo.GetShu()
	return shu.Cntr
}

func (s *Service) IncrCounter() {
	shu := s.repo.GetShu()
	shu.Cntr++
}

func (s *Service) GetAdresBase() string {
	shu := s.repo.GetShu()
	return shu.Settings.AdresBase
}

func (s *Service) GetDSN() (string, error) {
	shu := s.repo.GetShu()
	return shu.Settings.DBStorage, nil
}

func (s *Service) GetLongURL(i int) (string, error) {
	shu := s.repo.GetShu()
	return shu.MapURL[i], nil
}

func (s *Service) AddLongURL(i int, link string) {
	shu := s.repo.GetShu()
	shu.MapURL[i] = link
}

func (s *Service) GetAdresRun() string {
	shu := s.repo.GetShu()
	return shu.Settings.AdresRun
}

func (s *Service) GetLogLevel() string {
	shu := s.repo.GetShu()
	return shu.Settings.LogLevel
}

func (s *Service) GetShortURLAttr() *ShortURLAttr {
	return s.repo.GetShu()
	/*val := reflect.ValueOf(s.repo)
	return (*ShortURLAttr)(val.UnsafePointer())*/
}
