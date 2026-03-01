package repository

import (
	"github.com/google/uuid"

	"remy/internal/models"
)

type ListParams struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

type NoteRepository interface {
	Save(note *models.Note) error
	GetByID(id uuid.UUID) (*models.Note, error)
	List(params ListParams) ([]*models.Note, int64, error)
}
