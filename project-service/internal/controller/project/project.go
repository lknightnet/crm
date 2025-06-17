package project

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"project-service/internal/service"
	"project-service/internal/service/customServiceError"
	"strconv"
)

type ProjectController struct {
	ProjectService service.ProjectService
}

func NewProjectController(projectService service.ProjectService) *ProjectController {
	return &ProjectController{
		ProjectService: projectService,
	}
}

func (p *ProjectController) CreateProject(c *gin.Context) {
	token, _ := c.Get("token")
	var json ProjectCreateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, ProjectErrorResponse{Status: false, Error: err.Error()})
		return
	}

	err := p.ProjectService.CreateProject(token.(string), json.Name, json.Description, json.Users)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, ProjectErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, ProjectErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ProjectErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (p *ProjectController) GetProjectsByToken(c *gin.Context) {
	token, _ := c.Get("token")

	projects, err := p.ProjectService.GetProjectsByToken(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, ProjectErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, ProjectErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrProjectsNotFound) {
			c.JSON(http.StatusInternalServerError, ProjectErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ProjectErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	projectsResponse := NewProjectsResponse(projects)

	c.JSON(http.StatusOK, projectsResponse)
}

func (p *ProjectController) GetProjectByID(c *gin.Context) {
	token, _ := c.Get("token")
	id := c.Param("id")

	idI, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ProjectErrorResponse{
			Status: false,
			Error:  "id must be integer",
		})
		return
	}

	project, users, err := p.ProjectService.GetProjectByID(token.(string), idI)
	if err != nil {
		if errors.Is(err, customServiceError.ErrProjectNotFound) {
			c.JSON(http.StatusInternalServerError, ProjectErrorResponse{
				Status: false,
				Error:  customServiceError.ErrProjectNotFound.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, ProjectErrorResponse{
				Status: false,
				Error:  customServiceError.ErrProjectNotFound.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, ProjectErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, ProjectErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ProjectErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	projectsResponse := NewProjectWithUsersResponse(project, users)

	c.JSON(http.StatusOK, projectsResponse)
}

func (p *ProjectController) EditProject(c *gin.Context) {
	var json ProjectUpdateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, ProjectErrorResponse{Status: false, Error: err.Error()})
		return
	}

	err := p.ProjectService.EditProject(json.ProjectID, json.Name, json.Description, json.Users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ProjectErrorResponse{
			Status: true,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}
