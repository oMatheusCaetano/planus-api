package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/person/handler"
	"github.com/omatheuscaetano/planus-api/internal/person/service"
	"github.com/omatheuscaetano/planus-api/internal/person/store"
	"github.com/omatheuscaetano/planus-api/pkg/db"
)

func Routes(r *gin.RouterGroup) {

	personStore := store.NewPersonPgStore(db.GetDB())
	personService := service.NewPersonService(personStore)
	personHandler := handler.NewPersonHandler(personService)

	personGroup := r.Group("/person")
	{
		personGroup.POST("/list", personHandler.All)
		personGroup.GET("/:id", personHandler.Find)
		personGroup.POST("", personHandler.Create)
		personGroup.PUT("/:id", personHandler.Update)
		personGroup.DELETE("/:id", personHandler.Delete)
	}
}
