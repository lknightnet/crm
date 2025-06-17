package repository

import (
	"errors"
	"gorm.io/gorm"
	"project-service/internal/model"
	"project-service/internal/repository/customRepositoryError"
	"project-service/pkg/database"
)

type projectRepository struct {
	db *database.PostgreSQL
}

func (p *projectRepository) GetProjectByID(projectID int) (*model.Project, []model.ProjectUsers, error) {
	var project model.Project
	err := p.db.DB.
		Where("id = ?", projectID).
		First(&project).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, customRepositoryError.ErrProjectNotFound
		}
		return nil, nil, err
	}

	var users []model.ProjectUsers
	err = p.db.DB.
		Where("project_id = ?", project.ID).
		Find(&users).Error
	if err != nil {
		return nil, nil, err
	}
	return &project, users, nil
}

func (p *projectRepository) GetProjectsByUserID(userID int) ([]model.Project, error) {
	var projects []model.Project
	err := p.db.DB.Joins("JOIN project_users ON project_users.project_id = projects.id").
		Where("project_users.user_id = ?", userID).
		Find(&projects).Error
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, customRepositoryError.ErrProjectsNotFound
	}
	return projects, nil
}

func (p *projectRepository) CreateProject(project *model.Project) (int, error) {
	if err := p.db.DB.Create(project).Error; err != nil {
		return 0, err
	}
	return project.ID, nil
}

func newProjectRepository(db *database.PostgreSQL) *projectRepository {
	return &projectRepository{
		db: db,
	}
}
