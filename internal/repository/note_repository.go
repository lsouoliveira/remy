package repository

import (
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
	GetByID(id uint) (*models.Note, error)
	List(params ListParams) ([]*models.Note, int64, error)
}
