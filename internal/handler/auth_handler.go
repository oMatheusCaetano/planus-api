package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dto"
	"github.com/omatheuscaetano/planus-api/internal/response"
	"github.com/omatheuscaetano/planus-api/internal/service"
)

type AuthHandler struct{
	s *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{s: s}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var dto dto.LoginData
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err)
		return
	}
	res, err := h.s.Login(c.Request.Context(), &dto)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Ok(c, res)
}
