package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/company/services"
)

type CompanyHandler struct {
	service services.CompanyService
}

func NewCompanyHandler(service services.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) Find(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	company, err := h.service.Find(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Empresa não encontrada"})
		return
	}

	c.JSON(http.StatusOK, company)
}