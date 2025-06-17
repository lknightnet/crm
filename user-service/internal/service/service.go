package service

import (
	"user-service/internal/model"
	"user-service/internal/repository"
)

type UserService interface {
	GetUser(token string) (*model.User, error)
	ValidateToken(token string) error
	GetUserById(id string) (*model.User, error)
	GetUsersByUsername(username string) ([]model.User, error)
	UpdateUserByID(token string, name, email *string) error
}

type Service struct {
	UserService UserService
}

type DependenciesService struct {
	JWTServiceApi  string
	JWTServicePort string

	UserRepository repository.UserRepository
}

func NewService(deps *DependenciesService) *Service {
	return &Service{
		UserService: newUserService(deps.UserRepository, deps.JWTServiceApi, deps.JWTServicePort),
	}
}
