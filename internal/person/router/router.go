package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/auth/dto"
	middlewares "github.com/omatheuscaetano/planus-api/internal/auth/middleware"
	"github.com/omatheuscaetano/planus-api/internal/person/handler"
	"github.com/omatheuscaetano/planus-api/internal/person/service"
	"github.com/omatheuscaetano/planus-api/internal/person/store"
	"github.com/omatheuscaetano/planus-api/pkg/db"
)

func Routes(r *gin.RouterGroup) {
	personStore := store.NewPersonPgStore(db.GetDB())
	personService := service.NewPersonService(personStore)
	personHandler := handler.NewPersonHandler(personService)

	personGroup := r.Group("/person").Use(middlewares.JWTMiddleware())
	{
		personGroup.POST(
			"/paginate", 
			middlewares.AuthorizeMiddleware([]dto.P{{Module: "person", Action: "read"}}),
			personHandler.Paginate,
		)

		personGroup.POST(
			"/list", 
			middlewares.AuthorizeMiddleware([]dto.P{{Module: "person", Action: "read"}}),
			personHandler.List,
		)

		personGroup.GET(
			"/:id", 
			middlewares.AuthorizeMiddleware([]dto.P{{Module: "person", Action: "read"}}),
			personHandler.Find,
		)

		personGroup.POST(
			"", 
			middlewares.AuthorizeMiddleware([]dto.P{{Module: "person", Action: "create"}}),
			personHandler.Create,
		)

		personGroup.PUT(
			":id",
			middlewares.AuthorizeMiddleware([]dto.P{{Module: "person", Action: "update"}}),
			personHandler.Update,
		)

		personGroup.DELETE(
			":id",
			middlewares.AuthorizeMiddleware([]dto.P{{Module: "person", Action: "delete"}}),
			personHandler.Delete,
		)
	}
}
