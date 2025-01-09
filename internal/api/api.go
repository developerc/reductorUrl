package api

import (
	"bytes"
	"encoding/json"

	"go.uber.org/zap"
)

type LongURL struct {
	URL string `json:"url"`
}

type ShortURL struct {
	Result string `json:"result"`
}

func HandleAPIShorten(buf bytes.Buffer, logger *zap.Logger) (string, error) {
	var longURL LongURL
	if err := json.Unmarshal(buf.Bytes(), &longURL); err != nil {
		logger.Info("HandleApiShorten", zap.String("error", "demarshalling"))
		return "", err
	}

	return longURL.URL, nil
}

func ShortToJSON(strShortURL string, logger *zap.Logger) ([]byte, error) {
	shortURL := ShortURL{Result: strShortURL}
	jsonBytes, err := json.Marshal(shortURL)
	if err != nil {
		logger.Info("HandleApiShorten", zap.String("error", "marshalling"))
		return nil, err
	}
	return jsonBytes, nil
}
