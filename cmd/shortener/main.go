// main Главный пакет. Точка входа
package main

import (
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/developerc/reductorUrl/internal/server"
)

var (
	//go:embed version.txt
	buildVersion string
	//go:embed date.txt
	buildDate string
	//go:embed commit.txt
	buildCommit string
)

func main() {
	BuildVersion := strings.TrimSpace(buildVersion)
	if len(BuildVersion) > 0 {
		fmt.Printf("Build version: %q\n", BuildVersion)
	} else {
		fmt.Printf("Build version: N/A\n")
	}

	BuildDate := strings.TrimSpace(buildDate)
	if len(BuildDate) > 0 {
		fmt.Printf("Build date: %q\n", BuildDate)
	} else {
		fmt.Printf("Build date: N/A\n")
	}

	BuildCommit := strings.TrimSpace(buildCommit)
	if len(BuildCommit) > 0 {
		fmt.Printf("Build commit: %q\n", BuildCommit)
	} else {
		fmt.Printf("Build commit: N/A\n")
	}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
