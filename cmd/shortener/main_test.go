package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/developerc/reductorUrl/internal/app"
	"github.com/developerc/reductorUrl/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostHandler(t *testing.T) {
	settings := config.NewServerSettings()
	shortURLAttr := app.NewShortURLAttr(*settings)

	t.Run("#1_PostTest", func(t *testing.T) {
		longURL := strings.NewReader("http://blabla.ru")
		request := httptest.NewRequest(http.MethodPost, "/", longURL)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(app.PostHandler(*shortURLAttr))
		h(w, request)

		result := w.Result()

		assert.Equal(t, 201, result.StatusCode)
		defer result.Body.Close()
		resBody, err := io.ReadAll(result.Body)
		require.NoError(t, err)
		assert.Equal(t, "text/plain", result.Header.Get("Content-Type"))
		assert.Equal(t, `http://localhost:8080/1`, string(resBody))
	})

	t.Run("#2_BadTest", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodDelete, "/", nil)
		w := httptest.NewRecorder()
		app.BadHandler(w, request)

		result := w.Result()

		assert.Equal(t, 400, result.StatusCode)
		defer result.Body.Close()
	})
}
