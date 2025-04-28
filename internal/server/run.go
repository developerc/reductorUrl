package server

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/general"
	"github.com/developerc/reductorUrl/internal/service"
	dbstorage "github.com/developerc/reductorUrl/internal/service/db_storage"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
	"github.com/developerc/reductorUrl/internal/service/memory"
)

// Run метод запускает работу сервера и мягко останавливает.
func Run() error {
	var svc *service.Service
	var err error
	var needStop bool = false
	signalToClose := make(chan struct{})
	beforeStop := make(chan struct{})
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	settings := config.NewServerSettings()

	switch settings.TypeStorage {
	case config.MemoryStorage:
		svc, err = memory.NewServiceMemory(ctx, settings)
		if err != nil {
			return err
		}
	case config.FileStorage:
		svc, err = filestorage.NewServiceFile(ctx, settings)
		if err != nil {
			return err
		}
	case config.DBStorage:
		svc, err = dbstorage.NewServiceDB(ctx, settings)
		if err != nil {
			return err
		}
	}

	server, err := NewServer(svc)
	if err != nil {
		return err
	}
	server.logger.Info("Running server", zap.String("address", svc.Shu.Settings.AdresRun))
	routes := server.SetupRoutes()

	go func() {
		<-ctx.Done()
		server.logger.Info("Server", zap.String("shutdown", "begin"))
		server.httpSrv.Shutdown(ctx)
		server.logger.Info("Server", zap.String("shutdown", "end"))
		needStop = true
		close(signalToClose)
	}()

	go func() {
		defer close(beforeStop)

		for {
			select {
			case <-signalToClose:
				if needStop && general.CntrAtomVar.GetCntr() == 0 {
					err = svc.CloseDB()
					if err != nil {
						server.logger.Info("Close DB", zap.String("error", err.Error()))
					} else {
						server.logger.Info("Close DB", zap.String("success", "closed"))
					}
					return
				}
			case <-general.CntrAtomVar.GetChan():
				if needStop && general.CntrAtomVar.GetCntr() == 0 {
					err = svc.CloseDB()
					if err != nil {
						server.logger.Info("Close DB_", zap.String("error", err.Error()))
					} else {
						server.logger.Info("Close DB_", zap.String("success", "closed"))
					}
					return
				}
			}
		}
	}()

	server.httpSrv.Addr = svc.Shu.Settings.AdresRun
	server.httpSrv.Handler = routes
	if svc.Shu.Settings.EnableHTTPS {
		err = server.httpSrv.ListenAndServeTLS(svc.Shu.Settings.CertFile, svc.Shu.Settings.KeyFile)
	} else {
		err = server.httpSrv.ListenAndServe()
	}
	if err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			server.logger.Info("Close server", zap.String("success:", err.Error()))
		} else {
			return err
		}
	}

	<-beforeStop
	return nil
}
