package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
	"github.com/omatheuscaetano/planus-api/internal/middleware"
	"github.com/omatheuscaetano/planus-api/internal/model"
)

func UserRouter(g *gin.RouterGroup, i *dependency_container.Instances) {
	group := g.Group("/user")
	group.Use(middleware.JWTMiddleware())
	group.
		GET(
			"/:id",
			middleware.PermissionMiddleware(i, model.PermissionUserRead),
			i.User.Handlers.User.Find,
		).
		POST(
			"",
			middleware.PermissionMiddleware(i, model.PermissionUserCreate),
			i.User.Handlers.User.Create,
		).
		DELETE(
			"/:id",
			middleware.PermissionMiddleware(i, model.PermissionUserDelete),
			i.User.Handlers.User.Delete,
		)
}
