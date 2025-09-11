package validators

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var cnpj validator.Func = func(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    return Cnpj(value)
}

func Cnpj(value string) bool {
	// Remove tudo que não é número
	re := regexp.MustCompile(`\D`)
	cnpj := re.ReplaceAllString(value, "")

	if len(cnpj) != 14 {
		return false
	}

	// Verifica se todos os dígitos são iguais (inválido)
	for i := 0; i < 14; i++ {
		if strings.Count(cnpj, string(cnpj[i])) == 14 {
			return false
		}
	}

	// Calcula os dígitos verificadores
	dv1, dv2 := calcCnpjDigits(cnpj[:12])
	return cnpj[12:] == dv1+dv2
}

func calcCnpjDigits(numbers string) (string, string) {
	multipliers1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	multipliers2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	sum := 0
	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(numbers[i]))
		sum += digit * multipliers1[i]
	}
	dv1 := 0
	if mod := sum % 11; mod < 2 {
		dv1 = 0
	} else {
		dv1 = 11 - mod
	}

	sum = 0
	numbers += strconv.Itoa(dv1)
	for i := 0; i < 13; i++ {
		digit, _ := strconv.Atoi(string(numbers[i]))
		sum += digit * multipliers2[i]
	}
	dv2 := 0
	if mod := sum % 11; mod < 2 {
		dv2 = 0
	} else {
		dv2 = 11 - mod
	}

	return strconv.Itoa(dv1), strconv.Itoa(dv2)
}