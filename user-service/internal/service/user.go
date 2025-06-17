package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/repository/customRepositoryError"
	"user-service/internal/service/customServiceError"
	"user-service/pkg/tg"
)

type userService struct {
	JWTServiceApi  string
	JWTServicePort string

	UserRepository repository.UserRepository
}

func (u *userService) UpdateUserByID(token string, name, email *string) error {
	userModel, err := u.GetUser(token)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUserNotFound) {
			tg.SendError(err, "api/user/update")
			return customServiceError.ErrUserNotFound
		}
		tg.SendError(err, "api/user/update")
		return customServiceError.ErrUnknown
	}

	user := &model.User{
		ID:    userModel.ID,
		Name:  *name,
		Email: *email,
	}
	err = u.UserRepository.UpdateUserByID(userModel.ID, user)
	if err != nil {
		tg.SendError(err, "api/user/update")
		return customServiceError.ErrUnknown
	}
	return nil
}

func (u *userService) GetUsersByUsername(username string) ([]model.User, error) {
	user, err := u.UserRepository.GetUsersByUsername(username)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUsersNotFound) {
			tg.SendError(err, "api/user/get/like")
			return nil, customServiceError.ErrUsersNotFound
		}
	}
	return user, nil
}

func (u *userService) GetUserById(id string) (*model.User, error) {
	user, err := u.UserRepository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUserNotFound) {
			tg.SendError(err, "api/outside/user/get/:id")
			return nil, customServiceError.ErrUserNotFound
		}
	}
	return user, nil
}

func (u *userService) GetUser(token string) (*model.User, error) {
	accessToken, err := validate(token, u.JWTServiceApi, u.JWTServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/user/get")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/user/get")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/user/get")
		return nil, customServiceError.ErrUnknown
	}

	userID := accessToken.Subject

	user, err := u.UserRepository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUserNotFound) {
			tg.SendError(err, "api/user/get")
			return nil, customServiceError.ErrUserNotFound
		}
		tg.SendError(err, "api/user/get")
		return nil, customServiceError.ErrUnknown
	}
	return user, nil
}

func (u *userService) ValidateToken(token string) error {
	accessToken, err := validate(token, u.JWTServiceApi, u.JWTServicePort)
	if err != nil {
		tg.SendError(err, "api/user/validate")
		return customServiceError.ErrUnknown
	}

	userID := accessToken.Subject

	_, err = u.UserRepository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrUserNotFound) {
			return customServiceError.ErrTokenNotValid
		}
		tg.SendError(err, "api/user/validate")
		return customServiceError.ErrUnknown
	}
	return nil
}

func newUserService(userRepository repository.UserRepository, jwtApi, jwtPort string) *userService {
	return &userService{
		UserRepository: userRepository,
		JWTServiceApi:  jwtApi,
		JWTServicePort: jwtPort,
	}
}

func validate(token, api, port string) (*model.AccessToken, error) {
	client := http.Client{Timeout: time.Second * 30}

	// Формируем URL
	url := fmt.Sprintf("http://%s:%s/api/validate", api, port)

	// Формируем тело запроса
	payload := map[string]string{
		"access_token": token,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		var errResp struct {
			Status bool   `json:"status"`
			Error  string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, customServiceError.ErrTokenNotFound
		}
		return nil, customServiceError.ErrTokenNotFound
	}
	if resp.StatusCode == http.StatusUnauthorized {
		var errResp struct {
			Status bool   `json:"status"`
			Error  string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, customServiceError.ErrTokenExpired
		}
		return nil, customServiceError.ErrTokenExpired
	}

	// Проверяем код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// Декодируем ответ
	var accessToken model.AccessToken
	if err := json.NewDecoder(resp.Body).Decode(&accessToken); err != nil {
		return nil, err
	}

	return &accessToken, nil
}
