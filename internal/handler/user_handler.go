package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dto"
	"github.com/omatheuscaetano/planus-api/internal/model"
	"github.com/omatheuscaetano/planus-api/internal/response"
	"github.com/omatheuscaetano/planus-api/internal/service"
)

type UserHandler struct{
	s *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{s: s}
}

// @Summary Find User
// @Description Returns a user by ID
// @Tags User
// @Success 200 {object} model.User
// @Router /user/{id} [get]
func (h *UserHandler) Find(c *gin.Context) {
	id := c.Param("id")
	user, err := h.s.Find(c.Request.Context(), model.ID(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Ok(c, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.s.Delete(c.Request.Context(), model.ID(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

func (h *UserHandler) Create(c *gin.Context) {
	var dto dto.CreateUser
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err)
		return
	}

	user, err := h.s.Create(c.Request.Context(), &dto)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, user)
}
