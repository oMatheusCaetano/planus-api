package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
	"github.com/omatheuscaetano/planus-api/internal/middleware"
)

func UserRouter(g *gin.RouterGroup, i *dependency_container.UserInstances) {
	group := g.Group("/user")
	group.Use(middleware.JWTMiddleware())
	group.
		GET("/:id", i.Handlers.User.Find).
		POST("", i.Handlers.User.Create).
		DELETE("/:id", i.Handlers.User.Delete)
}
