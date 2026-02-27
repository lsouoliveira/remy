package models

type DomainEvent any

type NoteCreatedEvent struct {
	NoteID uint
}

func NewNoteCreatedEvent(noteID uint) *NoteCreatedEvent {
	return &NoteCreatedEvent{
		NoteID: noteID,
	}
}
