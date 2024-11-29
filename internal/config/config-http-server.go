package config

import (
	"flag"
	"fmt"
	"os"
)

type ServerSettings struct {
	AdresRun  string
	AdresBase string
}

var srvSetGlob ServerSettings

func GetSrvSetGlob() *ServerSettings {
	return &srvSetGlob
}

func NewServerSettings() *ServerSettings {
	srvSetGlob = ServerSettings{}
	// указываем имя флага, значение по умолчанию и описание
	ar := flag.String("a", "localhost:8080", "address running server")
	ab := flag.String("b", "http://localhost:8080", "base address shortener URL")
	// делаем разбор командной строки
	flag.Parse()
	//проверим переменную окружения SERVER_ADDRESS
	val, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || val == "" {
		srvSetGlob.AdresRun = *ar
		fmt.Println("AdresRun from flag: ", srvSetGlob.AdresRun)
	} else {
		srvSetGlob.AdresRun = val
		fmt.Println("AdresRun from env: ", srvSetGlob.AdresRun)
	}

	//проверим переменную окружения BASE_URL
	val, ok = os.LookupEnv("BASE_URL")
	if !ok || val == "" {
		srvSetGlob.AdresBase = *ab
		fmt.Println("AdresBase from flag: ", srvSetGlob.AdresBase)
	} else {
		srvSetGlob.AdresBase = val
		fmt.Println("AdresBase from env: ", srvSetGlob.AdresBase)
	}
	return &srvSetGlob
}
