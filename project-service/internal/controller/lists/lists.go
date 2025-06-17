package lists

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"project-service/internal/service"
	"project-service/internal/service/customServiceError"
	"strconv"
)

type InformationListController struct {
	InformationListService service.InformationListService
}

func NewInformationListController(informationListService service.InformationListService) *InformationListController {
	return &InformationListController{
		InformationListService: informationListService,
	}
}

func (i *InformationListController) GetListByProjectID(c *gin.Context) {
	id := c.Param("id")
	projectID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ListErrorResponse{
			Status: true,
			Error:  err.Error(),
		})
		return
	}

	lists, err := i.InformationListService.GetListByProjectID(projectID)
	if err != nil {
		if errors.Is(err, customServiceError.ErrInformationListNotFound) {
			c.JSON(http.StatusBadRequest, ListErrorResponse{
				Status: true,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ListErrorResponse{
			Status: true,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, lists)
}

func (i *InformationListController) CreateInformationList(c *gin.Context) {
	var json ListRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, ListErrorResponse{Status: false, Error: err.Error()})
		return
	}

	err := i.InformationListService.CreateInformationList(json.ProjectID, *json.Key, *json.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ListErrorResponse{
			Status: true,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (i *InformationListController) UpdateInformationList(c *gin.Context) {
	var json ListRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, ListErrorResponse{Status: false, Error: err.Error()})
		return
	}

	err := i.InformationListService.UpdateInformationList(json.ListID, json.Key, json.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ListErrorResponse{
			Status: true,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (i *InformationListController) DeleteInformationList(c *gin.Context) {
	id := c.Param("id")
	listID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ListErrorResponse{
			Status: true,
			Error:  err.Error(),
		})
		return
	}

	err = i.InformationListService.DeleteInformationList(listID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ListErrorResponse{
			Status: true,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}
