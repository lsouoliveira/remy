package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type InfrastructureError struct {
	Code    string
	Message string
}

type QueryValidationError struct {
	InfrastructureError
	OriginalError validator.ValidationErrors
}

func NewInfrastructureError(code string, message string) *InfrastructureError {
	return &InfrastructureError{
		Code:    code,
		Message: message,
	}
}

func (e *InfrastructureError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewQueryValidationError(originalError validator.ValidationErrors) *QueryValidationError {
	return &QueryValidationError{
		OriginalError: originalError,
	}
}

func (e *QueryValidationError) Error() string {
	return fmt.Sprintf("query validation error: %s", e.OriginalError.Error())
}
