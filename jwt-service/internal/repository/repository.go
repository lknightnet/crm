package repository

import (
	"jwt-service/internal/model"
	"jwt-service/pkg/database"
)

type JWTRepository interface {
	CreateTokens(accessToken *model.AccessToken, refreshToken *model.RefreshToken) error
	GetRefreshToken(refreshToken string) (*model.RefreshToken, error)
	GetAccessToken(accessToken string) (*model.AccessToken, error)
}

type Repository struct {
	JWTRepository JWTRepository
}

func NewRepositories(db *database.PostgreSQL) *Repository {
	return &Repository{
		JWTRepository: newJWTRepository(db),
	}
}
