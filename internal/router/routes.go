package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/shared/app"
)

func initRoutes(r *gin.Engine) {
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Bem-vindo à API %s!", app.Name())})
	})
}
