package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	"github.com/developerc/reductorUrl/internal/server"
	"github.com/developerc/reductorUrl/internal/service/memory"
)

// BenchmarkSrv служит для проверки скорости выполнения функций сервера.
func BenchmarkSrv(b *testing.B) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	svc, err := memory.NewInMemoryService(ctx)
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
		_, err := svc.AddLink(ctx, "http://blabla1.ru", "user1")
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
		_, _, err := svc.GetLongLink(ctx, "1")
		if err != nil {
			return
		}
	})
}
