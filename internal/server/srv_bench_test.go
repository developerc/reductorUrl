package server

import (
	"testing"

	"net/http/httptest"

	"github.com/developerc/reductorUrl/internal/service/memory"
)

func BenchmarkSrv(b *testing.B) {
	svc, err := memory.NewInMemoryService()
	if err != nil {
		return
	}
	srv, err := NewServer(svc)
	if err != nil {
		return
	}
	tsrv := httptest.NewServer(srv.SetupRoutes())
	defer tsrv.Close()

	b.Run("#1_PostTest_bench", func(b *testing.B) {
		_, err := svc.AddLink("http://blabla3.ru", "user1")
		if err != nil {
			return
		}
	})
}
