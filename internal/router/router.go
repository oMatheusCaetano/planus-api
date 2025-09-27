package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"

	_ "github.com/omatheuscaetano/planus-api/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AllRoutes(g *gin.Engine, i *dependency_container.Instances) {
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	AppRouter(g, i.App)
	UserRouter(g, i.User)
}
