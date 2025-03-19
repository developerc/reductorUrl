package server

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/developerc/reductorUrl/internal/service/memory"
)

func Run() error {
	service, err := memory.NewInMemoryService()
	if err != nil {
		return err
	}

	server, err := NewServer(service)
	if err != nil {
		return err
	}
	server.logger.Info("Running server", zap.String("address", service.GetAdresRun()))

	routes := server.SetupRoutes()
	err = http.ListenAndServe(service.GetAdresRun(), routes) //nolint:gosec // unnessesary error checking
	return err
}
