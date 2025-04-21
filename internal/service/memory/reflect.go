package memory

import "fmt"

// GetCounter возвращает значение счетчика для формирования идентификатора сокращенного URL
func (s *Service) GetCounter() int {
	return s.shu.Cntr
}

// IncrCounter увеличивает счетчик для создания нового короткого URL
func (s *Service) IncrCounter() {
	s.shu.Cntr++
}

// GetAdresBase получает базовый адрес сервера
func (s *Service) GetAdresBase() string {
	return s.shu.Settings.AdresBase
}

// GetDSN получает DSN базы данных
func (s *Service) GetDSN() (string, error) {
	return s.shu.Settings.DBStorage, nil
}

// GetLongURL получает длинный URL по ID
func (s *Service) GetLongURL(i int) (string, error) {
	return s.shu.MapURL[i].OriginalURL, nil
}

// AddLongURL добавляет длинный URL в Map
func (s *Service) AddLongURL(i int, link, usr string) {
	mapURLVal := MapURLVal{OriginalURL: link, Usr: usr}
	s.shu.MapURL[i] = mapURLVal
	fmt.Println("usr: ", usr)
	fmt.Println(s.shu.MapURL)
}

// AddLongURL получает адрес запуска сервера
func (s *Service) GetAdresRun() string {
	return s.shu.Settings.AdresRun
}

// GetLogLevel получает уровень логирования
func (s *Service) GetLogLevel() string {
	return s.shu.Settings.LogLevel
}

// GetShortURLAttr получает аттрибуты коротких URL
func (s *Service) GetShortURLAttr() *ShortURLAttr {
	return s.shu
}
