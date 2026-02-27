package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"remy/internal/helpers"
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
	PageSize string `form:"page_size" binding:"omitempty,min=1"`
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
		c.Error(err)
		return
	}

	var params services.ListNotesParams
	params = services.ListNotesParams{
		Page:     helpers.ParseInt(req.Page, 1),
		PageSize: helpers.ParseInt(req.PageSize, 10),
	}

	result, err := h.service.List(params)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(result.Notes, params.Page, params.PageSize, int(result.Total)))
}
