package integration

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"remy/internal/domainErrors/general"
	"remy/internal/infrastructure/db"
	infraRepo "remy/internal/infrastructure/db/repository"
	"remy/internal/mocks"
	"remy/internal/models"
	"remy/internal/services"
	"remy/internal/testhelpers"
)

type NoteServiceTestSuite struct {
	testhelpers.IntegrationSuite
}

func TestNoteServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NoteServiceTestSuite))
}

func (s *NoteServiceTestSuite) SetupSuite() {
	s.IntegrationSuite.SetupSuite()
}

func (s *NoteServiceTestSuite) newNoteService(publisher models.DomainEventPublisher) *services.NoteService {
	noteRepo := infraRepo.NewNoteRepository(s.DB, publisher)
	uowFactory := db.NewGormUnitOfWorkFactory(s.DB, publisher)
	return services.NewNoteService(noteRepo, uowFactory)
}

func (s *NoteServiceTestSuite) TestCreateNote_WhenParamsAreValid_CreatesNote() {
	mockPublisher := new(mocks.MockDomainEventPublisher)
	noteService := s.newNoteService(mockPublisher)

	mockPublisher.On("Publish", mock.Anything).Return(nil)

	note, err := noteService.Create(services.NoteCreate{Content: "Test Note"})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), note)
	assert.Equal(s.T(), "Test Note", note.Content)
}


func (s *NoteServiceTestSuite) TestListNotes_WhenNotesExist_ReturnsEmptyList() {
	noteService := s.newNoteService(&mocks.MockDomainEventPublisher{})

	notes, err := noteService.List(services.ListNotesParams{
		Page:     1,
		PageSize: 10,
		SortBy:   "created_at", Order: "asc",
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notes)
	assert.Equal(s.T(), 0, len(notes.Notes))
	assert.Equal(s.T(), int64(0), notes.Total)
}

func (s *NoteServiceTestSuite) TestListNotes_WhenNotesExist_ReturnsNotes() {
	noteService := s.newNoteService(&mocks.MockDomainEventPublisher{})

	s.DB.Create(&models.Note{Content: "Test Note 1"})
	s.DB.Create(&models.Note{Content: "Test Note 2"})

	notes, err := noteService.List(services.ListNotesParams{
		Page:     1,
		PageSize: 10,
		SortBy:   "created_at",
		Order:    "asc",
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notes)
	assert.Equal(s.T(), 2, len(notes.Notes))
	assert.Equal(s.T(), int64(2), notes.Total)
	assert.Contains(s.T(), []string{notes.Notes[0].Content, notes.Notes[1].Content}, "Test Note 1")
	assert.Contains(s.T(), []string{notes.Notes[0].Content, notes.Notes[1].Content}, "Test Note 2")
}

func (s *NoteServiceTestSuite) TestListNotes_WhenPageSizeIsLessThanOne_ReturnsDefaultPageSize() {
	noteService := s.newNoteService(&mocks.MockDomainEventPublisher{})

	s.DB.Create(&models.Note{Content: "Test Note 1"})
	s.DB.Create(&models.Note{Content: "Test Note 2"})

	notes, err := noteService.List(services.ListNotesParams{
		Page:     1,
		PageSize: 0,
		SortBy:   "created_at",
		Order:    "asc",
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notes)
	assert.Equal(s.T(), 0, len(notes.Notes))
	assert.Equal(s.T(), int64(2), notes.Total)
}

func (s *NoteServiceTestSuite) TestListNotes_WhenPageIsLessThanOne_ReturnsFirstPage() {
	noteService := s.newNoteService(&mocks.MockDomainEventPublisher{})

	s.DB.Create(&models.Note{Content: "Test Note 1"})
	s.DB.Create(&models.Note{Content: "Test Note 2"})

	notes, err := noteService.List(services.ListNotesParams{
		Page:     0,
		PageSize: 10,
		SortBy:   "created_at",
		Order:    "asc",
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notes)
	assert.Equal(s.T(), 2, len(notes.Notes))
	assert.Equal(s.T(), int64(2), notes.Total)
}

func (s *NoteServiceTestSuite) TestReview_WhenNoteExists_UpdatesSRSState() {
	noteService := s.newNoteService(&mocks.MockDomainEventPublisher{})

	note := &models.Note{
		Content: "Test Note",
		SRSState: models.SRSState{
			Repetitions: 0,
			Interval:    0,
			EaseFactor:  2.5,
			ReviewAt:    time.Now(),
		},
	}

	s.DB.Create(note)

	err := noteService.Review(services.ReviewParams{
		NoteID:  note.ID,
		Quality: 4,
	})

	assert.NoError(s.T(), err)

	var updatedNote models.Note
	s.DB.First(&updatedNote, note.ID)

	assert.Equal(s.T(), 1, updatedNote.SRSState.Repetitions)
	assert.Equal(s.T(), 1, updatedNote.SRSState.Interval)
	assert.InDelta(s.T(), 2.5, updatedNote.SRSState.EaseFactor, 0.01)
}

func (s *NoteServiceTestSuite) TestReview_WhenNoteDoesNotExist_ReturnsError() {
	noteService := s.newNoteService(&mocks.MockDomainEventPublisher{})

	err := noteService.Review(services.ReviewParams{
		NoteID:  999,
		Quality: 4,
	})

	assert.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, general.ErrNotFound))
}
