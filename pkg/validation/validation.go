package validation

import (
	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func ValidateStruct[T any](input T) error {
	return v.Struct(input)
}
