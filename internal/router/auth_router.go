package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
)

func AuthRouter(g *gin.RouterGroup, i *dependency_container.Instances) {
	g.Group("/auth").
		POST("/login", i.Auth.Handlers.Auth.Login)
}
