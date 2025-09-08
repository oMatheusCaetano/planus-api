package config

import (
	"github.com/joho/godotenv"
)

func initEnv() error {
	err := godotenv.Load()

	if err != nil {
		return err
	}

	return nil
}
