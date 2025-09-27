package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/response"
	"github.com/omatheuscaetano/planus-api/pkg/env"
)

type AppHandler struct{}

func NewAppHandler() *AppHandler {
	return &AppHandler{}
}

func (h *AppHandler) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, response.WelcomeResponse{
		Message:   "Bem-vindo à API " + env.AppName() + "!",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
