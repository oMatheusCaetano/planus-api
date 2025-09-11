package main

import (
	"github.com/omatheuscaetano/planus-api/internal/router"
	"github.com/omatheuscaetano/planus-api/pkg/app"
	"github.com/omatheuscaetano/planus-api/pkg/db"
	"github.com/omatheuscaetano/planus-api/pkg/validators"
)

func main() {
    app.Init()
    db.Init()
    validators.Init()
    r := router.Init()
    r.Run(":" + app.ServerPort())
}
