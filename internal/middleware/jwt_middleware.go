package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/omatheuscaetano/planus-api/internal/dto"
	"github.com/omatheuscaetano/planus-api/internal/model"
	"github.com/omatheuscaetano/planus-api/internal/response"
	"github.com/omatheuscaetano/planus-api/pkg/env"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")

        if authHeader == "" {
            response.Error(c, errs.New(http.StatusUnauthorized, "Header 'Authorization' não está presente na requisição"))
            return
        }

        parts := strings.Split(authHeader, " ")

        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            response.Error(c, errs.New(http.StatusUnauthorized, "Formato inválido para o header 'Authorization'"))
            return
        }

        tokenString := parts[1]

        token, err := jwt.ParseWithClaims(tokenString, &dto.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(env.JWTSecret()), nil
        })

        if err != nil || !token.Valid {
            response.Error(c, errs.New(http.StatusUnauthorized, "Token de autenticação inválido ou expirado"))
            return
        }

        claims, ok := token.Claims.(*dto.JWTClaims)
        if !ok {
            response.Error(c, errs.New(http.StatusUnauthorized, "Não foi possível extrair as claims do token"))
            return
        }

        c.Set("user", &model.User{
            Model: model.Model{
                ID: model.ID(claims.Sub),
            },
        })

        c.Next()
    }
}
