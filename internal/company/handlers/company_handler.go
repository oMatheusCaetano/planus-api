package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/company/models"
	"github.com/omatheuscaetano/planus-api/internal/company/services"
	"github.com/omatheuscaetano/planus-api/internal/shared/responses"
)

type CompanyHandler struct {
	service services.CompanyService
}

func NewCompanyHandler(service services.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) All(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))

	if page == 0 && perPage == 0 {
		data, err := h.service.All()
		if err != nil {
			responses.Error(c, err)
			return
		}
		responses.Ok(c, data, nil)
		return
	}

	paginated, err := h.service.Paginate(1, 10)
	if err != nil {
		responses.Error(c, err)
		return
	}
	responses.Ok(c, paginated.Data, paginated.Meta)
}

func (h *CompanyHandler) Find(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	company, err := h.service.Find(id)
	if err != nil {
		responses.Error(c, err)
		return
	}
	responses.Ok(c, company, nil)
}

func (h *CompanyHandler) Create(c *gin.Context) {
	var company models.Company

	if err := c.ShouldBindJSON(&company); err != nil {
		responses.BadRequest(c, err)
		return
	}

	err := h.service.Create(&company)
	if err != nil {
		responses.Error(c, err)
		return
	}

	responses.Ok(c, company, nil)
}
