package task

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"project-service/internal/service"
	"project-service/internal/service/customServiceError"
	"strconv"
)

type TaskController struct {
	TaskService service.TaskService
}

func NewTaskController(taskService service.TaskService) *TaskController {
	return &TaskController{
		TaskService: taskService,
	}
}

func (t *TaskController) CreateTask(c *gin.Context) {
	token, _ := c.Get("token")
	var json TaskCreateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, TaskErrorResponse{Status: false, Error: err.Error()})
		return
	}

	err := t.TaskService.CreateTask(token.(string), json.Name, json.Description, json.LevelPriority, json.Deadline, json.ProjectID, json.Executors)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (t *TaskController) GetTasksByProjectID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, TaskErrorResponse{
			Status: false,
			Error:  "id must be integer",
		})
		return
	}
	tasks, err := t.TaskService.GetTasksByProjectID(idInt)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tasksResponse := NewTasksProjectResponse(tasks)

	c.JSON(http.StatusOK, tasksResponse)
}

func (t *TaskController) GetTasksByToken(c *gin.Context) {
	token, _ := c.Get("token")

	tasks, err := t.TaskService.GetTasksByToken(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tasksResponse := NewTasksResponse(tasks)

	c.JSON(http.StatusOK, tasksResponse)
}

func (t *TaskController) GetTaskByID(c *gin.Context) {
	token, _ := c.Get("token")
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, TaskErrorResponse{
			Status: false,
			Error:  "id must be integer",
		})
		return
	}

	task, executors, err := t.TaskService.GetTaskByID(token.(string), idInt)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tasksResponse := NewTaskWithExecutorResponse(task, executors)

	c.JSON(http.StatusOK, tasksResponse)
}

func (t *TaskController) GetTaskExpired(c *gin.Context) {
	token, _ := c.Get("token")

	tasks, err := t.TaskService.GetTaskExpired(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tasksResponse := NewTasksResponse(tasks)

	c.JSON(http.StatusOK, tasksResponse)
}

func (t *TaskController) GetTaskToday(c *gin.Context) {
	token, _ := c.Get("token")

	tasks, err := t.TaskService.GetTaskToday(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tasksResponse := NewTasksResponse(tasks)

	c.JSON(http.StatusOK, tasksResponse)
}

func (t *TaskController) GetTaskThisWeek(c *gin.Context) {
	token, _ := c.Get("token")

	tasks, err := t.TaskService.GetTaskThisWeek(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tasksResponse := NewTasksResponse(tasks)

	c.JSON(http.StatusOK, tasksResponse)
}

func (t *TaskController) GetTaskNextWeek(c *gin.Context) {
	token, _ := c.Get("token")

	tasks, err := t.TaskService.GetTaskNextWeek(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tasksResponse := NewTasksResponse(tasks)

	c.JSON(http.StatusOK, tasksResponse)
}

func (t *TaskController) GetTaskNotDeadline(c *gin.Context) {
	token, _ := c.Get("token")

	tasks, err := t.TaskService.GetTaskNotDeadline(token.(string))
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		if errors.Is(err, customServiceError.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, TaskErrorResponse{
				Status: false,
				Error:  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: false,
			Error:  err.Error(),
		})
		return
	}

	tasksResponse := NewTasksResponse(tasks)

	c.JSON(http.StatusOK, tasksResponse)
}

func (t *TaskController) EditTask(c *gin.Context) {
	var json TaskUpdateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, TaskErrorResponse{Status: false, Error: err.Error()})
		return
	}

	err := t.TaskService.EditTask(json.TaskID, json.Name, json.Description, json.LevelPriority, json.Deadline, json.ProjectID, json.Executors)
	if err != nil {
		c.JSON(http.StatusInternalServerError, TaskErrorResponse{
			Status: true,
			Error:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}
