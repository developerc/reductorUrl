// api пакет обработки API запросов
package api

import (
	"bytes"
	"encoding/json"

	"go.uber.org/zap"
)

// LongURL структура длинный URL
type LongURL struct {
	URL string `json:"url"`
}

// ShortURL структура короткий URL
type ShortURL struct {
	Result string `json:"result"`
}

// HandleAPIShorten функция демаршаллинга длинного URL
func HandleAPIShorten(buf bytes.Buffer, logger *zap.Logger) (string, error) {
	var longURL LongURL
	if err := json.Unmarshal(buf.Bytes(), &longURL); err != nil {
		logger.Info("HandleApiShorten", zap.String("error", "demarshalling"))
		return "", err
	}

	return longURL.URL, nil
}

// ShortToJSON функция маршаллинга короткого URL
func ShortToJSON(strShortURL string, logger *zap.Logger) ([]byte, error) {
	shortURL := ShortURL{Result: strShortURL}
	jsonBytes, err := json.Marshal(shortURL)
	if err != nil {
		logger.Info("HandleApiShorten", zap.String("error", "marshaling"))
		return nil, err
	}
	return jsonBytes, nil
}
