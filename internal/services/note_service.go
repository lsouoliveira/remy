package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"remy/internal/domainErrors/general"
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
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	ReviewAt  time.Time `json:"review_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteList struct {
	Notes []NoteRead
	Total int64
}

type ListNotesParams struct {
	Page     int
	PageSize int
	SortBy   string
	Order    string
}

type ReviewParams struct {
	NoteID  uint
	Quality int
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

func (s *NoteService) Review(reviewParams ReviewParams) error {
	tx := s.db.Begin()

	var note models.Note
	if err := tx.First(&note, reviewParams.NoteID).Error; err != nil {
		tx.Rollback()
		return general.NotFound("note", reviewParams.NoteID)
	}

	if err := note.Review(reviewParams.Quality, models.NewSM2Algorithm()); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(&note).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save reviewed note: %w", err)
	}

	tx.Commit()

	return nil
}

func (s *NoteService) List(params ListNotesParams) (*NoteList, error) {
	var total int64
	if err := s.db.Model(&models.Note{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count notes: %w", err)
	}

	var notes []models.Note
	if err := s.db.Order(fmt.Sprintf("%s %s", params.SortBy, params.Order)).
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
		ReviewAt:  note.SRSState.ReviewAt,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}
