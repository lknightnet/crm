package repository

import (
	"errors"
	"gorm.io/gorm"
	"project-service/internal/model"
	"project-service/internal/repository/customRepositoryError"
	"project-service/pkg/database"
)

type projectUsersRepository struct {
	db *database.PostgreSQL
}

func (p *projectUsersRepository) AddProjectUser(projectUser *model.ProjectUsers) error {
	var project *model.Project
	err := p.db.DB.Where("id = ?", projectUser.ProjectID).First(project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customRepositoryError.ErrProjectNotFound
		}
		return err
	}

	if err := p.db.DB.Create(projectUser).Error; err != nil {
		return err
	}
	return nil
}

func (p *projectUsersRepository) RemoveProjectUser(projectUser *model.ProjectUsers) error {
	var project *model.Project
	err := p.db.DB.Where("id = ?", projectUser.ProjectID).First(project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customRepositoryError.ErrProjectNotFound
		}
		return err
	}

	if err := p.db.DB.
		Where("project_id = ? AND user_id = ?", projectUser.ProjectID, projectUser.UserID).
		Delete(&model.ProjectUsers{}).Error; err != nil {
		return err
	}
	return nil
}

func newProjectUsersRepository(db *database.PostgreSQL) *projectUsersRepository {
	return &projectUsersRepository{db: db}
}
