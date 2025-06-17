package service

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
)

type AuthService interface {
	Login(email, password string) (*model.Tokens, error)
	SignUp(email, name, password string) error
	Verify(oneTimeToken string) error
	Refresh(refreshToken string) (*model.Tokens, error)
}

type Service struct {
	AuthService AuthService
}

type DependenciesService struct {
	AuthSignature string

	JWTServiceApi  string
	JWTServicePort string

	UserRepository repository.UserRepository
	AuthRepository repository.AuthRepository
}

func NewService(deps *DependenciesService) *Service {
	return &Service{
		AuthService: newAuthService(deps.AuthSignature, deps.AuthRepository, deps.UserRepository, deps.JWTServiceApi, deps.JWTServicePort),
	}
}
