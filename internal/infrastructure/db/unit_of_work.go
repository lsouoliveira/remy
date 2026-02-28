package db

import (
	infraRepo "remy/internal/infrastructure/db/repository"
	"remy/internal/models"
	"remy/internal/repository"
	"remy/internal/unit_of_work"

	"gorm.io/gorm"
)

type GormUnitOfWorkFactory struct {
	db        *gorm.DB
	publisher models.DomainEventPublisher
}

func NewGormUnitOfWorkFactory(db *gorm.DB, publisher models.DomainEventPublisher) *GormUnitOfWorkFactory {
	return &GormUnitOfWorkFactory{db: db, publisher: publisher}
}

func (f *GormUnitOfWorkFactory) New() unit_of_work.UnitOfWork {
	return &GormUnitOfWork{tx: f.db.Begin(), publisher: f.publisher}
}

type GormUnitOfWork struct {
	tx        *gorm.DB
	publisher models.DomainEventPublisher
}

func (u *GormUnitOfWork) Notes() repository.NoteRepository {
	return infraRepo.NewNoteRepository(u.tx, u.publisher)
}

func (u *GormUnitOfWork) Commit() error {
	return u.tx.Commit().Error
}

func (u *GormUnitOfWork) Rollback() {
	u.tx.Rollback()
}
