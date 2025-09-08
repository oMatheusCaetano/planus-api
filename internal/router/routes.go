package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	company "github.com/omatheuscaetano/planus-api/internal/company/router"
	"github.com/omatheuscaetano/planus-api/internal/shared/app"
)

func initRoutes(r *gin.Engine) {
	apiRoutes := r.Group("/api")

	apiRoutes.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Bem-vindo à API %s!", app.Name())})
	}, )

	company.InitRoutes(apiRoutes)
}
