package server

import (
	//"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/developerc/reductorUrl/internal/service/memory"
)

// TestPost тестирует работу функций сервера.
func TestPost(t *testing.T) {
	svc, err := memory.NewInMemoryService()
	require.NoError(t, err)
	srv, err := NewServer(svc)
	require.NoError(t, err)
	tsrv := httptest.NewServer(srv.SetupRoutes())
	var cookie *http.Cookie
	defer tsrv.Close()

	t.Run("#1_PostTest", func(t *testing.T) {
		shortURL, err := svc.AddLink("http://blabla.ru", "user1")
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/1", shortURL)
	})

	t.Run("#2_HandleCookieTest", func(t *testing.T) {
		cookie, _, err = svc.HandleCookie("")
		require.NoError(t, err)
		fmt.Println(cookie)
	})

	t.Run("#3_PostJSONTest", func(t *testing.T) {
		longURL := strings.NewReader("{\"url\": \"http://blabla2.ru\"}")
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", longURL)
		w := httptest.NewRecorder()
		srv.AddLinkJSON(w, request)
		res := w.Result()
		res.Body.Close()
		assert.Equal(t, 201, res.StatusCode)
	})

	t.Run("#4_GetTest", func(t *testing.T) {
		resp, _, err := svc.GetLongLink("1")
		require.NoError(t, err)
		assert.Equal(t, "http://blabla.ru", resp)
	})

	/*t.Run("#5_Ping", func(t *testing.T) {
		err := svc.Ping()
		require.NoError(t, err)
	})

	t.Run("#6_GetUserURLs", func(t *testing.T) {
		jsonBytes, err := svc.FetchURLs(cookie.Value)
		require.NoError(t, err)
		assert.Equal(t, "[{\"short_url\":\"http://localhost:8080/1\",\"original_url\":\"http://blabla.ru\"},{\"short_url\":\"http://localhost:8080/2\",\"original_url\":\"http://blabla2.ru\"}]", string(jsonBytes))

	})

	t.Run("#7_DelURLs", func(t *testing.T) {
		var b bytes.Buffer
		b.WriteString("[\"1\"]")
		_, err = svc.DelURLs(cookie.Value, b)
		require.NoError(t, err)
	})

	t.Run("#8_PostBatchURLs", func(t *testing.T) {
		var b bytes.Buffer
		b.WriteString("[{\"correlation_id\":\"ident1\",\"original_url\":\"http://blabla17.ru\"}]")
		jsonBytes, err := svc.HandleBatchJSON(b, "user1")
		require.NoError(t, err)
		assert.Equal(t, "[{\"correlation_id\":\"ident1\",\"short_url\":\"http://localhost:8080/3\"}]", string(jsonBytes))
	})*/
}
