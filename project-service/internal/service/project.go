package service

import (
	"errors"
	"fmt"
	"project-service/internal/model"
	"project-service/internal/repository"
	"project-service/internal/repository/customRepositoryError"
	"project-service/internal/service/customServiceError"
	"project-service/pkg/tg"
	"strconv"
	"time"
)

type projectService struct {
	UserServiceApi         string
	UserServicePort        string
	ProjectRepository      repository.ProjectRepository
	ProjectUsersRepository repository.ProjectUsersRepository
}

func (p *projectService) GetProjectsByName(token string, projectName string) ([]model.ProjectWithProjectUsers, error) {
	user, err := GetUserByToken(token, p.UserServiceApi, p.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(fmt.Errorf("method: GetUserByToken, error: %s", err.Error()), "api/project/get/like")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(fmt.Errorf("method: GetUserByToken, error: %s", err.Error()), "api/project/get/like")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(fmt.Errorf("method: GetUserByToken, error: %s", err.Error()), "api/project/get/like")
		return nil, customServiceError.ErrUnknownError
	}

	projects, err := p.ProjectRepository.GetProjectsByName(user.ID, projectName)
	if err != nil {
		return nil, customServiceError.ErrUnknownError
	}

	var projectUser []model.ProjectWithProjectUsers
	for i, project := range projects {

		_, projectUsers, err := p.ProjectRepository.GetProjectByID(project.ID)
		if err != nil {
			if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
				tg.SendError(err, "api/project/get/like")
				return nil, customServiceError.ErrProjectNotFound
			}
		}

		projectUser[i] = model.ProjectWithProjectUsers{
			ID:           project.ID,
			Name:         project.Name,
			Description:  project.Description,
			CreatedID:    project.CreatedID,
			Visibility:   true,
			ProjectUsers: projectUsers,
		}
	}

	return projectUser, nil
}

func (p *projectService) EditProject(projectID int, name, description *string, projectUsers *[]int) error {
	project := &model.Project{
		ID: projectID,
	}

	if name != nil {
		project.Name = *name
	}
	if description != nil {
		project.Description = description
	}

	err := p.ProjectRepository.EditProject(project)
	if err != nil {
		tg.SendError(err.Error(), "api/project/update")
		return customServiceError.ErrUnknownError
	}

	err = p.ProjectUsersRepository.RemoveProjectUsers(projectID)
	if err != nil {
		tg.SendError(err.Error(), "api/project/update")
		return customServiceError.ErrUnknownError
	}
	if projectUsers != nil {
		for _, projectUser := range *projectUsers {
			err = p.ProjectUsersRepository.AddProjectUser(&model.ProjectUsers{
				UserID:    projectUser,
				ProjectID: projectID,
				CreatedAt: time.Now(),
			})
			if err != nil {
				tg.SendError(err.Error(), "api/project/update")
				return customServiceError.ErrUnknownError
			}
		}
	}

	return nil
}

func (p *projectService) GetProjectsByToken(token string) ([]model.ProjectWithCreatedUsername, error) {

	user, err := GetUserByToken(token, p.UserServiceApi, p.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(fmt.Errorf("method: GetUserByToken, error: %s", err.Error()), "api/project/get/list")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(fmt.Errorf("method: GetUserByToken, error: %s", err.Error()), "api/project/get/list")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(fmt.Errorf("method: GetUserByToken, error: %s", err.Error()), "api/project/get/list")
		return nil, customServiceError.ErrUnknownError
	}

	projects, err := p.ProjectRepository.GetProjectsByUserID(user.ID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrProjectsNotFound) {
			tg.SendError(fmt.Errorf("method: GetProjectsByUserID, error: %s", err.Error()), "api/project/get/list")
			return nil, customServiceError.ErrProjectsNotFound
		}
		tg.SendError(fmt.Errorf("method: GetProjectsByUserID, error: %s", err.Error()), "api/project/get/list")
		return nil, customServiceError.ErrUnknownError
	}

	var projectsWithResponsibleUsers []model.ProjectWithCreatedUsername
	for _, project := range projects {
		user, err := GetUserByID(strconv.Itoa(project.CreatedID), p.UserServiceApi, p.UserServicePort)
		if err != nil {
			tg.SendError(fmt.Errorf("method: GetUserByID, error: %s", err.Error()), "api/project/get/list")
			return nil, customServiceError.ErrUnknownError
		}

		projectsWithResponsibleUsers = append(projectsWithResponsibleUsers, model.ProjectWithCreatedUsername{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			CreatedName: user.Name,
			CreatedID:   project.CreatedID,
			Visibility:  project.Visibility,
		})
	}
	return projectsWithResponsibleUsers, nil
}

func (p *projectService) GetProjectByID(token string, projectID int) (*model.Project, []model.ProjectUsers, error) {
	user, err := GetUserByToken(token, p.UserServiceApi, p.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/project/get/:id")
			return nil, nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/project/get/:id")
			return nil, nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/project/get/:id")
		return nil, nil, customServiceError.ErrUnknownError
	}

	project, users, err := p.ProjectRepository.GetProjectByID(projectID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
			tg.SendError(err, "api/project/get/:id")
			return nil, nil, customServiceError.ErrProjectNotFound
		}
	}

	isFound := false
	for _, projectUsers := range users {
		if projectUsers.UserID == user.ID {
			isFound = true
			break
		}
	}

	if isFound {
		return project, users, nil
	} else {
		return nil, nil, customServiceError.ErrPermissionDenied
	}
}

func (p *projectService) CreateProject(token string, name string, description *string, projectUsers []int) error {
	user, err := GetUserByToken(token, p.UserServiceApi, p.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/project/create")
			return customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/project/create")
			return customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/project/create")
		return customServiceError.ErrUnknownError
	}

	project := &model.Project{
		Name:        name,
		Description: description,
		Visibility:  true,
		CreatedID:   user.ID,
	}

	projectID, err := p.ProjectRepository.CreateProject(project)
	if err != nil {
		tg.SendError(err, "api/project/create")
		return customServiceError.ErrUnknownError
	}

	err = p.ProjectUsersRepository.AddProjectUser(&model.ProjectUsers{
		UserID:    user.ID,
		ProjectID: projectID,
		CreatedAt: time.Now(),
	})

	for _, projectUser := range projectUsers {
		err = p.ProjectUsersRepository.AddProjectUser(&model.ProjectUsers{
			UserID:    projectUser,
			ProjectID: projectID,
			CreatedAt: time.Now(),
		})
	}

	if err != nil {
		tg.SendError(err, "api/project/create")
		return customServiceError.ErrUnknownError
	}
	return nil
}

func newProjectService(projectRepository repository.ProjectRepository, projectUsersRepository repository.ProjectUsersRepository, userServiceApi, userServicePort string) *projectService {
	return &projectService{
		UserServiceApi:         userServiceApi,
		UserServicePort:        userServicePort,
		ProjectRepository:      projectRepository,
		ProjectUsersRepository: projectUsersRepository,
	}
}
