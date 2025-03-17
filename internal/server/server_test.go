package server

import (
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/developerc/reductorUrl/internal/service/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPost(t *testing.T) {
	svc, err := memory.NewInMemoryService()
	require.NoError(t, err)
	srv, err := NewServer(svc)
	require.NoError(t, err)
	tsrv := httptest.NewServer(srv.SetupRoutes())
	defer tsrv.Close()

	t.Run("#1_PostTest", func(t *testing.T) {
		shortURL, err := svc.AddLink("http://blabla.ru", "user1")
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/1", shortURL)
	})

	t.Run("#2_PostJSONTest", func(t *testing.T) {
		longURL := strings.NewReader("{\"url\": \"http://blabla2.ru\"}")
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", longURL)
		w := httptest.NewRecorder()
		srv.AddLinkJSON(w, request)
		res := w.Result()
		res.Body.Close()
		assert.Equal(t, 201, res.StatusCode)
	})

	t.Run("#3_GetTest", func(t *testing.T) {
		resp, _, err := svc.GetLongLink("1")
		require.NoError(t, err)
		assert.Equal(t, "http://blabla.ru", resp)
	})
}
