package server

import (
	"net/http"

	//"reductorUrl/internal/server"
	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/service/memory"
	"go.uber.org/zap"
)

var service memory.Service
var server Server

func Run() error {
	//log.Println("hello")
	//var service memory.Service = memory.NewInMemoryService()
	service = memory.NewInMemoryService()

	if err := logger.Initialize(service.GetLogLevel()); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", service.GetAdresRun()))

	server = NewServer(service)

	routes := server.SetupRoutes()
	err := http.ListenAndServe(service.GetAdresRun(), routes)
	return err
}

func GetService() memory.Service {
	return service
}

func GetServer() Server {
	return server
}
