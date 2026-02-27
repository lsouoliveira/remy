package mocks

import (
	"github.com/stretchr/testify/mock"
	"remy/internal/models"
)

type MockDomainEventPublisher struct {
	mock.Mock
}

func (m *MockDomainEventPublisher) Publish(event models.DomainEvent) error {
	args := m.Called(event)
	return args.Error(0)
}
