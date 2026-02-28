package models

import (
	"time"

	"remy/internal/domainErrors"
)

type SRSAlgorithm interface {
	CalculateNextReview(srsState *SRSState, quality int) (*SRSState, error)
}

type SM2Algorithm struct{}

type Note struct {
	ID        uint      `gorm:"primaryKey"`
	Content   string    `gorm:"type:text;not null"`
	SRSState  SRSState  `gorm:"embedded"`
	Version   uint      `gorm:"type:int;default:1"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateNote(content string) *Note {
	srsState, _ := NewSRSState(0, 0, 2.5, time.Now())

	return &Note{
		SRSState: *srsState,
		Content:  content,
		Version:  1,
	}
}

func (n *Note) Review(quality int, algorithm SRSAlgorithm) error {
	newSRSState, err := algorithm.CalculateNextReview(&n.SRSState, quality)
	if err != nil {
		return err
	}

	n.SRSState = *newSRSState
	n.Version += 1

	return nil
}

func NewSM2Algorithm() *SM2Algorithm {
	return &SM2Algorithm{}
}

func (a *SM2Algorithm) CalculateNextReview(srsState *SRSState, quality int) (*SRSState, error) {
	if quality < 0 || quality > 5 {
		return nil, domainErrors.SRSState.InvalidQuality
	}

	var newRepetitions int
	var newInterval int
	var newEaseFactor float64

	if quality < 3 {
		newRepetitions = 0
		newInterval = 1
	} else {
		newRepetitions = srsState.Repetitions + 1
		newInterval = intervalForRepetitions(newRepetitions, srsState.EaseFactor, srsState.Interval)
	}

	newEaseFactor = srsState.EaseFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))

	if newEaseFactor < 1.3 {
		newEaseFactor = 1.3
	}

	return NewSRSState(newRepetitions, newInterval, newEaseFactor, time.Now().AddDate(0, 0, newInterval))
}

func intervalForRepetitions(repetitions int, easeFactor float64, interval int) int {
	switch repetitions {
	case 1:
		return 1
	case 2:
		return 6
	default:
		return int(float64(interval) * easeFactor)
	}
}
