package handlers

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/company/models"
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

func parseWhereCondition(condition any) db.WhereLogicBlock {
	var conditionOrBlock string

	conditionMap := condition.(map[string]interface{})

	if cond, ok := conditionMap["condition"].([]interface{}); ok {
		var conditions []db.WhereLogicBlock
		for _, c := range cond {
			conditions = append(conditions, parseWhereCondition(c))
		}
		return db.WhereLogicBlock{
			Operator:  conditionMap["operator"].(string),
			Condition: conditions,
		}
	} else if cond, ok := conditionMap["condition"].(map[string]interface{}); ok {
		conditionOrBlock = cond["operator"].(string)
		return db.WhereLogicBlock{
			Operator: conditionMap["operator"].(string),
			Condition: db.Where{
				Key:      cond["key"].(string),
				Operator: conditionOrBlock,
				Value:    cond["value"],
			},
		}
	}

	return db.WhereLogicBlock{}
}

func (h *CompanyHandler) All(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		responses.BadRequest(c, err)
		return
	}

	//? Where Conditions
	var conditions []db.WhereLogicBlock
	if where, ok := body["where"]; ok {
		for _, whereItem := range where.([]interface{}) {
			conditions = append(conditions, parseWhereCondition(whereItem) )
		}
	}

	//? Sort By
	var sortByDto []db.SortBy
	sortByQuery := body["sort_by"]

	if sortByQuery != "" && sortByQuery != nil {
		sortBy := strings.Split(sortByQuery.(string), ",")

		for _, sort := range sortBy {
			split := strings.Split(sort, ":")
			key := split[0]
			direction := "asc"
			if len(split) > 1 {
				direction = split[1]
			}
			sortByDto = append(sortByDto, db.SortBy{
				Key:       key,
				Direction: direction,
			})
		}
	}

	//? Pagination Params
	p := body["page"]
	pp := body["per_page"]
	var page, perPage int
	if p != nil {
		switch v := p.(type) {
		case float64:
			page = int(v)
		case int:
			page = v
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				page = i
			}
		}
	}
	if pp != nil {
		switch v := pp.(type) {
		case float64:
			perPage = int(v)
		case int:
			perPage = v
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				perPage = i
			}
		}
	}

	//?When Listing all without pagination
	if page == 0 && perPage == 0 {
		data, err := h.service.All(dto.ListingProps{
			SortBy: sortByDto,
			Where:  conditions,
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
		Where:  conditions,
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
