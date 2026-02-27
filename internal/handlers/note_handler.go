package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"remy/internal/response"
	"remy/internal/services"
)

type NoteHandler struct {
	service *services.NoteService
}

type CreateNoteRequest struct {
	Content string `json:"content" binding:"required"`
}

type NoteListRequest struct {
	Page     int `form:"page" binding:"omitempty"`
	PageSize int `form:"page_size" binding:"omitempty"`
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
		Content: req.Content,
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

	result, err := h.service.List(services.NewListNotesRequest(req.Page, req.PageSize))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.NewPaginatedResponse(result.Notes, req.Page, req.PageSize, int(result.Total)))
}
