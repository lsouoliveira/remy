package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"remy/internal/response"
	"remy/internal/services"
	"remy/internal/testhelpers"
)

type NoteHandlerTestSuite struct {
	testhelpers.IntegrationSuite
}

func TestNoteHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(NoteHandlerTestSuite))
}

func (s *NoteHandlerTestSuite) TestCreateNote_WhenParamsAreValid_CreateNote() {
	body := `{"content": "Test note content"}`

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Engine.ServeHTTP(w, req)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(s.T(), err)

	data, _ := json.Marshal(response.Data)
	var note services.NoteRead
	err = json.Unmarshal(data, &note)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), http.StatusCreated, w.Code)
	assert.Equal(s.T(), "Test note content", note.Content)
	assert.NotZero(s.T(), note.CreatedAt)
	assert.NotZero(s.T(), note.UpdatedAt)
}

func (s *NoteHandlerTestSuite) TestCreateNote_WhenContentIsMissing_ReturnBadRequest() {
	body := `{"content": null}`

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Engine.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *NoteHandlerTestSuite) TestCreateNote_WhenContentIsEmpty_ReturnBadRequest() {
	body := `{"content": null}`

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/notes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Engine.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *NoteHandlerTestSuite) TestListNotes_WhenNoNotes_ReturnEmptyList() {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/notes", nil)
	w := httptest.NewRecorder()
	s.Engine.ServeHTTP(w, req)

	var response response.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, w.Code)
	assert.NotNil(s.T(), response.Meta)
	assert.Equal(s.T(), 0, response.Meta.TotalItems)
	assert.Equal(s.T(), 0, response.Meta.TotalPages)
	assert.Equal(s.T(), 1, response.Meta.Page)
	assert.Equal(s.T(), 10, response.Meta.PageSize)

	data, _ := json.Marshal(response.Data)

	var notes []services.NoteRead
	err = json.Unmarshal(data, &notes)

	assert.NoError(s.T(), err)
	assert.IsType(s.T(), []services.NoteRead{}, notes)
	assert.Len(s.T(), notes, 0)
}

func (s *NoteHandlerTestSuite) TestListNotes_WhenPageIsInvalid_ReturnsFirstPage() {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/notes?page=abc", nil)
	w := httptest.NewRecorder()
	s.Engine.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *NoteHandlerTestSuite) TestListNotes_WhenPageSizeIsInvalid_ReturnsDefaultPageSize() {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/notes?page_size=abc", nil)
	w := httptest.NewRecorder()
	s.Engine.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusOK, w.Code)
}

func (s *NoteHandlerTestSuite) TestListNotes_WhenSortIsInvalid_ReturnsBadRequest() {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/notes?sort_by=invalid", nil)
	w := httptest.NewRecorder()
	s.Engine.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}
