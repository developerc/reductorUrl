package config

import (
	"flag"
	"log"
	"os"
)

type SrvSettings interface{}

type ServerSettings struct {
	ss          SrvSettings
	AdresRun    string
	AdresBase   string
	LogLevel    string
	FileStorage string
}

var serverSettings ServerSettings

func NewServerSettings() *ServerSettings {
	if serverSettings.ss != nil {
		return &serverSettings
	}
	serverSettings = ServerSettings{}
	ar := flag.String("a", "localhost:8080", "address running server")
	ab := flag.String("b", "http://localhost:8080", "base address shortener URL")
	logLevel := flag.String("l", "info", "log level")
	fileStorage := flag.String("f", "file_storage.txt", "file for storage data")
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
		log.Println("FileStorage from flag:", serverSettings.FileStorage)
	} else {
		serverSettings.FileStorage = val
		log.Println("FileStorage from env:", serverSettings.FileStorage)
	}

	return &serverSettings
}
