package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
	"github.com/omatheuscaetano/planus-api/internal/model"
	"github.com/omatheuscaetano/planus-api/internal/response"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

func PermissionMiddleware(i *dependency_container.Instances, requiredPermissions ...model.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        userVal, exists := c.Get("user")

        if !exists {
            claims, err := getJwtClaims(c)

            if err != nil {
                response.Error(c, err)
                return
            }

            userVal = &model.User{
                Model: model.Model{
                    ID: model.ID(claims.Sub),
                },
            }
        }

        user, _ := userVal.(*model.User)

        user, err := i.User.Services.User.Find(c.Request.Context(), user.ID)

        if err != nil {
            response.Error(c, err)
            return
        }

        can := user.Can(requiredPermissions...)

        if !can {
            response.Error(c, errs.New(403, "Você nao tem permissão para acessar este recurso"))
            return
        }

        c.Set("user", user)
        c.Next()
    }
}
