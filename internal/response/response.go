package response

import (
	"math"
)

type APIResponse struct {
	Data   any         `json:"data,omitempty"`
	Errors []*APIError `json:"errors,omitempty"`
	Meta   *Meta       `json:"meta,omitempty"`
}

type APIError struct {
	Status int     `json:"status"`
	Code   string  `json:"code"`
	Title  string  `json:"title"`
	Detail string  `json:"detail"`
	Source *Source `json:"source,omitempty"`
}

type Source struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
	Header    string `json:"header,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

func NewAPIError(status int, code string, title string, detail string) *APIError {
	return &APIError{
		Status: status,
		Code:   code,
		Title:  title,
		Detail: detail,
	}
}

func NewPaginatedResponse(data any, page int, pageSize int, totalItems int) APIResponse {
	var totalPages int

	if pageSize > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}

	return APIResponse{
		Data: data,
		Meta: &Meta{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			TotalItems: totalItems,
		},
	}
}
