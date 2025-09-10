package handlers

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/company/services"
	"github.com/omatheuscaetano/planus-api/internal/db"
	"github.com/omatheuscaetano/planus-api/internal/shared/dto"
	"github.com/omatheuscaetano/planus-api/internal/shared/responses"
)

type CompanyHandler struct {
	service services.CompanyService
}

func NewCompanyHandler(service services.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

func (h *CompanyHandler) All(c *gin.Context) {
	//? Sort By
	var sortByDto []db.SortBy
	sortByQuery := c.Query("sort_by")
	
	if sortByQuery != "" {
		sortBy := strings.Split(sortByQuery, ",")

		for _, sort := range sortBy {
			split := strings.Split(sort, ":")
			key := split[0]
			direction := "asc"
			if len(split) > 1 {
				direction = split[1]
			}
			log.Println("Sort:", key, direction, sortBy)
			sortByDto = append(sortByDto, db.SortBy{
				Key:       key,
				Direction: direction,
			})
		}
	}


	//? Pagination Params
	page, _ := strconv.Atoi(c.Query("page"))
	perPage, _ := strconv.Atoi(c.Query("per_page"))

	//?When Listing all without pagination
	if page == 0 && perPage == 0 {
		data, err := h.service.All(dto.ListingProps{
			SortBy: sortByDto,
		})
		if err != nil {
			responses.Error(c, err)
			return
		}
		responses.Ok(c, data, nil)
		return
	}

	//? Paginated Result
	paginated, err := h.service.Paginate(dto.PaginationProps{
		Page:    page,
		PerPage: perPage,
		SortBy:  sortByDto,
	})
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

// func (h *CompanyHandler) Create(c *gin.Context) {
// 	var company models.Company

// 	if err := c.ShouldBindJSON(&company); err != nil {
// 		responses.BadRequest(c, err)
// 		return
// 	}

// 	err := h.service.Create(&company)
// 	if err != nil {
// 		responses.Error(c, err)
// 		return
// 	}

// 	responses.Ok(c, company, nil)
// }
