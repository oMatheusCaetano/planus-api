package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
)

func UserRouter(g *gin.Engine, i *dependency_container.UserInstances) {
	g.Group("/user").
		GET("/:id", i.Handlers.User.Find).
		POST("", i.Handlers.User.Create).
		DELETE("/:id", i.Handlers.User.Delete)
}
