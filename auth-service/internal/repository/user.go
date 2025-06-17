package repository

import (
	"auth-service/internal/model"
	"auth-service/internal/repository/customRepositoryError"
	"auth-service/pkg/database"
	"errors"
	"gorm.io/gorm"
	"time"
)

type userRepository struct {
	db *database.PostgreSQL
}

func (u *userRepository) UpdateLastLoginUser(userID int, lastLoginUser time.Time) error {
	err := u.db.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Update("last_login", lastLoginUser).Error

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserByID(userID string) (*model.User, error) {
	var user model.User
	err := u.db.DB.Where("id = ? AND active = ?", userID, true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) UpdateActiveByUUID(userUUID string) error {
	var user *model.User
	result := u.db.DB.Model(user).Where("uuid = ?", userUUID).Update("active", true)
	return result.Error
}

func (u *userRepository) UpdateNameByUUID(userUUID string, name string) error {
	var user *model.User
	result := u.db.DB.Model(user).Where("uuid = ?", userUUID).Update("name", name)
	return result.Error
}

func (u *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := u.db.DB.Where("email = ? AND active = ?", email, true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customRepositoryError.ErrInvalidEmailOrPassword
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) CreateUser(user *model.User) error {
	if err := u.db.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func newUserRepository(db *database.PostgreSQL) *userRepository {
	return &userRepository{
		db: db,
	}
}
