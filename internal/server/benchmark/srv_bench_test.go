package server

import (
	//"strings"
	"testing"
	/*"net/http"
	"net/http/httptest"

	"github.com/developerc/reductorUrl/internal/server"
	"github.com/developerc/reductorUrl/internal/service/memory"*/)

// BenchmarkSrv служит для проверки скорости выполнения функций сервера.
func BenchmarkSrv(b *testing.B) {
	/*svc, err := memory.NewInMemoryService()
	if err != nil {
		return
	}
	srv, err := server.NewServer(svc)
	if err != nil {
		return
	}
	tsrv := httptest.NewServer(srv.SetupRoutes())
	defer tsrv.Close()

	b.Run("#1_PostTest_bench", func(b *testing.B) {
		_, err := svc.AddLink("http://blabla1.ru", "user1")
		if err != nil {
			return
		}
	})

	b.Run("#2_PostJSONTest_bench", func(b *testing.B) {
		longURL := strings.NewReader("{\"url\": \"http://blabla3.ru\"}")
		request := httptest.NewRequest(http.MethodPost, "/api/shorten", longURL)
		w := httptest.NewRecorder()
		srv.AddLinkJSON(w, request)
		res := w.Result()
		res.Body.Close()
	})

	b.Run("#3_GetTest_bench", func(b *testing.B) {
		_, _, err := svc.GetLongLink("1")
		if err != nil {
			return
		}
	})*/
}
