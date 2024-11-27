package config

import (
	"flag"
	"fmt"
	"os"
)

type settings struct {
	AdresRun  string
	AdresBase string
}

var ServerSettings settings

func InitSettings() {
	ServerSettings = settings{}
	// указываем имя флага, значение по умолчанию и описание
	ar := flag.String("a", "localhost:8080", "address running server")
	ab := flag.String("b", "http://localhost:8080", "base address shortener URL")
	// делаем разбор командной строки
	flag.Parse()
	//проверим переменную окружения SERVER_ADDRESS
	val, ok := os.LookupEnv("SERVER_ADDRESS")
	if !ok || val == "" {
		ServerSettings.AdresRun = *ar
		fmt.Println("AdresRun from flag: ", ServerSettings.AdresRun)
	} else {
		ServerSettings.AdresRun = val
		fmt.Println("AdresRun from env: ", ServerSettings.AdresRun)
	}

	//проверим переменную окружения BASE_URL
	val, ok = os.LookupEnv("BASE_URL")
	if !ok || val == "" {
		ServerSettings.AdresBase = *ab
		fmt.Println("AdresBase from flag: ", ServerSettings.AdresBase)
	} else {
		ServerSettings.AdresBase = val
		fmt.Println("AdresBase from env: ", ServerSettings.AdresBase)
	}

}
