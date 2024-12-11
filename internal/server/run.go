package server

import (
	"net/http"

	//"reductorUrl/internal/server"
	"github.com/developerc/reductorUrl/internal/service/memory"
)

var service memory.Service

func Run() error {
	//log.Println("hello")
	//var service memory.Service = memory.NewInMemoryService()
	service = memory.NewInMemoryService()
	var server Server = NewServer(service)
	routes := server.SetupRoutes()
	err := http.ListenAndServe(service.GetAdresRun(), routes)
	return err
}

func GetService() memory.Service {
	return service
}
