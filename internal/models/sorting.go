package models

import "github.com/google/uuid"

type (
	SortField int
	SortOrder int
)

type Cursor struct {
	ID    uuid.UUID
	Value any
	Field SortField
	Order SortOrder
}

const (
	Asc SortOrder = iota
	Desc
)

const (
	SortByCreatedAt SortField = iota
	SortByUpdatedAt
	SortByReviewAt
)
