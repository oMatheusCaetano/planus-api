package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
)

func AppRouter(g *gin.Engine, i *dependency_container.AppInstances) {
	g.GET("/", i.Handlers.App.Welcome)
}
