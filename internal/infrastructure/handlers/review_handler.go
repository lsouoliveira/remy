package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	infraErrors "remy/internal/infrastructure/errors"
	"remy/internal/infrastructure/helpers"
	"remy/internal/response"
	"remy/internal/services"
)

type ReviewHandler struct {
	service *services.NoteService
}

type ReviewRequest struct {
	Quality int `json:"quality" binding:"required,min=1,max=5"`
}

type ReviewListRequest struct {
	Page     string `form:"page" binding:"omitempty,min=1"`
	PageSize string `form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy   string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at review_at"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
}

type ReviewListCursor struct {
	ID     string `json:"id"`
	SortBy string `json:"sort_value"`
	Order  string `json:"order"`
}

func NewReviewHandler(service *services.NoteService) *ReviewHandler {
	return &ReviewHandler{
		service: service,
	}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	noteID, err := helpers.ParseUUID(c.Param("id"))
	if err != nil {
		c.Error(infraErrors.InvalidPathParameter("id"))
		return
	}

	err = h.service.Review(services.ReviewParams{
		NoteID:  noteID,
		Quality: req.Quality,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ReviewHandler) List(c *gin.Context) {
	var req ReviewListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			c.Error(infraErrors.NewQueryValidationError(validationErrs))
		} else {
			c.Error(err)
		}

		return
	}

	params := services.ListNotesParams{
		Page:     helpers.GetPageParam(c),
		PageSize: helpers.GetPageSizeParam(c),
		SortBy:   c.DefaultQuery("sort_by", "review_at"),
		Order:    c.DefaultQuery("order", "asc"),
	}

	result, err := h.service.List(params)
	if err != nil {
		c.Error(err)
		return
	}

	cursor := newReviewListCursor(result.Notes, params.SortBy, params.Order)

	c.JSON(http.StatusOK, response.NewPaginatedResponse(result.Notes, cursor, int(result.Total), params.PageSize))
}

func newReviewListCursor(notes []services.NoteRead, sortBy string, order string) any {
	if len(notes) == 0 {
		return NoteListCursor{}
	}

	lastNote := notes[len(notes)-1]

	return ReviewListCursor{
		ID:     lastNote.ID.String(),
		SortBy: sortBy,
		Order:  order,
	}
}
