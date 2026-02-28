package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

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

func (s *NoteServiceTestSuite) TestCreateNote_WhenParamsAreValid_CreatesNote() {
	mockPublisher := new(mocks.MockDomainEventPublisher)
	noteService := services.NewNoteService(s.DB, mockPublisher)

	mockPublisher.On("Publish", mock.Anything).Return(nil)

	note, err := noteService.Create(services.NoteCreate{Content: "Test Note"})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), note)
	assert.Equal(s.T(), "Test Note", note.Content)
}

func (s *NoteServiceTestSuite) TestCreateNote_WhenParamsAreValid_PublishesEvent() {
	mockPublisher := new(mocks.MockDomainEventPublisher)
	noteService := services.NewNoteService(s.DB, mockPublisher)

	mockPublisher.On("Publish", models.NoteCreatedEvent{
		NoteID: uint(1),
	}).Return(nil)

	note, err := noteService.Create(services.NoteCreate{Content: "Test Note"})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), note)
	mockPublisher.AssertCalled(s.T(), "Publish", models.NoteCreatedEvent{
		NoteID: note.ID,
	})
}

func (s *NoteServiceTestSuite) TestCreateNote_WhenPublisherReturnsError_ReturnsError() {
	mockPublisher := new(mocks.MockDomainEventPublisher)

	noteService := services.NewNoteService(s.DB, mockPublisher)

	mockPublisher.On("Publish", mock.Anything).Return(assert.AnError)

	note, err := noteService.Create(services.NoteCreate{Content: "Test Note"})

	assert.Error(s.T(), err)
	assert.Nil(s.T(), note)
}

func (s *NoteServiceTestSuite) TestCreateNote_WhenPublisherFails_RollsBackTransaction() {
	mockPublisher := new(mocks.MockDomainEventPublisher)

	noteService := services.NewNoteService(s.DB, mockPublisher)

	mockPublisher.On("Publish", mock.Anything).Return(assert.AnError)

	note, err := noteService.Create(services.NoteCreate{Content: "Test Note"})

	assert.Error(s.T(), err)
	assert.Nil(s.T(), note)

	var count int64

	s.DB.Model(&models.Note{}).Count(&count)

	assert.Equal(s.T(), int64(0), count)
}

func (s *NoteServiceTestSuite) TestListNotes_WhenNotesExist_ReturnsEmptyList() {
	noteService := services.NewNoteService(s.DB, &mocks.MockDomainEventPublisher{})

	notes, err := noteService.List(services.ListNotesParams{
		Page:     1,
		PageSize: 10,
		SortBy:   "created_at",
		Order:    "asc",
	})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notes)
	assert.Equal(s.T(), 0, len(notes.Notes))
	assert.Equal(s.T(), int64(0), notes.Total)
}

func (s *NoteServiceTestSuite) TestListNotes_WhenNotesExist_ReturnsNotes() {
	noteService := services.NewNoteService(s.DB, &mocks.MockDomainEventPublisher{})

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
	noteService := services.NewNoteService(s.DB, &mocks.MockDomainEventPublisher{})

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
	noteService := services.NewNoteService(s.DB, &mocks.MockDomainEventPublisher{})

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
