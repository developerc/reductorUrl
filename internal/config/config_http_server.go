package config

import (
	"flag"
	"log"
	"os"
)

type SrvSettings interface{}

type ServerSettings struct {
	TypeStorage string
	AdresRun    string
	AdresBase   string
	LogLevel    string
	FileStorage string
	DBStorage   string
}

//gocyclo:ignore
func NewServerSettings() *ServerSettings {
	serverSettings := ServerSettings{}
	serverSettings.TypeStorage = "MemoryStorage"

	ar := flag.String("a", "localhost:8080", "address running server")
	ab := flag.String("b", "http://localhost:8080", "base address shortener URL")
	logLevel := flag.String("l", "info", "log level")
	fileStorage := flag.String("f", "file_storage.txt", "file for storage data")
	dbStorage := flag.String("d", "host=localhost user=admin password=admin dbname=test sslmode=disable", "address connect to DB")
	flag.Parse()

	val, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || val == "" {
		serverSettings.AdresRun = *ar
		log.Println("AdresRun from flag:", serverSettings.AdresRun)
	} else {
		serverSettings.AdresRun = val
		log.Println("AdresRun from env:", serverSettings.AdresRun)
	}

	val, ok = os.LookupEnv("BASE_URL")
	if !ok || val == "" {
		serverSettings.AdresBase = *ab
		log.Println("AdresBase from flag:", serverSettings.AdresBase)
	} else {
		serverSettings.AdresBase = val
		log.Println("AdresBase from env:", serverSettings.AdresBase)
	}

	val, ok = os.LookupEnv("LOG_LEVEL")
	if !ok || val == "" {
		serverSettings.LogLevel = *logLevel
		log.Println("LogLevel from flag:", serverSettings.LogLevel)
	} else {
		serverSettings.LogLevel = val
		log.Println("LogLevel from env:", serverSettings.LogLevel)
	}

	val, ok = os.LookupEnv("FILE_STORAGE_PATH")
	if !ok || val == "" {
		serverSettings.FileStorage = *fileStorage
		if isFlagPassed("f") && (serverSettings.FileStorage != "") {
			serverSettings.TypeStorage = "FileStorage"
		}
		if serverSettings.FileStorage == "" {
			serverSettings.FileStorage = "file_storage.txt"
		}
		log.Println("FileStorage from flag:", serverSettings.FileStorage)
	} else {
		serverSettings.TypeStorage = "FileStorage"
		serverSettings.FileStorage = val
		log.Println("FileStorage from env:", serverSettings.FileStorage)
	}

	val, ok = os.LookupEnv("DATABASE_DSN")
	if !ok || val == "" {
		serverSettings.DBStorage = *dbStorage
		if isFlagPassed("d") && (serverSettings.DBStorage != "") {
			serverSettings.TypeStorage = "DBStorage"
		}
		log.Println("DbStorage from flag:", serverSettings.DBStorage)
	} else {
		serverSettings.TypeStorage = "DBStorage"
		serverSettings.DBStorage = val
		log.Println("DbStorage from env:", serverSettings.DBStorage)
	}

	log.Println("serverSettings.TypeStorage: ", serverSettings.TypeStorage)
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
