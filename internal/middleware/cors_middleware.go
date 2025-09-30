package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
    config := cors.DefaultConfig()

    origin := os.Getenv("CORS_ORIGINS")
    if origin != "" {
        config.AllowOrigins = strings.Split(origin, ",")
    }

    methods := os.Getenv("CORS_METHODS")
    if methods != "" {
        config.AllowMethods = strings.Split(methods, ",")
    }

    headers := os.Getenv("CORS_HEADERS")
    if headers != "" {
        config.AllowHeaders = strings.Split(headers, ",")
    }

    credentials := os.Getenv("CORS_CREDENTIALS")
    config.AllowCredentials = strings.ToLower(credentials) == "true"

    return cors.New(config)
}
