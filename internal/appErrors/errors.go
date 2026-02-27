package appErrors

import (
	"fmt"
)

type AppError struct {
	Status  int
	Code    string
	Message string
}

func NewError(status int, code string, message string) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
