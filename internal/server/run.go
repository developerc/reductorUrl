package server

import (
	"net/http"

	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/service/memory"
	"go.uber.org/zap"
)

var server *Server

func Run() error {
	service := memory.NewInMemoryService()

	zapLogger, err := logger.Initialize(service.GetLogLevel())
	if err != nil {
		return err
	}

	zapLogger.Info("Running server", zap.String("address", service.GetAdresRun()))

	server = NewServer(service)

	routes := server.SetupRoutes()
	err = http.ListenAndServe(service.GetAdresRun(), routes) //nolint:gosec // unnessesary error checking
	return err
}

func GetServer() *Server {
	return server
}
