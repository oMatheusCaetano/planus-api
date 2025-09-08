package config

import (
	"log"
)

func Init() {
	log.Println("Loading environment variables...")

	err := initEnv()
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}
}
