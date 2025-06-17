package service

import (
	"errors"
	"project-service/internal/model"
	"project-service/internal/repository"
	"project-service/internal/repository/customRepositoryError"
	"project-service/internal/service/customServiceError"
	"time"
)

type projectUserService struct {
	ProjectUsersRepository repository.ProjectUsersRepository
}

func (p *projectUserService) AddProjectUser(userID, projectID int) error {
	projectUser := &model.ProjectUsers{
		UserID:    userID,
		ProjectID: projectID,
		CreatedAt: time.Now(),
	}
	err := p.ProjectUsersRepository.AddProjectUser(projectUser)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
			return customServiceError.ErrProjectNotFound
		}
	}
	return nil
}

func (p *projectUserService) RemoveProjectUser(userID, projectID int) error {
	projectUser := &model.ProjectUsers{
		UserID:    userID,
		ProjectID: projectID,
	}
	err := p.ProjectUsersRepository.RemoveProjectUser(projectUser)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
			return customServiceError.ErrProjectNotFound
		}
	}
	return nil
}

func newProjectUserService(projectRepository repository.ProjectUsersRepository) *projectUserService {
	return &projectUserService{projectRepository}
}
