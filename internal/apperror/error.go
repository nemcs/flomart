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

/*
var (
	ErrNotFound = NewAppError(nil, "not found", "", "US-0003")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developerMessage,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage, code string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}


*/
