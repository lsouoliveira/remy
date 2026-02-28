package domainErrors

import (
	"fmt"
)

type DomainError struct {
	Status  int
	Code    string
	Message string
}

func NewError(code string, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
