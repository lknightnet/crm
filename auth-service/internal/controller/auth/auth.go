package auth

import (
	"auth-service/internal/controller/auth/request"
	"auth-service/internal/controller/auth/response"
	"auth-service/internal/service"
	"auth-service/internal/service/customServiceError"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	AuthService service.AuthService
}

func NewAuthController(AuthService service.AuthService) *AuthController {
	return &AuthController{
		AuthService: AuthService,
	}
}

func (authc *AuthController) Login(c *gin.Context) {
	var json request.LoginRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tokens, err := authc.AuthService.Login(json.Email, json.Password)
	if err != nil {
		if errors.Is(err, customServiceError.ErrInvalidEmailOrPassword) {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (authc *AuthController) Signup(c *gin.Context) {
	var json request.SignupRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	err := authc.AuthService.SignUp(json.Email, json.Name, json.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.StatusResponse{
		Status: true,
	})
}

func (authc *AuthController) Verify(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: false,
			Error:  "Token is required",
		})
		return
	}

	err := authc.AuthService.Verify(token)
	if err != nil {

		if errors.Is(err, customServiceError.ErrOneTimeTokenExpired) {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Status: false,
				Error:  customServiceError.ErrOneTimeTokenExpired.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.StatusResponse{
		Status: true,
	})
}

func (authc *AuthController) Refresh(c *gin.Context) {
	var json request.RefreshRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tokens, err := authc.AuthService.Refresh(json.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.TokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})

}
