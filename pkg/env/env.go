package env

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"github.com/omatheuscaetano/planus-api/pkg/path"
)

func Load() *errs.Error {
	err := godotenv.Load(path.RootDir(".env"))
	if err != nil {
		return errs.From(err)
	}
	return nil
}

func Get(key string) string {
	return os.Getenv(key)
}

func ApiPort() string {
	return Get("API_PORT")
}

func AppName() string {
	return Get("APP_NAME")
}
