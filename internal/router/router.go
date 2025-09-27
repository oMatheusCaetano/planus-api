package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
)

func AllRoutes(g *gin.Engine, i *dependency_container.Instances) {
	AppRouter(g, i.App)
}
