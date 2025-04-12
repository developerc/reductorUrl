// config пакет функций для получения параметров конфигурации запуска приложения
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

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
	EnableHTTPS bool
	CertFile    string
	KeyFile     string
	TypeStorage TypeStorage
}

// ConfigJSON структура для получения настроечных данных из JSON файла.
type ConfigJSON struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DataBaseDsn     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

// String метод возвращает тип хранилища данных
func (ts TypeStorage) String() string {
	return [...]string{"MemoryStorage", "FileStorage", "DBStorage"}[ts]
}

// NewServerSettings конструктор объекта хранения настроечных данных сервера
//
//gocyclo:ignore
func NewServerSettings() *ServerSettings {
	const (
		defaultConfig = "internal/config/configJSON.txt"
		usage         = "configuration by JSON file"
	)
	var fileJSON string
	var err error
	serverSettings := ServerSettings{}
	serverSettings.Logger, err = logger.Initialize("Info")
	if err != nil {
		log.Println(err)
	}
	serverSettings.TypeStorage = MemoryStorage

	ar := flag.String("a", "", "address running server")
	ab := flag.String("b", "", "base address shortener URL")
	logLevel := flag.String("l", "info", "log level")
	fileStorage := flag.String("f", "file_storage.txt", "file for storage data")
	dbStorage := flag.String("d", "", "address connect to DB")
	enableHTTPS := flag.Bool("s", false, "enable HTTPS")
	certFile := flag.String("cf", "certs/localhost.pem", "certificat file")
	keyFile := flag.String("kf", "certs/localhost-key.pem", "key file")
	flag.StringVar(&fileJSON, "c", defaultConfig, usage)
	flag.StringVar(&fileJSON, "config", defaultConfig, usage)
	flag.Parse()

	configJSON := getConfigJSON(fileJSON)
	fmt.Println(configJSON)

	val, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || val == "" {
		if !isFlagPassed("a") {
			serverSettings.AdresRun = configJSON.ServerAddress
			serverSettings.Logger.Info("AdresRun from fileJSON:", zap.String("address", serverSettings.AdresRun))
		} else {
			serverSettings.AdresRun = *ar
			serverSettings.Logger.Info("AdresRun from flag:", zap.String("address", serverSettings.AdresRun))
		}
	} else {
		serverSettings.AdresRun = val
		serverSettings.Logger.Info("AdresRun from env:", zap.String("address", serverSettings.AdresRun))
	}

	val, ok = os.LookupEnv("BASE_URL")
	if !ok || val == "" {
		if !isFlagPassed("b") {
			serverSettings.AdresBase = configJSON.BaseURL
			serverSettings.Logger.Info("AdresBase from fileJSON:", zap.String("address", serverSettings.AdresBase))
		} else {
			serverSettings.AdresBase = *ab
			serverSettings.Logger.Info("AdresBase from flag:", zap.String("address", serverSettings.AdresBase))
		}
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
		if !isFlagPassed("f") {
			serverSettings.FileStorage = configJSON.FileStoragePath
			serverSettings.Logger.Info("FileStorage from fileJSON:", zap.String("storage", serverSettings.FileStorage))
		} else {
			serverSettings.FileStorage = *fileStorage
			if isFlagPassed("f") && (serverSettings.FileStorage != "") {
				serverSettings.TypeStorage = FileStorage
			}
			if serverSettings.FileStorage == "" {
				serverSettings.FileStorage = "file_storage.txt"
			}
			serverSettings.Logger.Info("FileStorage from flag:", zap.String("storage", serverSettings.FileStorage))
		}

	} else {
		serverSettings.TypeStorage = FileStorage
		serverSettings.FileStorage = val
		serverSettings.Logger.Info("FileStorage from env:", zap.String("storage", serverSettings.FileStorage))
	}

	val, ok = os.LookupEnv("DATABASE_DSN")
	if !ok || val == "" {
		if !isFlagPassed("d") {
			serverSettings.DBStorage = configJSON.DataBaseDsn
			serverSettings.Logger.Info("DbStorage from fileJSON:", zap.String("storage", serverSettings.DBStorage))
		} else {
			serverSettings.DBStorage = *dbStorage
			if isFlagPassed("d") && (serverSettings.DBStorage != "") {
				serverSettings.TypeStorage = DBStorage
			}
			serverSettings.Logger.Info("DbStorage from flag:", zap.String("storage", serverSettings.DBStorage))
		}

	} else {
		serverSettings.TypeStorage = DBStorage
		serverSettings.DBStorage = val
		serverSettings.Logger.Info("DbStorage from env:", zap.String("storage", serverSettings.DBStorage))
	}

	serverSettings.Logger.Info("serverSettings.TypeStorage:", zap.String("storage", serverSettings.TypeStorage.String()))

	val, ok = os.LookupEnv("ENABLE_HTTPS")
	if !ok || val == "" {
		if !isFlagPassed("s") {
			serverSettings.EnableHTTPS = configJSON.EnableHTTPS
			serverSettings.Logger.Info("EnableHTTPS from fileJSON:", zap.String("enableHTTPS", strconv.FormatBool(serverSettings.EnableHTTPS)))
		} else {
			serverSettings.EnableHTTPS = *enableHTTPS
			serverSettings.Logger.Info("EnableHTTPS from flag:", zap.String("enableHTTPS", strconv.FormatBool(serverSettings.EnableHTTPS)))
		}

	} else {
		serverSettings.EnableHTTPS, err = strconv.ParseBool(val)
		if err != nil {
			serverSettings.EnableHTTPS = false
		}
		serverSettings.Logger.Info("EnableHTTPS from env:", zap.String("enableHTTPS", strconv.FormatBool(serverSettings.EnableHTTPS)))
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

func getConfigJSON(fileJSON string) *ConfigJSON {
	var configJSON ConfigJSON
	b, err := os.ReadFile(fileJSON)
	if err != nil {
		return nil
	}
	if err = json.Unmarshal(b, &configJSON); err != nil {
		return nil
	}
	return &configJSON
}
