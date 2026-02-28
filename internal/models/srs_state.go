package models

import (
	"time"

	"remy/internal/domainErrors/srs"
)

type SRSState struct {
	Repetitions int       `gorm:"type:int"`
	Interval    int       `gorm:"type:int"`
	EaseFactor  float64   `gorm:"type:decimal(3,2)"`
	ReviewAt    time.Time `gorm:"type:datetime"`
}

func NewSRSState(repetitions int, interval int, easeFactor float64, reviewAt time.Time) (*SRSState, error) {
	if repetitions < 0 {
		return nil, srs.ErrInvalidRepetitions
	}

	if interval < 0 {
		return nil, srs.ErrInvalidInterval
	}

	if easeFactor < 1.3 {
		return nil, srs.ErrInvalidEaseFactor
	}

	return &SRSState{
		Repetitions: repetitions,
		Interval:    interval,
		EaseFactor:  easeFactor,
		ReviewAt:    reviewAt,
	}, nil
}
