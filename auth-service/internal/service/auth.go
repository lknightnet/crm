package service

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
	"auth-service/internal/repository/customRepositoryError"
	"auth-service/internal/service/customServiceError"
	"auth-service/pkg/tg"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

type authService struct {
	AuthSignature string

	JWTServiceApi  string
	JWTServicePort string

	UserRepository repository.UserRepository
	AuthRepository repository.AuthRepository
}

func (a *authService) Refresh(refreshToken string) (*model.Tokens, error) {

	tokens, err := refresh(refreshToken, a.JWTServiceApi, a.JWTServicePort)
	if err != nil {
		tg.SendError(err, "api/auth/refresh")
		return nil, customServiceError.ErrUnknown
	}

	return tokens, nil
}

func (a *authService) Login(email, password string) (*model.Tokens, error) {
	signingPassword, err := signedPassword(password, a.AuthSignature)
	if err != nil {
		tg.SendError(err, "api/auth/login")
		return nil, customServiceError.ErrUnknown
	}

	user, err := a.UserRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrInvalidEmailOrPassword) {
			return nil, customServiceError.ErrInvalidEmailOrPassword
		}
		tg.SendError(err, "api/auth/login")
		return nil, customServiceError.ErrUnknown
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signingPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, customServiceError.ErrInvalidEmailOrPassword
		}
		tg.SendError(err, "api/auth/login")
		return nil, customServiceError.ErrUnknown
	}

	tokens, err := generate(strconv.Itoa(user.ID), a.JWTServiceApi, a.JWTServicePort)
	if err != nil {
		tg.SendError(err, "api/auth/login")
		return nil, customServiceError.ErrUnknown
	}

	err = a.UserRepository.UpdateLastLoginUser(user.ID, time.Now())
	if err != nil {
		tg.SendError(err, "api/auth/login")
		return nil, customServiceError.ErrUnknown
	}

	return tokens, nil
}

func (a *authService) SignUp(email, name, password string) error {
	signingPassword, err := signedPassword(password, a.AuthSignature)
	if err != nil {
		tg.SendError(err, "api/auth/signup")
		return customServiceError.ErrUnknown
	}

	hashedPassword, err := hashPassword(signingPassword)
	if err != nil {
		tg.SendError(err, "api/auth/signup")
		return customServiceError.ErrUnknown
	}

	userUUID := uuid.NewString()

	user := &model.User{
		Email:     email,
		Name:      name,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UUID:      userUUID,
		Active:    true,
	}

	err = a.UserRepository.CreateUser(user)
	if err != nil {
		tg.SendError(err, "api/auth/signup")
		return customServiceError.ErrUnknown
	}

	routeLink := model.GenerateRouteLink(userUUID, email)
	err = a.AuthRepository.CreateRouteLink(routeLink)
	if err != nil {
		tg.SendError(err, "api/auth/signup")
		return customServiceError.ErrUnknown
	}

	return nil
}

func (a *authService) Verify(oneTimeToken string) error {
	routeLink, err := a.AuthRepository.GetRouteLinkByToken(oneTimeToken)
	if err != nil {

		if errors.Is(err, customRepositoryError.ErrRouteLinkNotFound) {
			return customServiceError.ErrOneTimeTokenNotFound
		}
		tg.SendError(err, "api/auth/verify")
		return customServiceError.ErrUnknown
	}

	if routeLink.ExpirationAt.Before(time.Now()) {
		return customServiceError.ErrOneTimeTokenExpired
	}

	err = a.UserRepository.UpdateActiveByUUID(routeLink.UUID)
	if err != nil {
		tg.SendError(err, "api/auth/verify")
		return customServiceError.ErrUnknown
	}

	return nil
}

func newAuthService(AuthSignature string, AuthRepository repository.AuthRepository, UserRepository repository.UserRepository, jwtApi, jwtPort string) *authService {
	return &authService{
		AuthSignature:  AuthSignature,
		UserRepository: UserRepository,
		AuthRepository: AuthRepository,
		JWTServiceApi:  jwtApi,
		JWTServicePort: jwtPort,
	}
}

func signedPassword(password string, signature string) (string, error) {
	h := hmac.New(sha256.New, []byte(signature))
	_, err := h.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func generate(userID, api, port string) (*model.Tokens, error) {
	client := http.Client{Timeout: time.Second * 30}

	// Формируем URL
	url := fmt.Sprintf("http://%s:%s/api/generate", api, port)

	// Формируем тело запроса
	payload := map[string]string{
		"user_id": userID,
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

	// Проверяем код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// Декодируем ответ
	var tokens model.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

func refresh(refreshToken, api, port string) (*model.Tokens, error) {
	client := http.Client{Timeout: time.Second * 30}

	// Формируем URL
	url := fmt.Sprintf("http://%s:%s/api/refresh", api, port)

	// Формируем тело запроса
	payload := map[string]string{
		"refresh_token": refreshToken,
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

	// Проверяем код ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	// Декодируем ответ
	var tokens model.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}
