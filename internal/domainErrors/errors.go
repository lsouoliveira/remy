package domainErrors

import (
	"fmt"
)

type DomainError struct {
	Code    string
	Message string
	cause   error
}

func WrapError(err error, code string, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		cause:   err,
	}
}

func (e *DomainError) Unwrap() error {
	return e.cause
}

func (e *DomainError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s - caused by: %v", e.Code, e.Message, e.cause)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Is(target error) bool {
	if targetErr, ok := target.(*DomainError); ok {
		return e.Code == targetErr.Code
	}

	return false
}

func (e *DomainError) As(target any) bool {
	if targetErr, ok := target.(**DomainError); ok {
		*targetErr = e

		return true
	}

	return false
}

func New(code string, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}
