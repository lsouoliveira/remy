package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"remy/internal/testhelpers"
)

type NoteServiceTestSuite struct {
	testhelpers.TestSuite
}

func TestNoteServiceTestSuite(t *testing.T) {
	suite.Run(t, new(NoteServiceTestSuite))
}

func (s *NoteServiceTestSuite) TestCreateNote_WhenParamsAreValid_CreatesNote() {
	noteService := NewNoteService(s.DB)

	note, err := noteService.Create(NoteCreate{Content: "Test Note"})

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), note)
	assert.Equal(s.T(), "Test Note", note.Content)
}

func (s *NoteServiceTestSuite) TestListNotes_WhenNotesExist_ReturnsEmptyList() {
	noteService := NewNoteService(s.DB)

	notes, err := noteService.List(NewListNotesRequest(1, 10))

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notes)
	assert.Equal(s.T(), 0, len(notes.Notes))
	assert.Equal(s.T(), int64(0), notes.Total)
}

func (s *NoteServiceTestSuite) TestListNotes_WhenNotesExist_ReturnsNotes() {
	noteService := NewNoteService(s.DB)

	_, _ = noteService.Create(NoteCreate{Content: "Test Note 1"})
	_, _ = noteService.Create(NoteCreate{Content: "Test Note 2"})

	notes, err := noteService.List(NewListNotesRequest(1, 10))

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notes)
	assert.Equal(s.T(), 2, len(notes.Notes))
	assert.Equal(s.T(), int64(2), notes.Total)
	assert.Equal(s.T(), "Test Note 1", notes.Notes[0].Content)
	assert.Equal(s.T(), "Test Note 2", notes.Notes[1].Content)
}

func (s *NoteServiceTestSuite) TestListNotes_WhenPageSizeIsLessThanOne_ReturnsDefaultPageSize() {
	noteService := NewNoteService(s.DB)

	_, _ = noteService.Create(NoteCreate{Content: "Test Note 1"})
	_, _ = noteService.Create(NoteCreate{Content: "Test Note 2"})

	notes, err := noteService.List(NewListNotesRequest(1, 0))

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notes)
	assert.Equal(s.T(), 2, len(notes.Notes))
	assert.Equal(s.T(), int64(2), notes.Total)
}

func (s *NoteServiceTestSuite) TestListNotes_WhenPageIsLessThanOne_ReturnsFirstPage() {
	noteService := NewNoteService(s.DB)

	_, _ = noteService.Create(NoteCreate{Content: "Test Note 1"})
	_, _ = noteService.Create(NoteCreate{Content: "Test Note 2"})

	notes, err := noteService.List(NewListNotesRequest(0, 10))

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), notes)
	assert.Equal(s.T(), 2, len(notes.Notes))
	assert.Equal(s.T(), int64(2), notes.Total)
}
