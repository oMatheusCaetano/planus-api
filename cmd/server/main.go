package main

import (
	"github.com/omatheuscaetano/planus-api/internal/config"
	"github.com/omatheuscaetano/planus-api/internal/router"
	"github.com/omatheuscaetano/planus-api/internal/shared/app"
	"github.com/omatheuscaetano/planus-api/internal/shared/validators"
)

func main() {
	config.Init()
	validators.Init()
	r := router.Init()
	r.Run(":" + app.ServerPort())
}
