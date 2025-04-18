package server

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/developerc/reductorUrl/internal/general"
	"github.com/developerc/reductorUrl/internal/service/memory"
)

// Run метод запускает работу сервера и мягко останавливает.
func Run() error {
	var needStop bool = false
	//idleConnsClosed := make(chan struct{})
	idleConnsClosed := make(chan bool)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	service, err := memory.NewInMemoryService(ctx)
	if err != nil {
		return err
	}

	server, err := NewServer(service)
	if err != nil {
		return err
	}
	server.logger.Info("Running server", zap.String("address", service.GetAdresRun()))
	routes := server.SetupRoutes()

	go func() {
		<-ctx.Done()
		server.logger.Info("Server", zap.String("shutdown", "begin"))
		server.httpSrv.Shutdown(ctx)
		//close(idleConnsClosed)
		server.logger.Info("Server", zap.String("shutdown", "end"))
		needStop = true
		idleConnsClosed <- true
		//server.logger.Info("first gorutine", zap.String("end", "closing"))
	}()

	go func() {
	L:
		for {
			select {
			case <-idleConnsClosed:
				//fmt.Println("from idleConnsClosed:", needStop, general.CntrAtomVar.GetCntr(), x)
				server.logger.Info("catch idleConnsClosed", zap.String("begin", "succ"))
				if needStop && general.CntrAtomVar.GetCntr() == 0 {
					err = service.CloseDB()
					if err != nil {
						server.logger.Info("Close DB", zap.String("error", err.Error()))
					} else {
						server.logger.Info("Close DB", zap.String("success", "closed"))
					}
					break L
				}
			case <-general.CntrAtomVar.GetChan():
				//fmt.Println("from GetChan:", needStop, general.CntrAtomVar.GetCntr())
				if needStop && general.CntrAtomVar.GetCntr() == 0 {
					err = service.CloseDB()
					if err != nil {
						server.logger.Info("Close DB_", zap.String("error", err.Error()))
					} else {
						server.logger.Info("Close DB_", zap.String("success", "closed"))
					}
					break L
				}
			}
		}
		/*
			<-idleConnsClosed
			server.logger.Info("Close DB", zap.String("begin", "closing"))
			err = service.CloseDB()
			if err != nil {
				server.logger.Info("Close DB", zap.String("error", err.Error()))
			} else {
				server.logger.Info("Close DB", zap.String("success", "closed"))
			}*/
	}()

	server.httpSrv.Addr = service.GetAdresRun()
	server.httpSrv.Handler = routes
	if service.GetShortURLAttr().Settings.EnableHTTPS {
		err = server.httpSrv.ListenAndServeTLS(service.GetShortURLAttr().Settings.CertFile, service.GetShortURLAttr().Settings.KeyFile)
	} else {
		err = server.httpSrv.ListenAndServe()
	}
	if err != nil {
		//if err.Error() == "http: Server closed" {
		if errors.Is(err, http.ErrServerClosed) {
			server.logger.Info("Close server", zap.String("success:", err.Error()))
		} else {
			return err
		}
	}

	time.Sleep(time.Second)
	return nil
}
