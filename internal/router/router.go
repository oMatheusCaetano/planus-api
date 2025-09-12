package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	person "github.com/omatheuscaetano/planus-api/internal/person/router"
	auth "github.com/omatheuscaetano/planus-api/internal/auth/router"
	"github.com/omatheuscaetano/planus-api/pkg/app"
	"github.com/omatheuscaetano/planus-api/pkg/middlewares"
	"github.com/omatheuscaetano/planus-api/pkg/responses"
)

func Init() *gin.Engine {
	r := gin.Default()
	initRoutes(r)
	return r
}

func initRoutes(r *gin.Engine) {
	apiRoutes := r.Group("/api")

	apiRoutes.Use(middlewares.AppContextMiddleware())

	apiRoutes.GET("", func(c *gin.Context) {
		responses.Ok(c, gin.H{"message": fmt.Sprintf("Bem-vindo à API %s!", app.Name())})
	})

	person.Routes(apiRoutes)
	auth.Routes(apiRoutes)
}