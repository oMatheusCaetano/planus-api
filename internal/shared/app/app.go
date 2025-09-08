package app

import "os"

func Name() string {
	return os.Getenv("APP_NAME")
}
