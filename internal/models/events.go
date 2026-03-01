package models

import (
	"github.com/google/uuid"
)

type DomainEvent any

type NoteCreatedEvent struct {
	NoteID uuid.UUID
}

func NewNoteCreatedEvent(noteID uuid.UUID) *NoteCreatedEvent {
	return &NoteCreatedEvent{
		NoteID: noteID,
	}
}
