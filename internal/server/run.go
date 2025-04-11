package server

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/developerc/reductorUrl/internal/service/memory"
)

// Run метод запускает работу сервера.
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
	if service.GetShortURLAttr().Settings.EnableHTTPS == "true" {
		//err = http.ListenAndServeTLS(service.GetAdresRun(), "certs/localhost.pem", "certs/localhost-key.pem", routes)
		err = http.ListenAndServeTLS(service.GetAdresRun(), service.GetShortURLAttr().Settings.CertFile, service.GetShortURLAttr().Settings.KeyFile, routes)
	} else {
		err = http.ListenAndServe(service.GetAdresRun(), routes) //nolint:gosec // unnessesary error checking
	}

	return err
}
