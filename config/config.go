package config

import (
	"flag"
	"fmt"
)

type settings struct {
	AdresRun  string
	AdresBase string
}

var ServerSettings settings

func InitSettings() {
	ServerSettings = settings{}
	// указываем имя флага, значение по умолчанию и описание
	ServerSettings.AdresRun = *flag.String("a", "localhost:8080", "address running server")
	ServerSettings.AdresBase = *flag.String("b", "http://localhost:8080", "base address shortener URL")
	// делаем разбор командной строки
	flag.Parse()
	fmt.Println("AdresBase: ", ServerSettings.AdresBase)
	fmt.Println("AdresRun: ", ServerSettings.AdresRun)
}
