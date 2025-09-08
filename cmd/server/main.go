package main

import (
	"log"

	"github.com/omatheuscaetano/planus-api/internal/config"
	"github.com/omatheuscaetano/planus-api/internal/shared/app"
)

func main() {
	log.Println("[INITIALIZING SERVER...]\n")

	log.Println("[INITIALIZING CONFIG...]")
	config.Init()
	log.Println("[CONFIG INITIALIZED SUCCESSFULLY]\n")

	log.Println("[SERVER INITIALIZED SUCCESSFULLY]\n")

	log.Printf("Welcome to %s!\n", app.Name())
}
