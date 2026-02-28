package models

import (
	"time"

	"remy/internal/domainErrors"
)

type SRSState struct {
	Repetitions int       `gorm:"type:int"`
	Interval    int       `gorm:"type:int"`
	EaseFactor  float64   `gorm:"type:decimal(3,2)"`
	ReviewAt    time.Time `gorm:"type:datetime"`
}

func NewSRSState(repetitions int, interval int, easeFactor float64, reviewAt time.Time) (*SRSState, error) {
	if repetitions < 0 {
		return nil, domainErrors.SRSState.InvalidRepetitions
	}

	if interval < 0 {
		return nil, domainErrors.SRSState.InvalidInterval
	}

	if easeFactor < 1.3 {
		return nil, domainErrors.SRSState.InvalidEaseFactor
	}

	return &SRSState{
		Repetitions: repetitions,
		Interval:    interval,
		EaseFactor:  easeFactor,
		ReviewAt:    reviewAt,
	}, nil
}
