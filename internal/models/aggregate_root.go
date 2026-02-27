package models

import (
	"gorm.io/gorm"
)

type AggregateRoot struct {
	gorm.Model
	events []DomainEvent
}

func (a *AggregateRoot) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

func (a *AggregateRoot) FlushEvents() []DomainEvent {
	events := a.events
	a.events = nil

	return events
}
