package repositories

import (
	"context"
	"database/sql"

	repository_interfaces "github.com/RandySteven/go-cook/db"
)

type (
	Repositories struct {
		UserRepository
		OnboardingRepository
	}

	Repository[T any] interface {
		repository_interfaces.Saver[T]
		repository_interfaces.Finder[T]
		repository_interfaces.Updater[T]
		repository_interfaces.Deleter[T]
	}
)

func NewRepositories(db *sql.DB) *Repositories {
	dbx := func(ctx context.Context) repository_interfaces.Trigger {
		return db
	}

	return &Repositories{
		UserRepository:       newUserRepository(dbx),
		OnboardingRepository: newOnboardingRepository(dbx),
	}
}
