package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/model"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// @Summary Find User
// @Description Returns a user by ID
// @Tags User
// @Success 200 {object} model.User
// @Router /user/{id} [get]
func (h *UserHandler) Find(c *gin.Context) {
	c.JSON(http.StatusOK, model.User{
		ID:        "1",
		Name:      "John Doe",
		Username:  "johndoe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}
