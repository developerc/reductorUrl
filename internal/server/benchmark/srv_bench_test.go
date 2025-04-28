package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os/signal"
	"strings"
	"syscall"
	"testing"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/server"
	"github.com/developerc/reductorUrl/internal/service"
	dbstorage "github.com/developerc/reductorUrl/internal/service/db_storage"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
	"github.com/developerc/reductorUrl/internal/service/memory"
)

// BenchmarkSrv служит для проверки скорости выполнения функций сервера.
func BenchmarkSrv(b *testing.B) {
	var svc *service.Service
	var err error
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	//svc, err := memory.NewInMemoryService(ctx)
	settings := config.NewServerSettings()
	switch settings.TypeStorage {
	case config.MemoryStorage:
		svc, err = memory.NewServiceMemory(ctx, settings)
		if err != nil {
			return
		}
	case config.FileStorage:
		svc, err = filestorage.NewServiceFile(ctx, settings)
		if err != nil {
			return
		}
	case config.DBStorage:
		svc, err = dbstorage.NewServiceDB(ctx, settings)
		if err != nil {
			return
		}
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
