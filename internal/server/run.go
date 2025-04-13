package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		fmt.Println("Получен сигнал о прерывании работы")
		err = service.CloseDB()
		if err != nil {
			server.logger.Info("Close DB", zap.String("error", err.Error()))
		}
		os.Exit(0)
		if err := server.httpSrv.Shutdown(context.Background()); err != nil {
			server.logger.Info("Shutdown server", zap.String("error", err.Error()))
		}

	}()

	server.httpSrv.Addr = service.GetAdresRun()
	server.httpSrv.Handler = routes
	if service.GetShortURLAttr().Settings.EnableHTTPS {
		err = server.httpSrv.ListenAndServeTLS(service.GetShortURLAttr().Settings.CertFile, service.GetShortURLAttr().Settings.KeyFile)
	} else {
		err = server.httpSrv.ListenAndServe()
	}

	return err
}
