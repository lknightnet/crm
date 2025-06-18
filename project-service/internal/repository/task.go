package repository

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"project-service/internal/model"
	"project-service/internal/repository/customRepositoryError"
	"project-service/pkg/database"
)

type taskRepository struct {
	db *database.PostgreSQL
}

func (t *taskRepository) EditTask(task *model.Task) error {
	return t.db.DB.Model(&model.Task{}).Where("id = ?", task.ID).Updates(task).Error
}

func (t *taskRepository) RemoveAllTaskExecutors(taskID int) error {
	var task model.Task
	if err := t.db.DB.First(&task, taskID).Error; err != nil {
		return err
	}

	if err := t.db.DB.
		Where("task_id = ? AND executor_id != ?", taskID, task.CreatedID).
		Delete(&model.TaskExecutor{}).Error; err != nil {
		return err
	}

	return nil
}

func (t *taskRepository) CreateTask(task *model.Task) (int, error) {
	if err := t.db.DB.Create(task).Error; err != nil {
		return 0, err
	}
	return task.ID, nil
}

func (t *taskRepository) AddTaskExecutors(taskExecutor *model.TaskExecutor) error {
	log.Println(taskExecutor.ID)
	return t.db.DB.Create(taskExecutor).Error
}

func (t *taskRepository) GetTasksByProjectID(projectID int) ([]model.Task, error) {
	var tasks []model.Task
	err := t.db.DB.Where("project_id = ?", projectID).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, customRepositoryError.ErrTaskNotFound
	}
	return tasks, nil
}

func (t *taskRepository) GetTasksByUserID(userID int) ([]model.TaskWithExecutor, error) {
	var tasks []model.Task
	err := t.db.DB.Joins("JOIN task_executors ON task_executors.task_id = tasks.id").
		Where("task_executors.executor_id = ?", userID).
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, customRepositoryError.ErrTaskNotFound
	}

	taskWithExecutor := make([]model.TaskWithExecutor, len(tasks))
	for i, v := range tasks {
		var taskExecutor []model.TaskExecutor
		err = t.db.DB.
			Where("task_id = ?", v.ID).
			Find(&taskExecutor).Error
		if err != nil {
			return nil, err
		}

		taskExecutorWithoutName := make([]model.TaskExecutorWithoutName, len(taskExecutor))
		for i2, executor := range taskExecutor {
			taskExecutorWithoutName[i2] = model.TaskExecutorWithoutName{
				ExecutorID: executor.ExecutorID,
			}
		}
		taskWithExecutor[i] = model.TaskWithExecutor{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Priority:    v.Priority,
			Deadline:    v.Deadline,
			CreatedID:   v.CreatedID,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			ProjectID:   v.ProjectID,
			Status:      v.Status,
			Executors:   taskExecutorWithoutName,
		}
	}
	return taskWithExecutor, nil
}

func (t *taskRepository) GetTaskByID(taskID int) (*model.Task, []model.TaskExecutor, error) {
	var task model.Task
	err := t.db.DB.First(&task, taskID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, customRepositoryError.ErrTaskNotFound
		}
		return nil, nil, err
	}

	var executors []model.TaskExecutor
	err = t.db.DB.Where("task_id = ?", taskID).Find(&executors).Error
	if err != nil {
		return &task, nil, err
	}

	return &task, executors, nil
}

func newTaskRepository(db *database.PostgreSQL) *taskRepository {
	return &taskRepository{db: db}
}
