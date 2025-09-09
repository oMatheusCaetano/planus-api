package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/company/services"
	"github.com/omatheuscaetano/planus-api/internal/shared/responses"
)

type CompanyHandler struct {
	service services.CompanyService
}

func NewCompanyHandler(service services.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) Find(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	company, err := h.service.Find(id)
	if err != nil {
		responses.Error(c, err)
		return
	}

	c.JSON(http.StatusOK, company)
}