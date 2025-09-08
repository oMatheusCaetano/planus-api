package app

import "os"

func Name() string {
	return os.Getenv("APP_NAME")
}

func Mode() string {
	return os.Getenv("APP_MODE")
}

func ServerPort() string {
	return os.Getenv("API_PORT")
}
