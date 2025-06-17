package repository

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/customRepositoryError"
	"auth-service/pkg/database"
	"errors"
	"gorm.io/gorm"
)

type authRepository struct {
	db *database.PostgreSQL
}

func (a *authRepository) CreateRouteLink(routeLink *model.RouteLink) error {
	if err := a.db.DB.Create(routeLink).Error; err != nil {
		return err
	}
	return nil
}

func (a *authRepository) GetRouteLinkByToken(token string) (*model.RouteLink, error) {
	var routeLink model.RouteLink
	if err := a.db.DB.Where("one_time_token = ?", token).First(&routeLink).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrRouteLinkNotFound
		}
		return nil, err
	}
	return &routeLink, nil
}

func newAuthRepository(db *database.PostgreSQL) *authRepository {
	return &authRepository{
		db: db,
	}
}
