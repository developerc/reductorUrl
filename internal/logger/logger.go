package logger

import (
	"log"

	"go.uber.org/zap"
)

var zapLog *zap.Logger = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	if err := zl.Sync(); err != nil {
		log.Println(err)
	}
	zapLog = zl
	return nil
}

func GetLog() *zap.Logger {
	return zapLog
}
