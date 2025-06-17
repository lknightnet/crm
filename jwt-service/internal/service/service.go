package service

import (
	"jwt-service/internal/model"
	"jwt-service/internal/repository"
	"time"
)

type JWTService interface {
	GenerateTokens(userID string) (*model.Tokens, error)
	RefreshTokens(refreshToken string) (*model.Tokens, error)
	ValidateToken(accessToken string) (*model.AccessToken, error)
}

type Service struct {
	JWTService JWTService
}

type DependenciesService struct {
	SignKey       []byte
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	JWTRepository repository.JWTRepository
}

func NewService(deps *DependenciesService) *Service {
	return &Service{
		JWTService: newJWTService(deps.SignKey, deps.AccessExpiry, deps.RefreshExpiry, deps.JWTRepository),
	}
}
