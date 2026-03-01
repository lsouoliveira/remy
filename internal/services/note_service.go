package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"remy/internal/models"
	"remy/internal/repository"
	"remy/internal/unit_of_work"
)

type NoteService struct {
	repo       repository.NoteRepository
	uowFactory unit_of_work.UnitOfWorkFactory
}

type NoteCreate struct {
	Content string `json:"content" binding:"required"`
}

type NoteRead struct {
	ID        uuid.UUID `json:"id"`
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
	NoteID  uuid.UUID
	Quality int
}

func NewNoteService(repo repository.NoteRepository, uowFactory unit_of_work.UnitOfWorkFactory) *NoteService {
	return &NoteService{repo: repo, uowFactory: uowFactory}
}

func (s *NoteService) Create(request NoteCreate) (*NoteRead, error) {
	uow := s.uowFactory.New()
	defer uow.Rollback()

	note := models.CreateNote(request.Content)

	if err := uow.Notes().Save(note); err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	if err := uow.Commit(); err != nil {
		return nil, err
	}

	return mapToNoteRead(note), nil
}

func (s *NoteService) Review(reviewParams ReviewParams) error {
	uow := s.uowFactory.New()
	defer uow.Rollback()

	notes := uow.Notes()

	note, err := notes.GetByID(reviewParams.NoteID)
	if err != nil {
		return err
	}

	if err := note.Review(reviewParams.Quality, models.NewSM2Algorithm()); err != nil {
		return err
	}

	if err := notes.Save(note); err != nil {
		return err
	}

	return uow.Commit()
}

func (s *NoteService) List(params ListNotesParams) (*NoteList, error) {
	notes, total, err := s.repo.List(repository.ListParams{
		Page:     params.Page,
		PageSize: params.PageSize,
		SortBy:   params.SortBy,
		Order:    params.Order,
	})
	if err != nil {
		return nil, err
	}

	noteReads := make([]NoteRead, len(notes))
	for i, note := range notes {
		noteReads[i] = *mapToNoteRead(note)
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
