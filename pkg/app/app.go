package app

import (
	"context"
	"os"

	"github.com/joho/godotenv"
)

type AppContext struct {
	context.Context
	UserID int
}

func Init() error {
	err := godotenv.Load()

	if err != nil {
		return err
	}

	return nil
}

func Name() string {
	return os.Getenv("APP_NAME")
}

func Mode() string {
	return os.Getenv("APP_MODE")
}

func ServerPort() string {
	return os.Getenv("API_PORT")
}

func JWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

func DBHost() string {
	return os.Getenv("DB_HOST")
}

func DBPort() string {
	return os.Getenv("DB_CONTAINER_PORT")
}

func DBUser() string {
	return os.Getenv("DB_USER")
}

func DBPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func DBName() string {
	return os.Getenv("DB_NAME")
}
