package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// filled checks if a string field is nil or empty (after trimming spaces).
// Can replace the "required" validator for string fields.
func filled(fl validator.FieldLevel) bool {
	val := fl.Field().String()
	return strings.TrimSpace(val) != ""
}
