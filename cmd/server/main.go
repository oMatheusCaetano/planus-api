package main

import (
	"log"

	"github.com/omatheuscaetano/planus-api/internal/config"
	"github.com/omatheuscaetano/planus-api/internal/router"
	"github.com/omatheuscaetano/planus-api/internal/shared/app"
)

func main() {
	log.Println("[INITIALIZING SERVER...]\n")

	log.Println("[INITIALIZING CONFIG...]")
	config.Init()
	log.Println("[CONFIG INITIALIZED SUCCESSFULLY]\n")

	log.Println("[INITIALIZING ROUTER...]")
	r := router.Init()
	log.Println("[ROUTER INITIALIZED SUCCESSFULLY]\n")
	r.Run(":" + app.ServerPort())
}
