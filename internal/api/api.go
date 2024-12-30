package api

import (
	"bytes"
	"encoding/json"

	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/service/memory"
	"go.uber.org/zap"
)

/*type ArrLongURL struct {
	CorellationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}*/

type LongURL struct {
	URL string `json:"url"`
}

type ShortURL struct {
	Result string `json:"result"`
}

/*func HandleBatchJSON(buf bytes.Buffer) (string, error) {
	arrLongURL := make([]ArrLongURL, 0)
	if err := json.Unmarshal(buf.Bytes(), &arrLongURL); err != nil {
		zapLogger, err := logger.Initialize(memory.NewInMemoryService().GetLogLevel())
		zapLogger.Info("HandleBatchJSON", zap.String("error", "demarshalling"))
		return "", err
	}
	fmt.Println(arrLongURL)
	return "", nil
}*/

func HandleAPIShorten(buf bytes.Buffer) (string, error) {
	var longURL LongURL
	if err := json.Unmarshal(buf.Bytes(), &longURL); err != nil {
		zapLogger, err := logger.Initialize(memory.NewInMemoryService().GetLogLevel())
		if err != nil {
			return "", err
		}
		zapLogger.Info("HandleApiShorten", zap.String("error", "demarshalling"))
		return "", err
	}

	return longURL.URL, nil
}

func ShortToJSON(strShortURL string) ([]byte, error) {
	shortURL := ShortURL{Result: strShortURL}
	jsonBytes, err := json.Marshal(shortURL)
	if err != nil {
		zapLogger, err := logger.Initialize(memory.NewInMemoryService().GetLogLevel())
		if err != nil {
			return nil, err
		}
		zapLogger.Info("HandleApiShorten", zap.String("error", "marshalling"))
		return nil, err
	}
	return jsonBytes, nil
}
