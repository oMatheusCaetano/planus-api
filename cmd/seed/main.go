package main

import (
	"context"

	"github.com/omatheuscaetano/planus-api/database/seed"
	"github.com/omatheuscaetano/planus-api/internal/person/store"
	"github.com/omatheuscaetano/planus-api/pkg/app"
	"github.com/omatheuscaetano/planus-api/pkg/db"
)

func main() {
    app.Init()
    db.Init()

    personStore := store.NewPersonPgStore(db.GetDB())
	personSeeder := seed.NewPersonSeed(personStore)
    personSeeder.Generate(context.Background(), 3877)
}
