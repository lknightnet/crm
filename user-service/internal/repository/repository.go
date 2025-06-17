package repository

import (
	"user-service/internal/model"
	"user-service/pkg/database"
)

type UserRepository interface {
	GetUserByID(userID string) (*model.User, error)
	GetUsersByUsername(username string) ([]model.User, error)
	UpdateUserByID(userID int, user *model.User) error
}

type Repository struct {
	UserRepository UserRepository
}

func NewRepositories(db *database.PostgreSQL) *Repository {
	return &Repository{
		UserRepository: newUserRepository(db),
	}
}
