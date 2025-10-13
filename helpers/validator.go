package helpers

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates any struct and returns an error
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
