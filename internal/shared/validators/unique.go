package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/omatheuscaetano/planus-api/internal/db"
)

var unique validator.Func = func(fl validator.FieldLevel) bool {
    // O parâmetro vem no formato "table.field"
    param := fl.Param() // ex: "table.column"
    parts := strings.Split(param, ".")
    if len(parts) != 2 {
        return false // parâmetros inválidos
    }
    tableName := parts[0]
    fieldName := parts[1]
    value := fl.Field().Interface()

    return Unique(tableName, fieldName, value)
}

func Unique(tableName string, fieldName string, value interface{}) bool {
    exists, _ := db.
        From(tableName).
        Where(fieldName, "=", value).
        Exists()
    return !exists
}
