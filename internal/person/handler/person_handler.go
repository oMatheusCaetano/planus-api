package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/person/dto"
	"github.com/omatheuscaetano/planus-api/internal/person/service"
	"github.com/omatheuscaetano/planus-api/pkg/responses"
)

type PersonHandler struct {
	service *service.PersonService
}

func NewPersonHandler(service *service.PersonService) *PersonHandler {
	return &PersonHandler{service: service}
}

func (h *PersonHandler) List(c *gin.Context) {
	var dto dto.ListPerson
	if err := c.ShouldBindJSON(&dto); err != nil {
		responses.BadRequest(c, err)
		return
	}

	resources, appErr := h.service.All(c.Request.Context(), &dto)
	if appErr != nil {
		responses.Error(c, appErr)
		return
	}

	responses.Ok(c, resources)
}

func (h *PersonHandler) Find(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.BadRequest(c, err)
		return
	}

	resource, appErr := h.service.Find(c.Request.Context(), id)
	if appErr != nil {
		responses.Error(c, appErr)
		return
	}

	responses.Ok(c, resource)
}

func (h *PersonHandler) Create(c *gin.Context) {
	var dto dto.CreatePerson
	if err := c.ShouldBindJSON(&dto); err != nil {
		responses.BadRequest(c, err)
		return
	}

	resource, err := h.service.Create(c.Request.Context(), &dto)
	if err != nil {
		responses.Error(c, err)
		return
	}

	responses.Created(c, resource)
}

func (h *PersonHandler) Update(c *gin.Context) {
	var dto dto.UpdatePerson
	if err := c.ShouldBindJSON(&dto); err != nil {
		responses.BadRequest(c, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.BadRequest(c, err)
		return
	}

	resource, appErr := h.service.Update(c.Request.Context(), id, &dto)
	if appErr != nil {
		responses.Error(c, appErr)
		return
	}

	responses.Ok(c, resource)
}

func (h *PersonHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		responses.BadRequest(c, err)
		return
	}

	appErr := h.service.Delete(c.Request.Context(), id)
	if appErr != nil {
		responses.Error(c, appErr)
		return
	}

	responses.Ok(c, nil)
}
