package server

import (
	"net/http"

	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/service/memory"
	"go.uber.org/zap"
)

var service memory.Service
var server Server

func Run() error {
	service = memory.NewInMemoryService()

	if err := logger.Initialize(service.GetLogLevel()); err != nil {
		return err
	}

	logger.GetLog().Info("Running server", zap.String("address", service.GetAdresRun()))

	server = NewServer(service)

	routes := server.SetupRoutes()
	err := http.ListenAndServe(service.GetAdresRun(), routes) //nolint:gosec // unnessesary error checking
	return err
}

func GetService() memory.Service {
	return service
}

func GetServer() Server {
	return server
}
