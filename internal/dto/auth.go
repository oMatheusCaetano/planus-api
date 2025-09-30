package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/omatheuscaetano/planus-api/internal/model"
)

type LoginData struct {
    Username string `json:"username" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Token     string       `json:"token"`
    ExpiresIn int64        `json:"expires_in"`
    User      *model.User  `json:"user"`
}

type JWTClaims struct {
    Sub string `json:"sub"`
    jwt.RegisteredClaims
}
