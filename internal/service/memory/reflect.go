package memory

// GetCounter возвращает значение счетчика для формирования идентификатора сокращенного URL
func (s *Service) GetCounter() int {
	return s.repo.GetShu().Cntr
}

// IncrCounter увеличивает счетчик для создания нового короткого URL
func (s *Service) IncrCounter() {
	s.repo.GetShu().Cntr++
}

// GetAdresBase получает базовый адрес сервера
func (s *Service) GetAdresBase() string {
	return s.repo.GetShu().Settings.AdresBase
}

// GetDSN получает DSN базы данных
func (s *Service) GetDSN() (string, error) {
	return s.repo.GetShu().Settings.DBStorage, nil
}

// GetLongURL получает длинный URL по ID
func (s *Service) GetLongURL(i int) (string, error) {
	return s.repo.GetShu().MapURL[i], nil
}

// AddLongURL добавляет длинный URL в Map
func (s *Service) AddLongURL(i int, link string) {
	s.repo.GetShu().MapURL[i] = link
}

// AddLongURL получает адрес запуска сервера
func (s *Service) GetAdresRun() string {
	return s.repo.GetShu().Settings.AdresRun
}

// GetLogLevel получает уровень логирования
func (s *Service) GetLogLevel() string {
	return s.repo.GetShu().Settings.LogLevel
}

// GetShortURLAttr получает аттрибуты коротких URL
func (s *Service) GetShortURLAttr() *ShortURLAttr {
	return s.repo.GetShu()
}
