package repository

import (
	"project-service/internal/model"
	"project-service/pkg/database"
)

type ProjectRepository interface {
	CreateProject(project *model.Project) (int, error)
	GetProjectsByUserID(userID int) ([]model.Project, error)
	EditProject(project *model.Project) error
	GetProjectByID(projectID int) (*model.Project, []model.ProjectUsers, error)
}

type ProjectUsersRepository interface {
	AddProjectUser(projectUser *model.ProjectUsers) error
	RemoveProjectUser(projectUser *model.ProjectUsers) error
	RemoveProjectUsers(projectID int) error
}

type TaskRepository interface {
	CreateTask(task *model.Task) (int, error)
	AddTaskExecutors(taskExecutor *model.TaskExecutor) error
	RemoveAllTaskExecutors(taskID int) error
	GetTasksByProjectID(projectID int) ([]model.Task, error)
	GetTasksByUserID(userID int) ([]model.TaskWithExecutor, error)
	GetTaskByID(taskID int) (*model.Task, []model.TaskExecutor, error)
	EditTask(task *model.Task) error
}
type TimerRepository interface {
	StartTimerEntry(timerEntry *model.TimerEntry) error
	StopTimerEntry(userID int) error
	GetTimersByTaskID(taskID int) ([]model.TimerEntry, error)
	GetTimersByUserID(userID int) ([]model.TimerEntry, error)
}

type InformationListRepository interface {
	GetListByProjectID(projectID int) ([]model.InformationList, error)
	CreateInformationList(informationList *model.InformationList) error
	UpdateInformationList(informationList *model.InformationList) error
	DeleteInformationList(informationListID int) error
}

type Repository struct {
	ProjectRepository         ProjectRepository
	ProjectUsersRepository    ProjectUsersRepository
	TaskRepository            TaskRepository
	TimerRepository           TimerRepository
	InformationListRepository InformationListRepository
}

func NewRepository(db *database.PostgreSQL) *Repository {
	return &Repository{
		ProjectRepository:         newProjectRepository(db),
		ProjectUsersRepository:    newProjectUsersRepository(db),
		TaskRepository:            newTaskRepository(db),
		TimerRepository:           newTimerRepository(db),
		InformationListRepository: newInformationListRepository(db),
	}
}
