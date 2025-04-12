package server

import (
	//"net/http"

	"context"
	"fmt"
	"os"
	"os/signal"

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
	//---
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		<-sigint
		fmt.Println("Получен сигнал о прерывании работы")
		if err := server.httpSrv.Shutdown(context.Background()); err != nil {
			// ошибки закрытия Listener
			server.logger.Info("Shutdown server", zap.String("error", err.Error()))
			//log.Printf("HTTP server Shutdown: %v", err)
		}
	}()

	//---
	server.httpSrv.Addr = service.GetAdresRun()
	server.httpSrv.Handler = routes
	if service.GetShortURLAttr().Settings.EnableHTTPS {
		err = server.httpSrv.ListenAndServeTLS(service.GetShortURLAttr().Settings.CertFile, service.GetShortURLAttr().Settings.KeyFile)
	} else {
		err = server.httpSrv.ListenAndServe()
	}

	//---
	//---
	//idleConnsClosed := make(chan struct{})
	/*sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		<-sigint
		fmt.Println("Получен сигнал о прерывании работы")

		//close(idleConnsClosed)
	}()*/
	//---

	/*if service.GetShortURLAttr().Settings.EnableHTTPS {
		err = http.ListenAndServeTLS(service.GetAdresRun(), service.GetShortURLAttr().Settings.CertFile, service.GetShortURLAttr().Settings.KeyFile, routes)
	} else {
		err = http.ListenAndServe(service.GetAdresRun(), routes) //nolint:gosec // unnessesary error checking
	}*/
	//---
	//<-idleConnsClosed
	//fmt.Println("Server Shutdown gracefully")

	return err
}
