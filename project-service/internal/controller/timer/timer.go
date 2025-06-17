package timer

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"project-service/internal/service"
	"project-service/internal/service/customServiceError"
	"strconv"
)

type TimerController struct {
	TimerService service.TimerService
}

func NewTimerController(timerService service.TimerService) *TimerController {
	return &TimerController{
		TimerService: timerService,
	}
}

func (t *TimerController) StartTimer(c *gin.Context) {
	token, _ := c.Get("token")
	id := c.Param("id")
	idI, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, TimerErrorResponse{
			Status: false,
			Error:  "id must be integer",
		})
		return
	}

	err = t.TimerService.StartTimerEntry(token.(string), idI)
	if err != nil {
		if errors.Is(err, customServiceError.ErrAnotherTimerRunning) {
			c.JSON(http.StatusBadRequest, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, TimerErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (t *TimerController) StopTimer(c *gin.Context) {
	token, _ := c.Get("token")

	err := t.TimerService.StopTimerEntry(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTimerNotFound) {
			c.JSON(http.StatusNotFound, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, TimerErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (t *TimerController) GetTimersByTaskID(c *gin.Context) {
	token, _ := c.Get("token")
	id := c.Param("id")
	idI, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, TimerErrorResponse{
			Status: false,
			Error:  "id must be integer",
		})
		return
	}

	timers, err := t.TimerService.GetTimersByTaskID(token.(string), idI)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTimerNotFound) {
			c.JSON(http.StatusNotFound, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, TimerErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	var timerResponse []TimerGetResponse
	for i, timer := range timers {
		timerResponse[i] = *NewTimerGetResponse(&timer)
	}

	c.JSON(http.StatusOK, timerResponse)
}

func (t *TimerController) GetTimersByUserID(c *gin.Context) {
	token, _ := c.Get("token")

	timers, err := t.TimerService.GetTimersByUserID(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTimerNotFound) {
			c.JSON(http.StatusNotFound, TimerErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, TimerErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	var timerResponse []TimerGetResponse
	for i, timer := range timers {
		timerResponse[i] = *NewTimerGetResponse(&timer)
	}

	c.JSON(http.StatusOK, timerResponse)
}
