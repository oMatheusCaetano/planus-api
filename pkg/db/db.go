package db

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
	"github.com/omatheuscaetano/planus-api/pkg/app"
)

var (
	con  *sql.DB
	once sync.Once
)

func GetDB() *sql.DB {
	return con
}

func Init() {
	once.Do(func() {
		var err error
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
            app.DBHost(), app.DBUser(), app.DBPassword(), app.DBName(), app.DBPort(),
        )
		con, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal(err)
		}
		if err := con.Ping(); err != nil {
			log.Fatal(err)
		}
	})
}
