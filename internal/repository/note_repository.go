package repository

import (
	"github.com/google/uuid"

	"remy/internal/models"
)

type ListParams struct {
	SortBy models.SortField
	Order  models.SortOrder
	Limit  int
	Cursor *models.Cursor
}

type NoteRepository interface {
	Save(note *models.Note) error
	GetByID(id uuid.UUID) (*models.Note, error)
	List(params ListParams) ([]*models.Note, int64, error)
}
