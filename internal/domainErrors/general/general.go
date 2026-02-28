package general

import (
	"fmt"

	"remy/internal/domainErrors"
)

var ErrNotFound = domainErrors.New("general.not_found", "the requested resource was not found.")

func NotFound(entity string, id any) *domainErrors.DomainError {
	return domainErrors.New("general.not_found", entity+" with ID "+fmt.Sprintf("%v", id)+" not found")
}
