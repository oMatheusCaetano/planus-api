package main

import (
	"context"

	"github.com/omatheuscaetano/planus-api/database/seed"
	authStore "github.com/omatheuscaetano/planus-api/internal/auth/store"
	personStore "github.com/omatheuscaetano/planus-api/internal/person/store"
	"github.com/omatheuscaetano/planus-api/pkg/app"
	"github.com/omatheuscaetano/planus-api/pkg/db"
)

func main() {
    app.Init()
    db.Init()

    personStore := personStore.NewPersonPgStore(db.GetDB())
    authStore   := authStore.NewAuthPgStore(db.GetDB())
	personSeeder := seed.NewSeeder(personStore, authStore)
    personSeeder.Generate(context.Background(), 1393, true)
    personSeeder.Generate(context.Background(), 784, false)
}
