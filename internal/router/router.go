package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
	"github.com/omatheuscaetano/planus-api/internal/middleware"

	_ "github.com/omatheuscaetano/planus-api/docs"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AllRoutes(g *gin.Engine, i *dependency_container.Instances) {
	api := g.Group("")

	api.Use(middleware.CORSMiddleware())

	api.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	AuthRouter(api, i)
	AppRouter(api, i)
	UserRouter(api, i)
}
