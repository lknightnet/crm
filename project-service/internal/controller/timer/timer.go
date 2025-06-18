package timer

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"project-service/internal/service"
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
	var json TimerCreateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, TimerErrorResponse{Status: false, Error: err.Error()})
		return
	}

	timer, err := t.TimerService.StartTimer(token.(string), json.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, TimerErrorResponse{Status: false, Error: err.Error()})
	}

	c.JSON(http.StatusOK, NewTimerResponse(timer))
}

func (t *TimerController) StopTimer(c *gin.Context) {
	token, _ := c.Get("token")
	var json TimerCreateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, TimerErrorResponse{Status: false, Error: err.Error()})
		return
	}
	err := t.TimerService.StopTimer(token.(string), json.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, TimerErrorResponse{Status: false, Error: err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (t *TimerController) ResumeTimer(c *gin.Context) {
	id := c.Param("id")
	idI, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, TimerErrorResponse{
			Status: false,
			Error:  "id must be integer",
		})
		return
	}

	timer, err := t.TimerService.ResumeTimer(idI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, TimerErrorResponse{Status: false, Error: err.Error()})
	}

	c.JSON(http.StatusOK, NewTimerResponse(timer))
}
