package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/auth/dto"
	"github.com/omatheuscaetano/planus-api/pkg/app"
	"github.com/omatheuscaetano/planus-api/pkg/db"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
	"github.com/omatheuscaetano/planus-api/pkg/responses"
)

func AuthorizeMiddleware(actions []dto.P) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.Request.Context().(*app.Context).UserID

        if userID == nil {
            responses.Error(c, errs.New(http.StatusUnauthorized, "Não foi possível identificar o usuário autenticado"))
            return
        }

        db := db.GetDB()

        query := `SELECT TRUE FROM permissions WHERE user_id = ? AND (`

        args := []interface{}{*userID}

        for i, action := range actions {
            operator := "OR"

            if (i == 0) {
                operator = ""
            }

            query += fmt.Sprintf(" %s (module = ? AND action = ?)", operator)
            args = append(args, action.Module, action.Action)
        }

        query += ") LIMIT 1"

        var hasPermission bool

        err := db.QueryRow(query, args...).Scan(&hasPermission)

        if err != nil || !hasPermission {
            responses.Error(c, errs.New(http.StatusForbidden, "Usuário não possui permissão para acessar este recurso"))
            return
        }

        c.Next()
    }
}
