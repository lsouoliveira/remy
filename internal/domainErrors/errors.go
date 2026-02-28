package domainErrors

import (
	"fmt"
)

type DomainError struct {
	Code    string
	Message string
}

var SRSState = struct {
	InvalidRepetitions *DomainError
	InvalidInterval    *DomainError
	InvalidEaseFactor  *DomainError
	InvalidQuality     *DomainError
}{
	InvalidRepetitions: NewError("invalid_repetitions", "Repetitions must be non-negative"),
	InvalidInterval:    NewError("invalid_interval", "Interval must be non-negative"),
	InvalidEaseFactor:  NewError("invalid_ease_factor", "Ease factor must be at least 1.3"),
	InvalidQuality:     NewError("invalid_quality", "Quality must be between 0 and 5"),
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
