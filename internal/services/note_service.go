package services

import (
	"fmt"

	"gorm.io/gorm"

	"remy/internal/models"
)

type NoteService struct {
	db        *gorm.DB
	publisher models.DomainEventPublisher
}

type NoteCreate struct {
	Content string `json:"content" binding:"required"`
}

type NoteRead struct {
	ID        uint   `json:"id"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type NoteList struct {
	Notes []NoteRead
	Total int64
}

type ListNotesParams struct {
	Page     int
	PageSize int
}

func NewListNotesRequest(page int, pageSize int) ListNotesParams {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	return ListNotesParams{
		Page:     page,
		PageSize: pageSize,
	}
}

func NewNoteService(db *gorm.DB, publisher models.DomainEventPublisher) *NoteService {
	return &NoteService{db: db, publisher: publisher}
}

func (s *NoteService) Create(request NoteCreate) (*NoteRead, error) {
	tx := s.db.Begin()

	note := models.CreateNote(request.Content)

	if err := tx.Create(note).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	if err := s.publisher.Publish(models.NoteCreatedEvent{NoteID: note.ID}); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to publish note created event: %w", err)
	}

	tx.Commit()

	return mapToNoteRead(note), nil
}

func (s *NoteService) List(params ListNotesParams) (*NoteList, error) {
	var total int64
	if err := s.db.Model(&models.Note{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count notes: %w", err)
	}

	var notes []models.Note
	if err := s.db.Order("created_at desc").
		Limit(params.PageSize).
		Offset((params.Page - 1) * params.PageSize).
		Find(&notes).Error; err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	noteReads := make([]NoteRead, len(notes))
	for i, note := range notes {
		noteReads[i] = *mapToNoteRead(&note)
	}

	return &NoteList{
		Notes: noteReads,
		Total: total,
	}, nil
}

func mapToNoteRead(note *models.Note) *NoteRead {
	return &NoteRead{
		ID:        note.ID,
		Content:   note.Content,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
