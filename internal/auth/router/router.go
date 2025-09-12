package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/auth/handler"
	"github.com/omatheuscaetano/planus-api/internal/auth/service"
	"github.com/omatheuscaetano/planus-api/internal/auth/store"
	personStore "github.com/omatheuscaetano/planus-api/internal/person/store"
	"github.com/omatheuscaetano/planus-api/pkg/db"
)

func Routes(r *gin.RouterGroup) {

	personStore := personStore.NewPersonPgStore(db.GetDB())
	authStore := store.NewAuthPgStore(db.GetDB())
	authService := service.NewAuthService(authStore, personStore)
	authHandler := handler.NewAuthHandler(authService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
	}
}
