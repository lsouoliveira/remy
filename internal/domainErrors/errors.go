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
	InvalidRepetitions: NewError("srs_state.invalid_repetitions", "Repetitions must be non-negative"),
	InvalidInterval:    NewError("srs_state.invalid_interval", "Interval must be non-negative"),
	InvalidEaseFactor:  NewError("srs_state.invalid_ease_factor", "Ease factor must be at least 1.3"),
	InvalidQuality:     NewError("srs_state.invalid_quality", "Quality must be between 0 and 5"),
}

var NotFoundError = func(entity string, id any) *DomainError {
	return NewError("not_found", fmt.Sprintf("%s with ID %v not found", entity, id))
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
