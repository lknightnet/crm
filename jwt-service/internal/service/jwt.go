package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"jwt-service/internal/model"
	"jwt-service/internal/repository"
	"jwt-service/internal/repository/customRepositoryError"
	"jwt-service/internal/service/customServiceError"
	"jwt-service/pkg/tg"
	"time"
)

type jwtService struct {
	SignKey       []byte
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
	JWTRepository repository.JWTRepository
}

func (j *jwtService) GenerateTokens(userID string) (*model.Tokens, error) {
	accessTokenExp := time.Now().Add(j.AccessExpiry)
	accessTokenIss := time.Now()

	refreshTokenExp := time.Now().Add(j.RefreshExpiry)
	refreshTokenIss := time.Now()

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "jwt-service",
		Subject:   userID,
		Audience:  []string{"auth-service"},
		ExpiresAt: jwt.NewNumericDate(accessTokenExp),
		IssuedAt:  jwt.NewNumericDate(accessTokenIss),
	})

	accessToken, err := accessClaims.SignedString(j.SignKey)
	if err != nil {
		tg.SendError(err, "api/generate")
		return nil, err
	}

	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "jwt-service",
		Subject:   userID,
		Audience:  []string{"auth-service"},
		ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		IssuedAt:  jwt.NewNumericDate(refreshTokenIss),
	})

	refreshToken, err := refreshClaims.SignedString(j.SignKey)
	if err != nil {
		tg.SendError(err, "api/generate")
		return nil, err
	}

	modelAccessToken := &model.AccessToken{
		Issuer:       "jwt-service",
		Audience:     "auth-service",
		Subject:      userID,
		ExpirationAt: accessTokenExp,
		IssuedAt:     accessTokenIss,
		Token:        accessToken,
	}

	modelRefreshToken := &model.RefreshToken{
		Issuer:       "jwt-service",
		Audience:     "auth-service",
		Subject:      userID,
		ExpirationAt: refreshTokenExp,
		IssuedAt:     refreshTokenIss,
		Token:        refreshToken,
	}

	err = j.JWTRepository.CreateTokens(modelAccessToken, modelRefreshToken)
	if err != nil {
		tg.SendError(err, "api/generate")
		return nil, err
	}

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *jwtService) RefreshTokens(refreshToken string) (*model.Tokens, error) {
	modelRefreshToken, err := j.JWTRepository.GetRefreshToken(refreshToken)
	if err != nil {
		tg.SendError(err, "api/refresh")
		return nil, err
	}

	accessTokenExp := time.Now().Add(j.AccessExpiry)
	accessTokenIss := time.Now()

	refreshTokenExp := time.Now().Add(j.RefreshExpiry)
	refreshTokenIss := time.Now()

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "jwt-service",
		Subject:   modelRefreshToken.Subject,
		Audience:  []string{"auth-service"},
		ExpiresAt: jwt.NewNumericDate(accessTokenExp),
		IssuedAt:  jwt.NewNumericDate(accessTokenIss),
	})

	accessToken, err := accessClaims.SignedString(j.SignKey)
	if err != nil {
		tg.SendError(err, "api/refresh")
		return nil, err
	}

	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "jwt-service",
		Subject:   modelRefreshToken.Subject,
		Audience:  []string{"auth-service"},
		ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		IssuedAt:  jwt.NewNumericDate(refreshTokenIss),
	})

	newRefreshToken, err := refreshClaims.SignedString(j.SignKey)
	if err != nil {
		tg.SendError(err, "api/refresh")
		return nil, err
	}

	modelAccessToken := &model.AccessToken{
		Issuer:       "jwt-service",
		Audience:     "auth-service",
		Subject:      modelRefreshToken.Subject,
		ExpirationAt: accessTokenExp,
		IssuedAt:     accessTokenIss,
		Token:        accessToken,
	}

	modelRefreshToken = &model.RefreshToken{
		Issuer:       "jwt-service",
		Audience:     "auth-service",
		Subject:      modelRefreshToken.Subject,
		ExpirationAt: refreshTokenExp,
		IssuedAt:     refreshTokenIss,
		Token:        newRefreshToken,
	}

	err = j.JWTRepository.CreateTokens(modelAccessToken, modelRefreshToken)
	if err != nil {
		tg.SendError(err, "api/refresh")
		return nil, err
	}

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (j *jwtService) ValidateToken(accessToken string) (*model.AccessToken, error) {
	modelAccessToken, err := j.JWTRepository.GetAccessToken(accessToken)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrAccessTokenNotFound) {
			tg.SendError(err, "api/validate")
			return nil, customServiceError.ErrAccessTokenNotFound
		}
		tg.SendError(err, "api/validate")
		return nil, err
	}

	if modelAccessToken.ExpirationAt.Before(time.Now()) {
		tg.SendError(customServiceError.ErrTokenExpired, "api/validate")
		return nil, customServiceError.ErrTokenExpired
	}

	if modelAccessToken.ExpirationAt.Before(modelAccessToken.IssuedAt) {
		tg.SendError(customServiceError.ErrTokenExpTimeBeforeIss, "api/validate")
		return nil, customServiceError.ErrTokenExpTimeBeforeIss
	}

	return modelAccessToken, nil
}

func newJWTService(signKey []byte, accessExpiry, refreshExpiry time.Duration, JWTRepository repository.JWTRepository) *jwtService {
	return &jwtService{
		SignKey:       signKey,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
		JWTRepository: JWTRepository,
	}
}
