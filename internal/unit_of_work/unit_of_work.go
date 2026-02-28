package unit_of_work

import "remy/internal/repository"

type UnitOfWork interface {
	Notes() repository.NoteRepository
	Commit() error
	Rollback()
}

type UnitOfWorkFactory interface {
	New() UnitOfWork
}
