package infrastructure

import (
	"remy/internal/models"
)

type InMemoryPublisher struct {
	handles []models.DomainEventHandler
}

func NewInMemoryPublisher() *InMemoryPublisher {
	return &InMemoryPublisher{
		handles: []models.DomainEventHandler{},
	}
}

func (p *InMemoryPublisher) Publish(event models.DomainEvent) error {
	for _, handle := range p.handles {
		if err := handle.Handle(event); err != nil {
			return err
		}
	}

	return nil
}

func (p *InMemoryPublisher) RegisterHandle(handle models.DomainEventHandler) {
	p.handles = append(p.handles, handle)
}
