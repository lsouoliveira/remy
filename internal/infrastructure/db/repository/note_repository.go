package repository

import (
	"fmt"

	"github.com/google/uuid"
	gormpkg "gorm.io/gorm"
	"gorm.io/gorm/clause"

	"remy/internal/domainErrors/general"
	infraErrors "remy/internal/infrastructure/errors"
	"remy/internal/models"
	"remy/internal/repository"
)

type paginationFilter struct {
	sortBy    string
	operator  string
	isDesc    bool
	value     any
	id        uuid.UUID
	hasCursor bool
}

var SortFieldMap = map[models.SortField]string{
	models.SortByCreatedAt: "created_at",
	models.SortByUpdatedAt: "updated_at",
	models.SortByReviewAt:  "review_at",
}

type NoteRepository struct {
	db        *gormpkg.DB
	publisher models.DomainEventPublisher
}

func NewNoteRepository(db *gormpkg.DB, publisher models.DomainEventPublisher) repository.NoteRepository {
	return &NoteRepository{db: db, publisher: publisher}
}

func (r *NoteRepository) Save(note *models.Note) error {
	if note.IsNew {
		if err := r.db.Create(note).Error; err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}
	} else {
		result := r.db.Model(note).Where("version = ?", note.Version-1).Updates(note)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return infraErrors.ErrVersionConflict
		}
	}

	for _, event := range note.FlushEvents() {
		if err := r.publisher.Publish(event); err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}

	return nil
}

func (r *NoteRepository) GetByID(id uuid.UUID) (*models.Note, error) {
	var note models.Note
	if err := r.db.First(&note, id).Error; err != nil {
		return nil, general.NotFound("note", id)
	}

	return &note, nil
}

func (r *NoteRepository) List(params repository.ListParams) ([]*models.Note, bool, error) {
	var notes []*models.Note
	query := r.db.Model(&models.Note{})

	filter, err := parseCursor(params.SortBy, params.Order, params.Cursor)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse cursor: %w", err)
	}

	if filter.hasCursor {
		query = query.Where(clause.Expr{
			SQL:  fmt.Sprintf("(%s, id) %s (?, ?)", filter.sortBy, filter.operator),
			Vars: []any{filter.value, filter.id},
		})
	}

	query = query.Order(clause.OrderBy{Columns: []clause.OrderByColumn{
		{Column: clause.Column{Name: filter.sortBy}, Desc: filter.isDesc},
		{Column: clause.Column{Name: "id"}, Desc: filter.isDesc},
	}})

	query = query.Limit(params.Limit + 1)

	if err := query.Find(&notes).Error; err != nil {
		return nil, false, fmt.Errorf("failed to list notes: %w", err)
	}

	hasNextPage := len(notes) > params.Limit
	if hasNextPage {
		notes = notes[:params.Limit]
	}

	return notes, hasNextPage, nil
}

func parseCursor(sortBy models.SortField, order models.SortOrder, cursor *models.Cursor) (*paginationFilter, error) {
	filter := &paginationFilter{}

	if cursor != nil {
		sortBy = cursor.Field
		order = cursor.Order
		filter.hasCursor = true
		filter.value = cursor.Value
		filter.id = cursor.ID
	}

	sortByField, err := mapSortField(sortBy)
	if err != nil {
		return nil, err
	}

	operator, isDesc, err := mapSortOrder(order)
	if err != nil {
		return nil, err
	}

	filter.sortBy = sortByField
	filter.isDesc = isDesc
	filter.operator = operator

	return filter, nil
}

func mapSortField(sortField models.SortField) (string, error) {
	if field, ok := SortFieldMap[sortField]; ok {
		return field, nil
	}

	return "", fmt.Errorf("unsupported sort field: %d", sortField)
}

func mapSortOrder(order models.SortOrder) (operator string, isDesc bool, err error) {
	switch order {
	case models.Asc:
		return ">", false, nil
	case models.Desc:
		return "<", true, nil
	default:
		return "", false, fmt.Errorf("unsupported sort order: %d", order)
	}
}
