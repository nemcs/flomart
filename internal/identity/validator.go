package identity

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateRegisterInput(input RegisterInput) error {
	return validate.Struct(input)
}
