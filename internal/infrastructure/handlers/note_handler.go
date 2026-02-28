package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"remy/internal/helpers"
	infraErrors "remy/internal/infrastructure/errors"
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
	Page     string `form:"page" binding:"omitempty,min=1"`
	PageSize string `form:"page_size" binding:"omitempty,min=1,max=100"`
	SortBy   string `form:"sort_by" binding:"omitempty,oneof=created_at updated_at review_at"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc"`
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

	params := services.ListNotesParams{
		Page:     getPageParam(c),
		PageSize: getPageSizeParam(c),
		SortBy:   c.DefaultQuery("sort_by", "created_at"),
		Order:    c.DefaultQuery("order", "asc"),
	}

	result, err := h.service.List(params)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(result.Notes, params.Page, params.PageSize, int(result.Total)))
}

func getPageParam(c *gin.Context) int {
	pageStr := c.Query("page")
	return max(helpers.ParseInt(pageStr, 1), 1)
}

func getPageSizeParam(c *gin.Context) int {
	pageSizeStr := c.Query("page_size")
	pageSize := helpers.ParseInt(pageSizeStr, 10)

	if pageSize < 1 {
		pageSize = 10
	}

	return pageSize
}
