// config пакет функций для получения параметров конфигурации запуска приложения
package config

import (
	"flag"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/developerc/reductorUrl/internal/logger"
)

// TypeStorage enum переменная для определения типа хранилища данных
type TypeStorage int

// TypeStorage определение переменной
const (
	MemoryStorage TypeStorage = iota
	FileStorage
	DBStorage
)

// ServerSettings структура для хранения настроечных данных сервера
type ServerSettings struct {
	Logger      *zap.Logger
	AdresRun    string
	AdresBase   string
	LogLevel    string
	FileStorage string
	DBStorage   string
	EnableHTTPS string
	CertFile    string
	KeyFile     string
	TypeStorage TypeStorage
}

// String метод возвращает тип хранилища данных
func (ts TypeStorage) String() string {
	return [...]string{"MemoryStorage", "FileStorage", "DBStorage"}[ts]
}

// NewServerSettings конструктор объекта хранения настроечных данных сервера
//
//gocyclo:ignore
func NewServerSettings() *ServerSettings {
	var err error
	serverSettings := ServerSettings{}
	serverSettings.Logger, err = logger.Initialize("Info")
	if err != nil {
		log.Println(err)
	}
	serverSettings.TypeStorage = MemoryStorage

	ar := flag.String("a", "localhost:8080", "address running server")
	ab := flag.String("b", "http://localhost:8080", "base address shortener URL")
	logLevel := flag.String("l", "info", "log level")
	fileStorage := flag.String("f", "file_storage.txt", "file for storage data")
	dbStorage := flag.String("d", "", "address connect to DB")
	enableHTTPS := flag.String("s", "", "enable HTTPS")
	certFile := flag.String("cf", "certs/localhost.pem", "certificat file")
	keyFile := flag.String("kf", "certs/localhost-key.pem", "key file")
	flag.Parse()

	val, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || val == "" {
		serverSettings.AdresRun = *ar
		serverSettings.Logger.Info("AdresRun from flag:", zap.String("address", serverSettings.AdresRun))
	} else {
		serverSettings.AdresRun = val
		serverSettings.Logger.Info("AdresRun from env:", zap.String("address", serverSettings.AdresRun))
	}

	val, ok = os.LookupEnv("BASE_URL")
	if !ok || val == "" {
		serverSettings.AdresBase = *ab
		serverSettings.Logger.Info("AdresBase from flag:", zap.String("address", serverSettings.AdresBase))
	} else {
		serverSettings.AdresBase = val
		serverSettings.Logger.Info("AdresBase from env:", zap.String("address", serverSettings.AdresBase))
	}

	val, ok = os.LookupEnv("LOG_LEVEL")
	if !ok || val == "" {
		serverSettings.LogLevel = *logLevel
		serverSettings.Logger.Info("LogLevel from flag:", zap.String("level", serverSettings.LogLevel))
	} else {
		serverSettings.LogLevel = val
		serverSettings.Logger.Info("LogLevel from env:", zap.String("level", serverSettings.LogLevel))
	}

	val, ok = os.LookupEnv("FILE_STORAGE_PATH")
	if !ok || val == "" {
		serverSettings.FileStorage = *fileStorage
		if isFlagPassed("f") && (serverSettings.FileStorage != "") {
			serverSettings.TypeStorage = FileStorage
		}
		if serverSettings.FileStorage == "" {
			serverSettings.FileStorage = "file_storage.txt"
		}
		serverSettings.Logger.Info("FileStorage from flag:", zap.String("storage", serverSettings.FileStorage))
	} else {
		serverSettings.TypeStorage = FileStorage
		serverSettings.FileStorage = val
		serverSettings.Logger.Info("FileStorage from env:", zap.String("storage", serverSettings.FileStorage))
	}

	val, ok = os.LookupEnv("DATABASE_DSN")
	if !ok || val == "" {
		serverSettings.DBStorage = *dbStorage
		if isFlagPassed("d") && (serverSettings.DBStorage != "") {
			serverSettings.TypeStorage = DBStorage
		}
		serverSettings.Logger.Info("DbStorage from flag:", zap.String("storage", serverSettings.DBStorage))
	} else {
		serverSettings.TypeStorage = DBStorage
		serverSettings.DBStorage = val
		serverSettings.Logger.Info("DbStorage from env:", zap.String("storage", serverSettings.DBStorage))
	}

	serverSettings.Logger.Info("serverSettings.TypeStorage:", zap.String("storage", serverSettings.TypeStorage.String()))

	val, ok = os.LookupEnv("ENABLE_HTTPS")
	if !ok || val == "" {
		serverSettings.EnableHTTPS = *enableHTTPS
		serverSettings.Logger.Info("EnableHTTPS from flag:", zap.String("enableHTTPS", serverSettings.EnableHTTPS))
	} else {
		serverSettings.EnableHTTPS = val
		serverSettings.Logger.Info("EnableHTTPS from env:", zap.String("enableHTTPS", serverSettings.EnableHTTPS))
	}

	val, ok = os.LookupEnv("CERT_FILE")
	if !ok || val == "" {
		serverSettings.CertFile = *certFile
		serverSettings.Logger.Info("certFile from flag:", zap.String("certFile", serverSettings.CertFile))
	} else {
		serverSettings.CertFile = val
		serverSettings.Logger.Info("certFile from env:", zap.String("certFile", serverSettings.CertFile))
	}

	val, ok = os.LookupEnv("KEY_FILE")
	if !ok || val == "" {
		serverSettings.KeyFile = *keyFile
		serverSettings.Logger.Info("keyFile from flag:", zap.String("keyFile", serverSettings.KeyFile))
	} else {
		serverSettings.KeyFile = val
		serverSettings.Logger.Info("keyFile from env:", zap.String("keyFile", serverSettings.KeyFile))
	}
	return &serverSettings
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
