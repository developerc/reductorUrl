package main

import (
	"log"

	"github.com/developerc/reductorUrl/internal/server"
)

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
