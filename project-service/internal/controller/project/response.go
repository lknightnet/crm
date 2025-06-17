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

type ProjectProjectUsersResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProjectID int       `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ProjectWithProjectUsersResponse struct {
	ID           int                           `json:"id"`
	Name         string                        `json:"name"`
	Description  *string                       `json:"description"`
	CreatedID    int                           `json:"created_id"`
	Visibility   bool                          `json:"visibility"`
	ProjectUsers []ProjectProjectUsersResponse `json:"project_users"`
}

func NewProjectWithProjectUsersResponse(p *model.ProjectWithProjectUsers) *ProjectWithProjectUsersResponse {
	projectUsers := make([]ProjectProjectUsersResponse, len(p.ProjectUsers))
	for i, user := range p.ProjectUsers {
		projectUsers[i] = ProjectProjectUsersResponse{
			ID:        user.ID,
			UserID:    user.UserID,
			ProjectID: user.ProjectID,
			CreatedAt: user.CreatedAt,
		}
	}

	return &ProjectWithProjectUsersResponse{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		CreatedID:    p.CreatedID,
		Visibility:   p.Visibility,
		ProjectUsers: projectUsers,
	}
}

func NewProjectsWithProjectUsersResponse(projects []model.ProjectWithProjectUsers) []ProjectWithProjectUsersResponse {
	res := make([]ProjectWithProjectUsersResponse, len(projects))
	for i, p := range projects {
		res[i] = *NewProjectWithProjectUsersResponse(&p)
	}
	return res
}
