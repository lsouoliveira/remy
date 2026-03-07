package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	infraErrors "remy/internal/infrastructure/errors"
	"remy/internal/infrastructure/helpers"
	"remy/internal/models"
	"remy/internal/response"
	"remy/internal/services"
)

type NoteHandler struct {
	service *services.NoteService
}

type CreateNoteRequest struct {
	Content *string `json:"content" binding:"required,min=1"`
}

type NoteListRequest struct {
	SortBy string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at review_at"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc"`
	Cursor string `form:"next_cursor" binding:"omitempty"`
	Limit  string `form:"limit" binding:"omitempty"`
}

type NoteListCursor struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	SortBy    string    `json:"sort_by"`
	Order     string    `json:"order"`
	SortValue string    `json:"sort_value"`
}

func NewNoteHandler(service *services.NoteService) *NoteHandler {
	return &NoteHandler{
		service: service,
	}
}

func (h *NoteHandler) Create(c *gin.Context) {
	var req CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	note, err := h.service.Create(services.NoteCreate{
		Content: *req.Content,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, response.APIResponse{
		Data: note,
	})
}

func (h *NoteHandler) List(c *gin.Context) {
	var req NoteListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			c.Error(infraErrors.NewQueryValidationError(validationErrs))
		} else {
			c.Error(err)
		}

		return
	}

	cursor, err := parseCursor(req.Cursor)
	if err != nil {
		c.Error(infraErrors.InvalidCursorParameter())
		return
	}

	limit := helpers.ParseInt(req.Limit, 20)

	result, err := h.service.List(services.ListNotesParams{})
	if err != nil {
		c.Error(err)
		return
	}

	nextCursor := newNoteListCursor(result.Notes, req.SortBy, req.Order)

	c.JSON(http.StatusOK, response.NewPaginatedResponse(result.Notes, nextCursor, int(result.Total), limit))
}

func parseSortBy(s string) (models.SortField, error) {
	switch s {
	case "created_at":
		return models.SortByCreatedAt, nil
	case "updated_at":
		return models.SortByUpdatedAt, nil
	case "review_at":
		return models.SortByReviewAt, nil
	default:
		return 0, errors.New("invalid sort_by value")
	}
}

func parseOrder(s string) (models.SortOrder, error) {
	switch s {
	case "asc":
		return models.Asc, nil
	case "desc":
		return models.Desc, nil
	default:
		return 0, errors.New("invalid order value")
	}
}

func parseCursor(s string) (*NoteListCursor, error) {
	cursor, err := helpers.DecodeCursor[NoteListCursor](s)
	if err != nil {
		return nil, err
	}

	if cursor == nil {
		defaultCursor := defaultCursor()
		return &defaultCursor, nil
	}

	return cursor, nil
}

func defaultCursor() NoteListCursor {
	return NoteListCursor{
		ID:        uuid.Nil,
		CreatedAt: time.Time{},
		SortBy:    "created_at",
		Order:     "asc",
		SortValue: "",
	}
}

func newNoteListCursor(notes []services.NoteRead, sortBy string, order string) NoteListCursor {
	if len(notes) == 0 {
		return NoteListCursor{}
	}

	lastNote := notes[len(notes)-1]

	return NoteListCursor{
		ID:        lastNote.ID,
		CreatedAt: lastNote.CreatedAt,
		SortBy:    sortBy,
		Order:     order,
		SortValue: getSortValue(lastNote, sortBy),
	}
}

func getSortValue(note services.NoteRead, sortBy string) string {
	switch sortBy {
	case "created_at":
		return note.CreatedAt.Format(time.RFC3339Nano)
	case "updated_at":
		return note.UpdatedAt.Format(time.RFC3339Nano)
	case "review_at":
		return note.ReviewAt.Format(time.RFC3339Nano)
	default:
		return ""
	}
}
