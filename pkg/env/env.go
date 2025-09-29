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

func DBHost() string {
	return Get("DB_HOST")
}

func DBPort() string {
	return Get("DB_CONTAINER_PORT")
}

func DBName() string {
	return Get("DB_NAME")
}

func DBUser() string {
	return Get("DB_USER")
}

func DBPassword() string {
	return Get("DB_PASSWORD")
}

func DBTestHost() string {
	return Get("DB_TEST_HOST")
}

func DBTestPort() string {
	return Get("DB_TEST_CONTAINER_PORT")
}

func DBTestName() string {
	return Get("DB_TEST_NAME")
}

func DBTestUser() string {
	return Get("DB_TEST_USER")
}

func DBTestPassword() string {
	return Get("DB_TEST_PASSWORD")
}

func Timezone() string {
	return Get("TIMEZONE")
}
