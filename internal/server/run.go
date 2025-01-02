package server

import (
	"net/http"

	"github.com/developerc/reductorUrl/internal/service/memory"
	"go.uber.org/zap"
)

//var server *Server

func Run() error {
	service, err := memory.NewInMemoryService()
	if err != nil {
		return err
	}
	/*zapLogger, err := logger.Initialize(service.GetLogLevel())
	if err != nil {
		return err
	}

	zapLogger.Info("Running server", zap.String("address", service.GetAdresRun()))*/

	server, err := NewServer(service)
	if err != nil {
		return err
	}
	server.logger.Info("Running server", zap.String("address", service.GetAdresRun()))

	routes := server.SetupRoutes()
	err = http.ListenAndServe(service.GetAdresRun(), routes) //nolint:gosec // unnessesary error checking
	return err
}

func (s *Server) GetServer() *Server {
	return s
}

/*func (s *Server) GetService() *memory.Service{
	val := reflect.ValueOf(s.service)
	service := val.Elem().FieldByName("Service")

}*/
