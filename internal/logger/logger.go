package logger

import (
	"log"

	"go.uber.org/zap"
)

type Logger interface{}

type ZapLogger struct {
	log    Logger
	zapLog *zap.Logger
}

var zapLogger ZapLogger

func Initialize(level string) (*zap.Logger, error) {
	if zapLogger.log != nil {
		return zapLogger.zapLog, nil
	}
	zapLogger = ZapLogger{zapLog: zap.NewNop()}
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	if err := zl.Sync(); err != nil {
		log.Println(err)
	}
	zapLogger.zapLog = zl
	return zapLogger.zapLog, nil
}
