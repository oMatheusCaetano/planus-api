package response

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/omatheuscaetano/planus-api/pkg/errs"
)

func Abort(c *gin.Context, code int, message string, payload any, meta any) {
	c.AbortWithStatusJSON(code, gin.H{
		"message": message,
		"code":    code,
		"data":    payload,
		"meta":    meta,
	})
}

func JSON(c *gin.Context, code int, payload any) {
	c.JSON(code, payload)
}

func Error(c *gin.Context, err error) {
    apiError := errs.From(err) 
    Abort(c, apiError.Code, apiError.Message, nil, nil)
}

func BadRequest(c *gin.Context, err error) {
	Abort(c, http.StatusBadRequest, "Dados inválidos", translateValidationErrors(err), nil)
}

func Ok(c *gin.Context, payload any) {
	JSON(c, http.StatusOK, payload)
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Created(c *gin.Context, payload any) {
	JSON(c, http.StatusCreated, payload)
}

func translateValidationErrors(err error) map[string]string {
	errors := map[string]string{}

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := mapField(e)

			switch e.Tag() {
			case "required":
				errors[field] = makeMessage(e, "Campo obrigatório")
			case "filled":
				errors[field] = makeMessage(e, "Campo não pode estar vazio")
			case "min":
				errors[field] = makeMessage(e, "Precisa ter no mínimo {param} caracteres")
			case "max":
				errors[field] = makeMessage(e, "Precisa ter no máximo {param} caracteres")
			case "len":
				errors[field] = makeMessage(e, "Precisa ter exatamente {param} caracteres")
			case "unique":
				errors[field] = makeMessage(e, "{field} já existe")
			case "email":
				errors[field] = makeMessage(e, "E-mail inválido")
			case "cnpj":
				errors[field] = makeMessage(e, "CNPJ inválido")
			default:
				errors[field] = makeMessage(e, "Valor inválido")
			}
		}
	}

	return errors
}

func mapField(e validator.FieldError) string {
	field := toSnakeCase(e.Field())
	return field
}

func makeMessage(e validator.FieldError, msg string) string {
	msg = strings.ReplaceAll(msg, "{field}", e.Field())
	msg = strings.ReplaceAll(msg, "{param}", e.Param())
	return msg
}

func toSnakeCase(s string) string {
	regex := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := regex.ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}
