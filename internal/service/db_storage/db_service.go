package dbstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/service/memory"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func CheckPing() error {
	zapLogger, err := logger.Initialize(memory.NewInMemoryService().GetLogLevel())
	if err != nil {
		return err
	}

	service := memory.NewInMemoryService()
	dsn, err := service.GetDSN()
	if err != nil {
		zapLogger.Info("CheckPing", zap.String("error", err.Error()))
		return err
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		zapLogger.Info("CheckPing", zap.String("error", err.Error()))
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		zapLogger.Info("CheckPing", zap.String("error", err.Error()))
		return err
	}
	return nil
}
