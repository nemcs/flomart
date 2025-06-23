package apperror

import (
	"fmt"
)

type AppError struct {
	UserMsg string
	DevMsg  string
	Err     error
	Code    int
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s : %s", e.DevMsg, e.Err.Error())
	}
	return e.DevMsg
}

// Позволяет использовать errors.Is() и errors.As() с AppError, чтобы доставать вложенные ошибки.
func (e *AppError) Unwrap() error {
	return e.Err
}
func New(err error, userMsg, devMsg string, code int) *AppError {
	return &AppError{
		UserMsg: userMsg,
		DevMsg:  devMsg,
		Err:     err,
		Code:    code,
	}
}
func Wrap(err error, userMsg, devMsg string, code int) *AppError {
	return &AppError{
		UserMsg: userMsg,
		DevMsg:  devMsg,
		Err:     err,
		Code:    code,
	}
}
