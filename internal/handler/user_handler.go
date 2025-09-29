package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/model"
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
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
