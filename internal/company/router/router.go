package router

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/company/handlers"
	"github.com/omatheuscaetano/planus-api/internal/company/repositories"
	"github.com/omatheuscaetano/planus-api/internal/company/services"
	db "github.com/omatheuscaetano/planus-api/internal/database"
)

func InitRoutes(r *gin.RouterGroup) {
	sqlDB := db.GetDB()
	companyRepo := repositories.NewCompanyRepository(sqlDB)
	companyService := services.NewCompanyService(companyRepo)
	companyHandler := handlers.NewCompanyHandler(companyService)

	company := r.Group("/company")
	{
		company.GET("/:id", companyHandler.Find)
		company.POST("", companyHandler.Create)
	}
}