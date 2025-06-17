package project

import (
	"project-service/internal/model"
	"time"
)

type ProjectResponse struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedID   int     `json:"created_id"` //кто создал
}

type ProjectWithCreatedUsernameResponse struct {
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description"`
	CreatedUsername string  `json:"created_username"` //кто создал
	CreatedID       int     `json:"created_id"`       //кто создал
}

type ProjectUsersResponse struct {
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ProjectWithUsersResponse struct {
	Project      *ProjectResponse       `json:"project"`
	ProjectUsers []ProjectUsersResponse `json:"project_users"`
}

func NewProjectResponse(project *model.Project) *ProjectResponse {
	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedID:   project.CreatedID,
	}
}

func NewProjectWithCreatedUsernameResponse(project *model.ProjectWithCreatedUsername) *ProjectWithCreatedUsernameResponse {
	return &ProjectWithCreatedUsernameResponse{
		ID:              project.ID,
		Name:            project.Name,
		Description:     project.Description,
		CreatedUsername: project.CreatedName,
		CreatedID:       project.CreatedID,
	}
}

func NewProjectsResponse(project []model.ProjectWithCreatedUsername) []ProjectWithCreatedUsernameResponse {
	projects := make([]ProjectWithCreatedUsernameResponse, len(project))
	for i, v := range project {
		projects[i] = *NewProjectWithCreatedUsernameResponse(&v)
	}
	return projects
}

func NewProjectWithUsersResponse(project *model.Project, users []model.ProjectUsers) *ProjectWithUsersResponse {

	projectUsers := make([]ProjectUsersResponse, len(users))
	for i, user := range users {
		projectUsers[i] = ProjectUsersResponse{
			UserID:    user.UserID,
			CreatedAt: user.CreatedAt,
		}
	}
	return &ProjectWithUsersResponse{
		Project:      NewProjectResponse(project),
		ProjectUsers: projectUsers,
	}
}

type ProjectErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}
