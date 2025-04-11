// general пакет для размещения типов данных, общих для нескольких пакетов.
// Служит для ухода от перекрестных зависимостей.
package general

// ArrLongURL структура списка длинных URL.
type ArrLongURL struct {
	CorellationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ArrRepoURL структура списка URL присланных пользователем.
type ArrRepoURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
