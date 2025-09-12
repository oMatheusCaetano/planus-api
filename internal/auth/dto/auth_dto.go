package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/omatheuscaetano/planus-api/internal/person/model"
)

type CreateUser struct {
    ID       int    `json:"id".       binding:"required,min=1"`
    Email    string `json:"email"     binding:"required,email,max=255"`
    Password string `json:"password"  binding:"required,min=8,max=100"`
}

type Login struct {
    Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginData struct {
    Token     string         `json:"token"`
    ExpiresIn int64          `json:"expires_in"`
    User      *model.Person  `json:"user"`
}

type JWTClaims struct {
    Sub int `json:"sub"`
    jwt.RegisteredClaims
}
