package server

import (
	"testing"

	"net/http/httptest"

	"github.com/developerc/reductorUrl/internal/service/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLink(t *testing.T) {
	svc := memory.NewInMemoryService()
	srv := NewServer(&svc)
	tsrv := httptest.NewServer(srv.SetupRoutes())
	defer tsrv.Close()

	t.Run("#1_PostTest", func(t *testing.T) {
		//shortURL, err := th.svc.AddLink("http://blabla.ru")
		shortURL, err := svc.AddLink("http://blabla.ru")
		require.NoError(t, err)
		assert.Equal(t, "http://localhost:8080/1", shortURL)
	})

	/*t.Run("#2_PostJSONTest", func(t *testing.T) {
		longURL := strings.NewReader("\"url\": \"https://blabla2.ru\"")
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", longURL)
		w := httptest.NewRecorder()
		srv.addLinkJSON(w, request)
		res := w.Result()
		assert.Equal(t, 201, res.StatusCode)
	})*/

	t.Run("#3_GetTest", func(t *testing.T) {
		resp, err := svc.GetLongLink("1")
		require.NoError(t, err)
		assert.Equal(t, "http://blabla.ru", resp)
	})
}
