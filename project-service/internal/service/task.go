package service

import (
	"errors"
	"fmt"
	"log"
	"project-service/internal/model"
	"project-service/internal/repository"
	"project-service/internal/repository/customRepositoryError"
	"project-service/internal/service/customServiceError"
	"project-service/pkg/tg"
	"strconv"
	"time"
)

type taskService struct {
	UserServiceApi  string
	UserServicePort string

	TaskRepository    repository.TaskRepository
	ProjectRepository repository.ProjectRepository
}

func (t *taskService) EditTask(taskID int, name *string, description *string, priority *int, deadline *time.Time, projectID *int, executors *[]int) error {
	task := &model.Task{
		ID: taskID,
	}

	if name != nil {
		task.Name = *name
	}
	if description != nil {
		task.Description = description
	}
	if priority != nil {
		task.Priority = *priority
	}
	if deadline != nil {
		task.Deadline = deadline
	}
	if projectID != nil {
		task.ProjectID = *projectID
	}

	err := t.TaskRepository.EditTask(task)
	if err != nil {
		tg.SendError(err.Error(), "api/project/update")
		return customServiceError.ErrUnknownError
	}

	if executors != nil {
		err := t.TaskRepository.RemoveAllTaskExecutors(taskID)
		if err != nil {
			tg.SendError(err, "api/task/edit/deleteExecutors")
			return customServiceError.ErrUnknownError
		}

		// Добавляем новых исполнителей
		for _, executorID := range *executors {
			err := t.AddTaskExecutors(executorID, taskID)
			if err != nil {
				tg.SendError(err, "api/task/edit/addExecutor")
				return customServiceError.ErrUnknownError
			}
		}
	}

	return nil
}

func (t *taskService) GetTaskExpired(token string) ([]model.TaskWithExecutorWithName, error) {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/task/get/expired")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/task/get/expired")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/task/get/expired")
		return nil, customServiceError.ErrUnknownError
	}

	tasksWithExecutor, err := t.TaskRepository.GetTasksByUserID(user.ID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTaskNotFound) {
			tg.SendError(err, "api/task/get/expired")
			return nil, customServiceError.ErrTaskNotFound
		}
		tg.SendError(err, "api/task/get/expired")
		return nil, customServiceError.ErrUnknownError
	}

	now := time.Now()
	var tasksExpired []model.TaskWithExecutorWithName

	for _, task := range tasksWithExecutor {
		if task.Deadline != nil {
			if task.Deadline.Before(now) {
				var tasksUsers []model.TaskExecutorWithName
				for _, executor := range task.Executors {
					user, err := GetUserByID(strconv.Itoa(executor.ExecutorID), t.UserServiceApi, t.UserServicePort)
					if err != nil {
						tg.SendError(fmt.Errorf("method: GetUserByID, error: %s", err.Error()), "api/task/get/expired")
						return nil, customServiceError.ErrUnknownError
					}

					tasksUsers = append(tasksUsers, model.TaskExecutorWithName{
						ExecutorID:   executor.ExecutorID,
						ExecutorName: user.Name,
					})
				}

				project, _, err := t.ProjectRepository.GetProjectByID(task.ProjectID)
				if err != nil {
					if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
						tg.SendError(err, "api/task/get/expired")
						return nil, customServiceError.ErrProjectNotFound
					}
					tg.SendError(err, "api/task/get/expired")
					return nil, customServiceError.ErrUnknownError
				}

				taskWithName := model.TaskWithExecutorWithName{
					ID:          task.ID,
					Name:        task.Name,
					Description: task.Description,
					Priority:    task.Priority,
					Deadline:    task.Deadline,
					CreatedID:   task.CreatedID,
					CreatedAt:   task.CreatedAt,
					UpdatedAt:   task.UpdatedAt,
					ProjectID:   task.ProjectID,
					ProjectName: project.Name,
					Status:      task.Status,
					Executors:   tasksUsers,
				}
				tasksExpired = append(tasksExpired, taskWithName)
			}
		}

	}

	return tasksExpired, nil
}

func (t *taskService) GetTaskToday(token string) ([]model.TaskWithExecutorWithName, error) {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/task/get/today")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/task/get/today")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/task/get/today")
		return nil, customServiceError.ErrUnknownError
	}

	tasks, err := t.TaskRepository.GetTasksByUserID(user.ID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTaskNotFound) {
			tg.SendError(err, "api/task/get/today")
			return nil, customServiceError.ErrTaskNotFound
		}
		tg.SendError(err, "api/task/get/today")
		return nil, customServiceError.ErrUnknownError
	}

	now := time.Now()
	var tasksExpired []model.TaskWithExecutorWithName

	for _, task := range tasks {
		if task.Deadline != nil {
			if task.Deadline.Day() == now.Day() {
				var tasksUsers []model.TaskExecutorWithName
				for _, executor := range task.Executors {
					executorUser, err := GetUserByID(strconv.Itoa(executor.ExecutorID), t.UserServiceApi, t.UserServicePort)
					if err != nil {
						tg.SendError(err, "api/task/get/today")
						return nil, customServiceError.ErrUnknownError
					}

					tasksUsers = append(tasksUsers, model.TaskExecutorWithName{
						ExecutorID:   executor.ExecutorID,
						ExecutorName: executorUser.Name,
					})
				}

				project, _, err := t.ProjectRepository.GetProjectByID(task.ProjectID)
				if err != nil {
					if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
						tg.SendError(err, "api/task/get/today")
						return nil, customServiceError.ErrProjectNotFound
					}
					tg.SendError(err, "api/task/get/today")
					return nil, customServiceError.ErrUnknownError
				}

				taskWithName := model.TaskWithExecutorWithName{
					ID:          task.ID,
					Name:        task.Name,
					Description: task.Description,
					Priority:    task.Priority,
					Deadline:    task.Deadline,
					CreatedID:   task.CreatedID,
					CreatedAt:   task.CreatedAt,
					UpdatedAt:   task.UpdatedAt,
					ProjectID:   task.ProjectID,
					ProjectName: project.Name,
					Status:      task.Status,
					Executors:   tasksUsers,
				}
				tasksExpired = append(tasksExpired, taskWithName)
			}
		}
	}
	return tasksExpired, nil
}

func (t *taskService) GetTaskThisWeek(token string) ([]model.TaskWithExecutorWithName, error) {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/task/get/thisweek")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/task/get/thisweek")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/task/get/thisweek")
		return nil, customServiceError.ErrUnknownError
	}

	tasks, err := t.TaskRepository.GetTasksByUserID(user.ID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTaskNotFound) {
			tg.SendError(err, "api/task/get/thisweek")
			return nil, customServiceError.ErrTaskNotFound
		}
		tg.SendError(err, "api/task/get/thisweek")
		return nil, customServiceError.ErrUnknownError
	}

	now := time.Now()

	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())

	// Конец недели (воскресенье 23:59:59)
	endOfWeek := startOfWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

	var tasksThisWeek []model.TaskWithExecutorWithName

	for _, task := range tasks {
		if task.Deadline != nil {
			if task.Deadline.After(startOfWeek) && task.Deadline.Before(endOfWeek) {
				log.Println("сервисы, эта неделя таск айди")
				log.Println(task.ID)
				var tasksUsers []model.TaskExecutorWithName
				for _, executor := range task.Executors {
					executorUser, err := GetUserByID(strconv.Itoa(executor.ExecutorID), t.UserServiceApi, t.UserServicePort)
					if err != nil {
						tg.SendError(err, "api/task/get/thisweek")
						return nil, customServiceError.ErrUnknownError
					}

					tasksUsers = append(tasksUsers, model.TaskExecutorWithName{
						ExecutorID:   executor.ExecutorID,
						ExecutorName: executorUser.Name,
					})
				}

				project, _, err := t.ProjectRepository.GetProjectByID(task.ProjectID)
				if err != nil {
					if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
						tg.SendError(err, "api/task/get/thisweek")
						return nil, customServiceError.ErrProjectNotFound
					}
					tg.SendError(err, "api/task/get/thisweek")
					return nil, customServiceError.ErrUnknownError
				}

				taskWithName := model.TaskWithExecutorWithName{
					ID:          task.ID,
					Name:        task.Name,
					Description: task.Description,
					Priority:    task.Priority,
					Deadline:    task.Deadline,
					CreatedID:   task.CreatedID,
					CreatedAt:   task.CreatedAt,
					UpdatedAt:   task.UpdatedAt,
					ProjectName: project.Name,
					ProjectID:   task.ProjectID,
					Status:      task.Status,
					Executors:   tasksUsers,
				}
				tasksThisWeek = append(tasksThisWeek, taskWithName)
			}
		}
	}
	log.Println("сервисы, эта неделя таски которые выше падают")
	log.Println(tasksThisWeek)
	return tasksThisWeek, nil
}

func (t *taskService) GetTaskNextWeek(token string) ([]model.TaskWithExecutorWithName, error) {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/task/get/nextweek")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/task/get/nextweek")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/task/get/nextweek")
		return nil, customServiceError.ErrUnknownError
	}

	tasks, err := t.TaskRepository.GetTasksByUserID(user.ID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTaskNotFound) {
			tg.SendError(err, "api/task/get/nextweek")
			return nil, customServiceError.ErrTaskNotFound
		}
		tg.SendError(err, "api/task/get/nextweek")
		return nil, customServiceError.ErrUnknownError
	}

	now := time.Now()

	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfThisWeek := time.Date(
		now.Year(), now.Month(), now.Day()-weekday+1,
		0, 0, 0, 0,
		now.Location(),
	)

	startOfNextWeek := startOfThisWeek.AddDate(0, 0, 7)
	endOfNextWeek := startOfNextWeek.AddDate(0, 0, 7).Add(-time.Nanosecond)

	var tasksNextWeek []model.TaskWithExecutorWithName
	for _, task := range tasks {
		deadline := task.Deadline
		if deadline != nil {
			if (deadline.Equal(startOfNextWeek) || deadline.After(startOfNextWeek)) &&
				(deadline.Before(endOfNextWeek) || deadline.Equal(endOfNextWeek)) {
				log.Println("сервисы, дедлайн неделя след")
				log.Println(deadline)
				var tasksUsers []model.TaskExecutorWithName
				for _, executor := range task.Executors {
					executorUser, err := GetUserByID(strconv.Itoa(executor.ExecutorID), t.UserServiceApi, t.UserServicePort)
					if err != nil {
						tg.SendError(err, "api/task/get/nextweek")
						return nil, customServiceError.ErrUnknownError
					}

					tasksUsers = append(tasksUsers, model.TaskExecutorWithName{
						ExecutorID:   executor.ExecutorID,
						ExecutorName: executorUser.Name,
					})
				}
				project, _, err := t.ProjectRepository.GetProjectByID(task.ProjectID)
				if err != nil {
					if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
						tg.SendError(err, "api/task/get/nextweek")
						return nil, customServiceError.ErrProjectNotFound
					}
					tg.SendError(err, "api/task/get/nextweek")
					return nil, customServiceError.ErrUnknownError
				}

				taskWithName := model.TaskWithExecutorWithName{
					ID:          task.ID,
					Name:        task.Name,
					Description: task.Description,
					Priority:    task.Priority,
					Deadline:    task.Deadline,
					CreatedID:   task.CreatedID,
					CreatedAt:   task.CreatedAt,
					UpdatedAt:   task.UpdatedAt,
					ProjectName: project.Name,
					ProjectID:   task.ProjectID,
					Status:      task.Status,
					Executors:   tasksUsers,
				}
				tasksNextWeek = append(tasksNextWeek, taskWithName)
			}
		}

	}
	log.Println("сервисы, след неделя какие таски идут выше")
	log.Println(tasksNextWeek)
	return tasksNextWeek, nil
}

func (t *taskService) GetTaskNotDeadline(token string) ([]model.TaskWithExecutorWithName, error) {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/task/get/notdeadline")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/task/get/notdeadline")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/task/get/notdeadline")
		return nil, customServiceError.ErrUnknownError
	}

	tasks, err := t.TaskRepository.GetTasksByUserID(user.ID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTaskNotFound) {
			tg.SendError(err, "api/task/get/notdeadline")
			return nil, customServiceError.ErrTaskNotFound
		}
		tg.SendError(err, "api/task/get/notdeadline")
		return nil, customServiceError.ErrUnknownError
	}

	var tasksNextWeek []model.TaskWithExecutorWithName
	for _, task := range tasks {
		deadline := task.Deadline
		if deadline == nil {
			log.Println("сервисы, таски, которые не имеют дедлайна")
			log.Println(task)
			var tasksUsers []model.TaskExecutorWithName
			for _, executor := range task.Executors {
				executorUser, err := GetUserByID(strconv.Itoa(executor.ExecutorID), t.UserServiceApi, t.UserServicePort)
				if err != nil {
					tg.SendError(err, "api/task/get/nextweek")
					return nil, customServiceError.ErrUnknownError
				}

				tasksUsers = append(tasksUsers, model.TaskExecutorWithName{
					ExecutorID:   executor.ExecutorID,
					ExecutorName: executorUser.Name,
				})
			}

			project, _, err := t.ProjectRepository.GetProjectByID(task.ProjectID)
			if err != nil {
				if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
					tg.SendError(err, "api/task/get/nextweek")
					return nil, customServiceError.ErrProjectNotFound
				}
				tg.SendError(err, "api/task/get/nextweek")
				return nil, customServiceError.ErrUnknownError
			}

			taskWithName := model.TaskWithExecutorWithName{
				ID:          task.ID,
				Name:        task.Name,
				Description: task.Description,
				Priority:    task.Priority,
				Deadline:    task.Deadline,
				CreatedID:   task.CreatedID,
				CreatedAt:   task.CreatedAt,
				ProjectName: project.Name,
				UpdatedAt:   task.UpdatedAt,
				ProjectID:   task.ProjectID,
				Status:      task.Status,
				Executors:   tasksUsers,
			}
			tasksNextWeek = append(tasksNextWeek, taskWithName)
		}
	}

	log.Println("сервисы все таски, которые идут дальше")
	log.Println(tasksNextWeek)
	return tasksNextWeek, nil
}

func (t *taskService) CreateTask(token string, name string, description *string, levelPriority int, deadline *time.Time, projectID int, executors []int) error {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/task/create")
			return customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/task/create")
			return customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/task/create")
		return customServiceError.ErrUnknownError
	}

	task := &model.Task{
		Name:      name,
		Priority:  levelPriority,
		CreatedID: user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ProjectID: projectID,
		Status:    model.STATUS_NOT_STARTED,
	}

	if description != nil {
		task.Description = description
	}
	if deadline != nil {
		task.Deadline = deadline
	}

	taskID, err := t.TaskRepository.CreateTask(task)
	if err != nil {
		tg.SendError(err, "api/task/create")
		return customServiceError.ErrUnknownError
	}

	if executors != nil {
		for _, executor := range executors {
			err := t.AddTaskExecutors(executor, taskID)
			if err != nil {
				tg.SendError(err, "api/task/create")
				return customServiceError.ErrUnknownError
			}
		}
	}

	return nil

}

func (t *taskService) AddTaskExecutors(executorID, taskID int) error {
	taskExecutor := &model.TaskExecutor{
		TaskID:     taskID,
		ExecutorID: executorID,
	}
	err := t.TaskRepository.AddTaskExecutors(taskExecutor)
	if err != nil {
		tg.SendError(err, "api/task/executor/add")
		return customServiceError.ErrUnknownError
	}
	return nil
}

func (t *taskService) GetTasksByProjectID(projectID int) ([]model.Task, error) {
	tasks, err := t.TaskRepository.GetTasksByProjectID(projectID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTaskNotFound) {
			tg.SendError(err, "api/task/get/project/:id")
			return nil, customServiceError.ErrTaskNotFound
		}
		tg.SendError(err, "api/task/get/project/:id")
		return nil, customServiceError.ErrUnknownError
	}
	return tasks, nil
}

func (t *taskService) GetTasksByToken(token string) ([]model.TaskWithExecutorWithName, error) {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/task/get/user")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/task/get/user")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/task/get/user")
		return nil, customServiceError.ErrUnknownError
	}

	tasks, err := t.TaskRepository.GetTasksByUserID(user.ID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTaskNotFound) {
			tg.SendError(err, "api/task/get/user")
			return nil, customServiceError.ErrTaskNotFound
		}
		tg.SendError(err, "api/task/get/user")
		return nil, customServiceError.ErrUnknownError
	}

	var taskss []model.TaskWithExecutorWithName
	for _, task := range tasks {
		var tasksUsers []model.TaskExecutorWithName
		for _, executor := range task.Executors {
			executorUser, err := GetUserByID(strconv.Itoa(executor.ExecutorID), t.UserServiceApi, t.UserServicePort)
			if err != nil {
				tg.SendError(err, "api/task/get/user")
				return nil, customServiceError.ErrUnknownError
			}

			tasksUsers = append(tasksUsers, model.TaskExecutorWithName{
				ExecutorID:   executor.ExecutorID,
				ExecutorName: executorUser.Name,
			})
		}

		project, _, err := t.ProjectRepository.GetProjectByID(task.ProjectID)
		if err != nil {
			if errors.Is(err, customRepositoryError.ErrProjectNotFound) {
				tg.SendError(err, "api/task/get/user")
				return nil, customServiceError.ErrProjectNotFound
			}
			tg.SendError(err, "api/task/get/user")
			return nil, customServiceError.ErrUnknownError
		}

		taskWithName := model.TaskWithExecutorWithName{
			ID:          task.ID,
			Name:        task.Name,
			Description: task.Description,
			Priority:    task.Priority,
			Deadline:    task.Deadline,
			CreatedID:   task.CreatedID,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			ProjectID:   task.ProjectID,
			ProjectName: project.Name,
			Status:      task.Status,
			Executors:   tasksUsers,
		}
		taskss = append(taskss, taskWithName)
	}

	return taskss, nil
}

func (t *taskService) GetTaskByID(token string, taskID int) (*model.Task, []model.TaskExecutor, error) {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/task/get/single/:id")
			return nil, nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/task/get/single/:id")
			return nil, nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/task/get/user")
		return nil, nil, customServiceError.ErrUnknownError
	}

	task, executors, err := t.TaskRepository.GetTaskByID(taskID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTaskNotFound) {
			tg.SendError(err, "api/task/get/single/:id")
			return nil, nil, customServiceError.ErrTaskNotFound
		}
	}

	isFound := false
	for _, executor := range executors {
		if executor.ExecutorID == user.ID {
			isFound = true
			break
		}
	}

	if isFound {
		return task, executors, nil
	} else {
		return nil, nil, customServiceError.ErrPermissionDenied
	}
}

func newTaskService(taskRepository repository.TaskRepository, projectRepository repository.ProjectRepository, userServiceApi, userServicePort string) *taskService {
	return &taskService{
		UserServiceApi:    userServiceApi,
		UserServicePort:   userServicePort,
		TaskRepository:    taskRepository,
		ProjectRepository: projectRepository,
	}
}
