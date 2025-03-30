// logger пакет служит для размещения инструмента логирования.
package logger

import (
	"log"

	"go.uber.org/zap"
)

// Initialize конструктор логгера
func Initialize(level string) (*zap.Logger, error) {
	var err error
	var zapLevel zap.AtomicLevel
	if level == "Info" {
		zapLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	} else {
		zapLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zapLevel
	zapConfig.OutputPaths = []string{"stderr"}
	zapLogger, _ := zapConfig.Build()
	if err = zapLogger.Sync(); err != nil {
		log.Println(err)
	}

	return zapLogger, err
}
