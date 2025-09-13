package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/pkg/app"
)

func AppContextMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        appCtx := &app.Context{
            Context:  c.Request.Context(),
        }

        c.Request = c.Request.WithContext(appCtx)
        c.Next()
    }
}
