package identity

import (
	"flomart/internal/identity/dto"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRegisterInput(input dto.RegisterInput) error {
	return validate.Struct(input)
}
