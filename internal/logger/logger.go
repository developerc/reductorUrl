package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

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
	defer zl.Sync()
	Log = zl
	return nil
}
