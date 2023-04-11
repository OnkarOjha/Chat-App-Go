package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func CheckValidation(data interface{}) error {
	fmt.Println("dadtaaskh: ", data)
	validationErr := Validate.Var(data,"required")
	if validationErr != nil {
		return validationErr
	}
	return nil
}

func CheckValidationStruct(data interface{}) error {
	fmt.Println("dadtaaskh: ", data)

	validationErr := Validate.Var(data,"required")
	if validationErr != nil {
		return validationErr
	}
	return nil
}
