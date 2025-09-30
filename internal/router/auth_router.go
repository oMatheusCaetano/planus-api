package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
)

func AuthRouter(g *gin.RouterGroup, i *dependency_container.AuthInstances) {
	g.Group("/auth").
		POST("/login", i.Handlers.Auth.Login)
}
