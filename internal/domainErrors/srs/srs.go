package srs

import (
	"remy/internal/domainErrors"
)

var (
	ErrInvalidRepetitions = domainErrors.New("srs.invalid_repetitions", "repetitions must be non-negative")
	ErrInvalidInterval    = domainErrors.New("srs.invalid_interval", "interval must be non-negative")
	ErrInvalidEaseFactor  = domainErrors.New("srs.invalid_ease_factor", "ease factor must be at least 1.3")
	ErrInvalidQuality     = domainErrors.New("srs.invalid_quality", "quality must be between 0 and 5")
)
