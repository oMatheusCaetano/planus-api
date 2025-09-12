package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/omatheuscaetano/planus-api/internal/auth/dto"
	"github.com/omatheuscaetano/planus-api/pkg/app"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"github.com/omatheuscaetano/planus-api/pkg/responses"
)

func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")

        if authHeader == "" {
            responses.Error(c, errs.New(http.StatusUnauthorized, "Header 'Authorization' não está presente na requisição"))
            return
        }

        parts := strings.Split(authHeader, " ")

        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            responses.Error(c, errs.New(http.StatusUnauthorized, "Formato inválido para o header 'Authorization'"))
            return
        }

        tokenString := parts[1]

        token, err := jwt.ParseWithClaims(tokenString, &dto.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(app.JWTSecret()), nil
        })
        
        if err != nil || !token.Valid {
            responses.Error(c, errs.New(http.StatusUnauthorized, "Token de autenticação inválido ou expirado"))
            return
        }

        claims, ok := token.Claims.(*dto.JWTClaims)
        if !ok {
            responses.Error(c, errs.New(http.StatusUnauthorized, "Não foi possível extrair as claims do token"))
            return
        }

        if appCtx, ok := c.Request.Context().(*app.AppContext); ok {
            appCtx.UserID = &claims.Sub
        }

        c.Next()
    }
}