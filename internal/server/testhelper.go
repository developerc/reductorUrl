package server

import (
	"net/http/httptest"

	"github.com/developerc/reductorUrl/internal/service/memory"
)

type TestHelper struct {
	srv *httptest.Server
	svc svc
}

func NewTestHelper() TestHelper {
	s := memory.NewInMemoryService()
	srv := NewServer(&s)

	tsrv := httptest.NewServer(srv.SetupRoutes())
	return TestHelper{
		svc: &s,
		srv: tsrv,
	}
}
