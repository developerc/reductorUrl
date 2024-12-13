package api

import (
	"bytes"
	"encoding/json"

	"github.com/developerc/reductorUrl/internal/logger"
	"go.uber.org/zap"
)

type LongURL struct {
	URL string `json:"url"`
}

type ShortURL struct {
	Result string `json:"result"`
}

func HandleApiShorten(buf bytes.Buffer) (string, error) {
	var longURL LongURL
	// десериализуем
	if err := json.Unmarshal(buf.Bytes(), &longURL); err != nil {
		logger.Log.Info("HandleApiShorten", zap.String("error", "demarshalling"))
		return "", err
	}

	return longURL.URL, nil
}

func ShortToJSON(strShortURL string) ([]byte, error) {
	shortURL := ShortURL{Result: strShortURL}
	jsonBytes, err := json.Marshal(shortURL)
	if err != nil {
		logger.Log.Info("ShortToJSON", zap.String("error", "marshalling"))
		return nil, nil
	}
	//fmt.Println(string(jsonBytes))
	return jsonBytes, nil
}
