package validators

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	db "github.com/omatheuscaetano/planus-api/internal/database"
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
    sqlDB := db.GetDB()
    var exists bool
    query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE %s = $1 LIMIT 1)`, tableName, fieldName)
    sqlDB.QueryRow(query, value).Scan(&exists)
    return !exists
}
