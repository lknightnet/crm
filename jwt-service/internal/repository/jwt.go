package repository

import (
	"errors"
	"gorm.io/gorm"
	"jwt-service/internal/model"
	"jwt-service/internal/repository/customRepositoryError"
	"jwt-service/pkg/database"
)

type jwtRepository struct {
	db *database.PostgreSQL
}

func (j *jwtRepository) CreateTokens(accessToken *model.AccessToken, refreshToken *model.RefreshToken) error {
	tx := j.db.DB.Begin()

	if tx.Error != nil {
		return tx.Error
	}

	if err := j.db.DB.Create(accessToken).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := j.db.DB.Create(refreshToken).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (j *jwtRepository) GetRefreshToken(refreshToken string) (*model.RefreshToken, error) {
	var refresh model.RefreshToken

	if err := j.db.DB.Where("token = ?", refreshToken).First(&refresh).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, err
	}
	return &refresh, nil
}

func (j *jwtRepository) GetAccessToken(accessToken string) (*model.AccessToken, error) {
	var access model.AccessToken

	if err := j.db.DB.Where("token = ?", accessToken).First(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrAccessTokenNotFound
		}
		return nil, err
	}
	return &access, nil
}

func newJWTRepository(db *database.PostgreSQL) *jwtRepository {
	return &jwtRepository{db: db}
}
