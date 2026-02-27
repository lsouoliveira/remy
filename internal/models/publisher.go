package models

type DomainEventPublisher interface {
	Publish(event DomainEvent) error
}

type DomainEventHandler interface {
	Handle(event DomainEvent) error
}
