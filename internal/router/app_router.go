package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
)

func AppRouter(g *gin.RouterGroup, i *dependency_container.Instances) {
	g.GET("/", i.App.Handlers.App.Welcome)
}
