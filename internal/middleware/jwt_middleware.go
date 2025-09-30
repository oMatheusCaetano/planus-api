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

func getJwtClaims(c *gin.Context) (*dto.JWTClaims, *errs.Error) {
    authHeader := c.GetHeader("Authorization")

    if authHeader == "" {
        return nil, errs.New(http.StatusUnauthorized, "Header 'Authorization' não está presente na requisição")
    }

    parts := strings.Split(authHeader, " ")

    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        return nil, errs.New(http.StatusUnauthorized, "Formato inválido para o header 'Authorization'")
    }

    tokenString := parts[1]

    token, err := jwt.ParseWithClaims(tokenString, &dto.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(env.JWTSecret()), nil
    })

    if err != nil || !token.Valid {
        return nil, errs.New(http.StatusUnauthorized, "Token de autenticação inválido ou expirado")
    }

    claims, ok := token.Claims.(*dto.JWTClaims)
    if !ok {
        return nil, errs.New(http.StatusUnauthorized, "Não foi possível extrair as claims do token")
    }

    return claims, nil
}

func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, err := getJwtClaims(c)

        if err != nil {
            response.Error(c, err)
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
