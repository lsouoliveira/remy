package repository

import (
	"fmt"

	gormpkg "gorm.io/gorm"

	"remy/internal/domainErrors/general"
	infraErrors "remy/internal/infrastructure/errors"
	"remy/internal/models"
	"remy/internal/repository"
)

type NoteRepository struct {
	db        *gormpkg.DB
	publisher models.DomainEventPublisher
}

func NewNoteRepository(db *gormpkg.DB, publisher models.DomainEventPublisher) repository.NoteRepository {
	return &NoteRepository{db: db, publisher: publisher}
}

func (r *NoteRepository) Save(note *models.Note) error {
	if note.ID == 0 {
		if err := r.db.Create(note).Error; err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}
	} else {
		result := r.db.Model(note).Where("version = ?", note.Version-1).Save(note)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return infraErrors.ErrVersionConflict
		}
	}

	for _, event := range note.FlushEvents() {
		if err := r.publisher.Publish(event); err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}

	return nil
}

func (r *NoteRepository) GetByID(id uint) (*models.Note, error) {
	var note models.Note
	if err := r.db.First(&note, id).Error; err != nil {
		return nil, general.NotFound("note", id)
	}

	return &note, nil
}

func (r *NoteRepository) List(params repository.ListParams) ([]*models.Note, int64, error) {
	var total int64
	if err := r.db.Model(&models.Note{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count notes: %w", err)
	}

	var notes []*models.Note
	if err := r.db.Order(fmt.Sprintf("%s %s", params.SortBy, params.Order)).
		Limit(params.PageSize).
		Offset((params.Page - 1) * params.PageSize).
		Find(&notes).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list notes: %w", err)
	}

	return notes, total, nil
}
