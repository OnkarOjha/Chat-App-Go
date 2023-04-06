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

// func MapValidation(mp map[string]interface{}, rules map[string]interface{}, w http.ResponseWriter) {
// 	errs := Validate.ValidateMap(mp, rules)
// 	if len(errs) > 0 {
// 		response.ShowResponse("Failure", 400,"", errs, w)
// 		return
// 	}
// }
