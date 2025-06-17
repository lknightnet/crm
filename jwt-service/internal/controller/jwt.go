package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"jwt-service/internal/controller/request"
	"jwt-service/internal/controller/response"
	"jwt-service/internal/service"
	"jwt-service/internal/service/customServiceError"
	"net/http"
)

type jwtController struct {
	JWTService service.JWTService
}

func newJWTController(jwtService service.JWTService) *jwtController {
	return &jwtController{
		JWTService: jwtService,
	}
}

// @BasePath /api

// Generate godoc
// @Summary Generate JWT tokens
// @Description This endpoint generates JWT tokens for a user based on the provided user ID.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body request.GenerateRequest true "User ID for token generation"
// @Success 200 {object} response.TokensResponse "Successfully generated tokens"
// @Failure 400 {object} response.ErrorResponse "Invalid input data"
// @Failure 500 {object} response.ErrorResponse "Error generating tokens"
// @Router /api/generate [post]
func (jwtc *jwtController) Generate(c *gin.Context) {
	var json request.GenerateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	tokens, err := jwtc.JWTService.GenerateTokens(json.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Refresh godoc
// @Summary Refresh JWT tokens
// @Description This endpoint refreshes the JWT tokens (access and refresh tokens) based on the provided refresh token.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body request.RefreshRequest true "Refresh token"
// @Success 200 {object} response.TokensResponse "Successfully refreshed tokens"
// @Failure 400 {object} response.ErrorResponse "Invalid input data"
// @Failure 500 {object} response.ErrorResponse "Error generating new tokens"
// @Router /api/refresh [post]
func (jwtc *jwtController) Refresh(c *gin.Context) {
	var json request.RefreshRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	tokens, err := jwtc.JWTService.RefreshTokens(json.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Validate godoc
// @Summary Validate JWT token
// @Description This endpoint validates the provided access token and returns its claims (issuer, audience, subject, etc.).
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body request.ValidateRequest true "Access token"
// @Success 200 {object} response.AccessTokenResponse "Successfully validated token"
// @Failure 400 {object} response.ErrorResponse "Invalid input data"
// @Failure 401 {object} response.ErrorResponse "token is expired"
// @Failure 500 {object} response.ErrorResponse "Error validating token"
// @Router /api/validate [post]
func (jwtc *jwtController) Validate(c *gin.Context) {
	var json request.ValidateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	token, err := jwtc.JWTService.ValidateToken(json.AccessToken)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{Status: false, Error: err.Error()})
			return
		}
		if errors.Is(err, customServiceError.ErrAccessTokenNotFound) {
			c.JSON(http.StatusNotFound, response.ErrorResponse{Status: false, Error: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Status: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.AccessTokenResponse{
		Issuer:       token.Issuer,
		Audience:     token.Audience,
		Subject:      token.Subject,
		ExpirationAt: token.ExpirationAt,
		IssuedAt:     token.IssuedAt,
		Token:        token.Token,
	})
}
