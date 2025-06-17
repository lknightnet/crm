package repository

import (
	"errors"
	"gorm.io/gorm"
	"user-service/internal/model"
	"user-service/internal/repository/customRepositoryError"
	"user-service/pkg/database"
)

type userRepository struct {
	db *database.PostgreSQL
}

func (u *userRepository) UpdateUserByID(userID int, user *model.User) error {
	updates := make(map[string]interface{})

	if user.Name != "" {
		updates["name"] = user.Name
	}

	if user.Email != "" {
		updates["email"] = user.Email
	}

	if len(updates) == 0 {
		return nil
	}

	err := u.db.DB.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates).Error

	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUsersByUsername(username string) ([]model.User, error) {
	var users []model.User

	if err := u.db.DB.Where("name ILIKE ?", "%"+username+"%").Find(&users).Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, customRepositoryError.ErrUsersNotFound
	}

	return users, nil
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

func newUserRepository(db *database.PostgreSQL) *userRepository {
	return &userRepository{
		db: db,
	}
}
