package repository

import (
	"auth-service/internal/model"
	"auth-service/pkg/database"
	"time"
)

type AuthRepository interface {
	CreateRouteLink(routeLink *model.RouteLink) error
	GetRouteLinkByToken(token string) (*model.RouteLink, error)
}

type UserRepository interface {
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateLastLoginUser(userID int, lastLoginUser time.Time) error
	UpdateActiveByUUID(userUUID string) error
}

type Repository struct {
	UserRepository UserRepository
	AuthRepository AuthRepository
}

func NewRepositories(db *database.PostgreSQL) *Repository {
	return &Repository{
		UserRepository: newUserRepository(db),
		AuthRepository: newAuthRepository(db),
	}
}
