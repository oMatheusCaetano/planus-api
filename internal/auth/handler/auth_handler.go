package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/auth/dto"
	"github.com/omatheuscaetano/planus-api/internal/auth/service"
	"github.com/omatheuscaetano/planus-api/pkg/responses"
)

type AuthHandler struct {
    service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
    return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
    var dto dto.Login
    if err := c.ShouldBindJSON(&dto); err != nil {
        responses.BadRequest(c, err)
        return
    }

    response, err := h.service.Login(c.Request.Context(), &dto)
    if err != nil {
        responses.Error(c, err)
        return
    }

    responses.Ok(c, response)
}
