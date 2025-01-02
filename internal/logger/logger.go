package logger

import (
	"log"

	"go.uber.org/zap"
)

func Initialize(level string) (*zap.Logger, error) {
	var err error
	var zapLevel zap.AtomicLevel
	if level == "Info" {
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	//encoder := zap.NewDevelopmentConfig()
	zapConfig := zap.NewDevelopmentConfig()
	//zapConfig.EncoderConfig = encoder
	zapConfig.Level = zapLevel
	zapConfig.OutputPaths = []string{"stderr"}
	zapLogger, err := zapConfig.Build()
	//zapLogger := zap.Must(zapConfig.Build())
	if err := zapLogger.Sync(); err != nil {
		log.Println(err)
	}

	return zapLogger, err
}

/*type Logger interface{}

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
}*/
