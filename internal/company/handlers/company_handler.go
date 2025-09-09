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
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		responses.Error(c, err)
		return
	}

	company, appErr := h.service.Find(id)
	if appErr != nil {
		responses.Error(c, appErr)
		return
	}

	c.JSON(http.StatusOK, company)
}