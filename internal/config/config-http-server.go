package config

import (
	"flag"
	"fmt"
	"os"
)

type ServerSettings struct {
	AdresRun  string
	AdresBase string
	LogLevel  string
}

var srvSetGlob ServerSettings

func GetSrvSetGlob() *ServerSettings {
	return &srvSetGlob
}

func NewServerSettings() *ServerSettings {
	srvSetGlob = ServerSettings{}
	ar := flag.String("a", "localhost:8080", "address running server")
	ab := flag.String("b", "http://localhost:8080", "base address shortener URL")
	logLevel := flag.String("l", "info", "log level")
	flag.Parse()

	val, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || val == "" {
		srvSetGlob.AdresRun = *ar
		fmt.Println("AdresRun from flag: ", srvSetGlob.AdresRun)
	} else {
		srvSetGlob.AdresRun = val
		fmt.Println("AdresRun from env: ", srvSetGlob.AdresRun)
	}

	val, ok = os.LookupEnv("BASE_URL")
	if !ok || val == "" {
		srvSetGlob.AdresBase = *ab
		fmt.Println("AdresBase from flag: ", srvSetGlob.AdresBase)
	} else {
		srvSetGlob.AdresBase = val
		fmt.Println("AdresBase from env: ", srvSetGlob.AdresBase)
	}

	val, ok = os.LookupEnv("LOG_LEVEL")
	if !ok || val == "" {
		srvSetGlob.LogLevel = *logLevel
		fmt.Println("LogLevel from flag: ", srvSetGlob.LogLevel)
	} else {
		srvSetGlob.LogLevel = val
		fmt.Println("LogLevel from env: ", srvSetGlob.LogLevel)
	}

	return &srvSetGlob
}