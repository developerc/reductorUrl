package general

type ArrLongURL struct {
	CorellationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ArrRepoURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
