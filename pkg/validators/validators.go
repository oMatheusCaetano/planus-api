package validators

import (
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Init() {
    v, ok := binding.Validator.Engine().(*validator.Validate)

    if !ok {
        log.Fatalf("Error getting validator engine")
    }

    v.RegisterValidation("cnpj", cnpj)
    v.RegisterValidation("filled", filled)
}