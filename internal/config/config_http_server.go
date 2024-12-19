package config

import (
	"flag"
	"log"
	"os"
)

type ServerSettings struct {
	AdresRun    string
	AdresBase   string
	LogLevel    string
	FileStorage string
}

func NewServerSettings() *ServerSettings {
	srvSetGlob := ServerSettings{}
	ar := flag.String("a", "localhost:8080", "address running server")
	ab := flag.String("b", "http://localhost:8080", "base address shortener URL")
	logLevel := flag.String("l", "info", "log level")
	fileStorage := flag.String("f", "file_storage.txt", "file for storage data")
	flag.Parse()

	val, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || val == "" {
		srvSetGlob.AdresRun = *ar
		log.Println("AdresRun from flag:", srvSetGlob.AdresRun)
	} else {
		srvSetGlob.AdresRun = val
		log.Println("AdresRun from env:", srvSetGlob.AdresRun)
	}

	val, ok = os.LookupEnv("BASE_URL")
	if !ok || val == "" {
		srvSetGlob.AdresBase = *ab
		log.Println("AdresBase from flag:", srvSetGlob.AdresBase)
	} else {
		srvSetGlob.AdresBase = val
		log.Println("AdresBase from env:", srvSetGlob.AdresBase)
	}

	val, ok = os.LookupEnv("LOG_LEVEL")
	if !ok || val == "" {
		srvSetGlob.LogLevel = *logLevel
		log.Println("LogLevel from flag:", srvSetGlob.LogLevel)
	} else {
		srvSetGlob.LogLevel = val
		log.Println("LogLevel from env:", srvSetGlob.LogLevel)
	}

	val, ok = os.LookupEnv("FILE_STORAGE_PATH")
	if !ok || val == "" {
		srvSetGlob.FileStorage = *fileStorage
		log.Println("FileStorage from flag:", srvSetGlob.FileStorage)
	} else {
		srvSetGlob.FileStorage = val
		log.Println("FileStorage from env:", srvSetGlob.FileStorage)
	}

	return &srvSetGlob
}
