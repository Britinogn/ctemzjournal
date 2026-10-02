package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct runs go-playground/validator on a struct.
func ValidateStruct(v any) error {
	return validate.Struct(v)
}
