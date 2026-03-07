package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type InfrastructureError struct {
	Code    string
	Message string
	cause   error
}

var ErrVersionConflict = NewInfrastructureError("infra.version_conflict", "the resource was modified by another request.")

func InvalidPathParameter(param string) *InfrastructureError {
	return NewInfrastructureError("infra.invalid_path_parameter", fmt.Sprintf("invalid path parameter: %s", param))
}

func InvalidCursorParameter() *InfrastructureError {
	return NewInfrastructureError("infra.invalid_cursor_parameter", "invalid cursor parameter")
}

func NewInfrastructureError(code string, message string) *InfrastructureError {
	return &InfrastructureError{
		Code:    code,
		Message: message,
	}
}

func WrapInfrastructureError(err error, code string, message string) *InfrastructureError {
	return &InfrastructureError{
		Code:    code,
		Message: message,
		cause:   err,
	}
}

func (e *InfrastructureError) Unwrap() error {
	return e.cause
}

func (e *InfrastructureError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s - caused by: %v", e.Code, e.Message, e.cause)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *InfrastructureError) Is(target error) bool {
	if targetErr, ok := target.(*InfrastructureError); ok {
		return e.Code == targetErr.Code
	}

	return false
}

func (e *InfrastructureError) As(target any) bool {
	if targetErr, ok := target.(**InfrastructureError); ok {
		*targetErr = e

		return true
	}

	return false
}

type QueryValidationError struct {
	InfrastructureError
	OriginalError validator.ValidationErrors
}

func NewQueryValidationError(originalError validator.ValidationErrors) *QueryValidationError {
	return &QueryValidationError{
		OriginalError: originalError,
	}
}

func (e *QueryValidationError) Error() string {
	return fmt.Sprintf("query validation error: %s", e.OriginalError.Error())
}
