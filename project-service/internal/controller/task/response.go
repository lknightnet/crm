package task

import (
	"project-service/internal/model"
	"time"
)

type TaskErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

type TaskWithExecutorsWithNameResponse struct {
	ID          int                          `json:"id"`
	Name        string                       `json:"name"`
	Description *string                      `json:"description"`
	Status      string                       `json:"status"`
	Priority    int                          `json:"priority"`
	Deadline    *time.Time                   `json:"deadline"`
	CreatedID   int                          `json:"created_id"` //кто создал
	CreatedAt   time.Time                    `json:"created_at"`
	Executors   []model.TaskExecutorWithName `json:"executors"`
	ProjectID   int                          `json:"project_id"`
	ProjectName string                       `json:"project_name"`
}

type TaskExecutorResponse struct {
	ExecutorID int `json:"executor_id"`
}

type TaskResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	CreatedID   int        `json:"created_id"` //кто создал
	CreatedAt   time.Time  `json:"created_at"`
}

type TaskWithExecutorResponse struct {
	TaskResponse *TaskResponse          `json:"task"`
	TaskExecutor []TaskExecutorResponse `json:"task_executor"`
}

type TaskProjectResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	CreatedID   int        `json:"created_id"` //кто создал
	CreatedAt   time.Time  `json:"created_at"`
}

func NewTaskResponse(task *model.TaskWithExecutorWithName) *TaskWithExecutorsWithNameResponse {
	return &TaskWithExecutorsWithNameResponse{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		Deadline:    task.Deadline,
		CreatedID:   task.CreatedID,
		ProjectName: task.ProjectName,
		ProjectID:   task.ProjectID,
		CreatedAt:   task.CreatedAt,
		Executors:   task.Executors,
	}
}

func NewTasksResponse(tasks []model.TaskWithExecutorWithName) []TaskWithExecutorsWithNameResponse {
	tasksGetByProjectID := make([]TaskWithExecutorsWithNameResponse, len(tasks))
	for i, task := range tasks {
		tasksGetByProjectID[i] = *NewTaskResponse(&task)
	}
	return tasksGetByProjectID
}

func NewTaskWithExecutorResponse(task *model.Task, executors []model.TaskExecutor) *TaskWithExecutorResponse {

	taskUsers := make([]TaskExecutorResponse, len(executors))
	for i, user := range executors {
		taskUsers[i] = TaskExecutorResponse{
			ExecutorID: user.ExecutorID,
		}
	}
	return &TaskWithExecutorResponse{
		TaskResponse: (*TaskResponse)(NewTaskkResponse(task)),
		TaskExecutor: taskUsers,
	}
}

func NewTaskkResponse(task *model.Task) *TaskProjectResponse {
	return &TaskProjectResponse{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		Deadline:    task.Deadline,
		CreatedID:   task.CreatedID,
		CreatedAt:   task.CreatedAt,
	}
}
func NewTasksProjectResponse(tasks []model.Task) []TaskProjectResponse {
	tasksGetByProjectID := make([]TaskProjectResponse, len(tasks))
	for i, task := range tasks {
		tasksGetByProjectID[i] = *NewTaskkResponse(&task)
	}
	return tasksGetByProjectID
}
