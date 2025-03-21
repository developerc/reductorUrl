// general пакет для размещения типов данных, общих для нескольких пакетов.
// Служит для ухода от перекрестных зависимостей.
package general

type ArrLongURL struct {
	CorellationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ArrRepoURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
