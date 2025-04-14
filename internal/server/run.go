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
	idleConnsClosed := make(chan struct{})
	//idleConnsClosed2 := make(chan struct{})
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
		fmt.Println("Получен сигнал о прерывании работы")
		os.Exit(0)
		server.httpSrv.Shutdown(ctx)
		close(idleConnsClosed)
		fmt.Println("сервер остановлен")
		//os.Exit(0)
	}()

	go func() {
		<-idleConnsClosed
		fmt.Println("останавливаем БД")
		err = service.CloseDB()
		if err != nil {
			server.logger.Info("Close DB", zap.String("error", err.Error()))
		}
		fmt.Println("БД остановлена")
		//close(idleConnsClosed2)
		//os.Exit(0)
	}()

	/*go func() {
		<-idleConnsClosed2
		fmt.Println("закрываем приложение")
		time.Sleep(time.Second)
		//os.Exit(0)
	}()*/

	/*sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		fmt.Println("Получен сигнал о прерывании работы")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.httpSrv.Shutdown(ctx)
		fmt.Println("server shutted down")*/
	//close(idleConnsClosed)

	/*err = service.CloseDB()
	if err != nil {
		server.logger.Info("Close DB", zap.String("error", err.Error()))
	}*/
	//os.Exit(0)
	/*if err := server.httpSrv.Shutdown(context.Background()); err != nil {
		server.logger.Info("Shutdown server", zap.String("error", err.Error()))
	}*/

	//}()

	server.httpSrv.Addr = service.GetAdresRun()
	server.httpSrv.Handler = routes
	if service.GetShortURLAttr().Settings.EnableHTTPS {
		err = server.httpSrv.ListenAndServeTLS(service.GetShortURLAttr().Settings.CertFile, service.GetShortURLAttr().Settings.KeyFile)
	} else {
		err = server.httpSrv.ListenAndServe()
	}

	//<-idleConnsClosed

	return err
}
