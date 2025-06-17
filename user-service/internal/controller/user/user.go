package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-service/internal/controller/user/request"
	"user-service/internal/controller/user/response"
	"user-service/internal/service"
	"user-service/internal/service/customServiceError"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(UserService service.UserService) *UserController {
	return &UserController{UserService: UserService}
}

func (userc *UserController) Get(c *gin.Context) {
	token, _ := c.Get("token")

	log.Println("token", token)

	user, err := userc.UserService.GetUser(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, response.ErrorResponse{
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

	c.JSON(http.StatusOK, response.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		LastLogin: user.LastLogin,
		CreatedAt: user.CreatedAt,
	})
}

func (userc *UserController) Validate(c *gin.Context) {
	token, _ := c.Get("token")

	err := userc.UserService.ValidateToken(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenNotValid) {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{
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

	c.JSON(http.StatusOK, response.StatusResponse{
		Status: true,
	})
}

func (userc *UserController) OutsideGetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := userc.UserService.GetUserById(id)
	if err != nil {
		if errors.Is(err, customServiceError.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, response.ErrorResponse{
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

	c.JSON(http.StatusOK, user)
}

func (userc *UserController) GetUsersByUsername(c *gin.Context) {
	var json request.GetUserByUserNameRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	users, err := userc.UserService.GetUsersByUsername(json.Username)
	if err != nil {
		if errors.Is(err, customServiceError.ErrUsersNotFound) {
			c.JSON(http.StatusOK, response.MessageResponse{
				Status:  true,
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	usersResponse := response.NewUsersByUserNameResponse(users)

	c.JSON(http.StatusOK, usersResponse)
}

func (userc *UserController) UpdateUserByID(c *gin.Context) {
	token, _ := c.Get("token")
	var json request.UpdateUser
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	err := userc.UserService.UpdateUserByID(token.(string), json.Name, json.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}
