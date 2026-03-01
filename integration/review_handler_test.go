package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"remy/internal/models"
	"remy/internal/testhelpers"
)

type ReviewHandlerSuite struct {
	testhelpers.IntegrationSuite
}

func TestReviewHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(ReviewHandlerSuite))
}

func (s *ReviewHandlerSuite) SetupSuite() {
	s.IntegrationSuite.SetupSuite()
}

func (s *ReviewHandlerSuite) TestCreateReview_WhenParamsAreValid_ReturnsNoContent() {
	note := models.CreateNote("Test note content")
	err := s.DB.Create(note).Error

	assert.NoError(s.T(), err)

	body := `{"quality": 4}`

	w := s.Post(fmt.Sprintf("/api/v1/notes/%s/review", note.ID.String()), body)

	assert.Equal(s.T(), http.StatusNoContent, w.Code)
}

func (s *ReviewHandlerSuite) TestCreateReview_WhenQualityIsMissing_ReturnBadRequest() {
	note := models.CreateNote("Test note content")
	err := s.DB.Create(note).Error

	assert.NoError(s.T(), err)

	body := `{"quality": null}`

	w := s.Post("/api/v1/notes/1/review", body)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
}

func (s *ReviewHandlerSuite) TestCreateReview_WhenNoteDoesNotExist_ReturnNotFound() {
	body := `{"quality": 4}`

	w := s.Post(fmt.Sprintf("/api/v1/notes/%s/review", uuid.New().String()), body)

	assert.Equal(s.T(), http.StatusNotFound, w.Code)
}
