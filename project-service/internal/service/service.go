package service

import (
	"project-service/internal/model"
	"project-service/internal/repository"
	"time"
)

type ProjectService interface {
	CreateProject(token string, name string, description *string, projectUsers *[]int) error
	GetProjectsByToken(token string) ([]model.ProjectWithCreatedUsername, error)
	GetProjectByID(token string, projectID int) (*model.Project, []model.ProjectUsers, error)
	EditProject(projectID int, name, description *string, projectUsers *[]int) error
	GetProjectsByName(token string, projectName string) ([]model.ProjectWithProjectUsers, error)
}

type ProjectUsersService interface {
	AddProjectUser(userID, projectID int) error
	RemoveProjectUser(userID, projectID int) error
}

type TaskService interface {
	CreateTask(token string, name string, description *string, levelPriority int, deadline *time.Time, projectID int, executors []int) error
	AddTaskExecutors(executorID, taskID int) error
	GetTasksByProjectID(projectID int) ([]model.Task, error)
	GetTasksByToken(token string) ([]model.TaskWithExecutorWithName, error)
	GetTaskByID(token string, taskID int) (*model.Task, []model.TaskExecutor, error)
	EditTask(taskID int, name *string, description *string, priority *int, deadline *time.Time, projectID *int, executors *[]int) error

	GetTaskExpired(token string) ([]model.TaskWithExecutorWithName, error)
	GetTaskToday(token string) ([]model.TaskWithExecutorWithName, error)
	GetTaskThisWeek(token string) ([]model.TaskWithExecutorWithName, error)
	GetTaskNextWeek(token string) ([]model.TaskWithExecutorWithName, error)
	GetTaskNotDeadline(token string) ([]model.TaskWithExecutorWithName, error)
}

type TimerService interface {
	StartTimer(token string, taskID *int) (*model.TimerEntry, error)
	StopTimer(token string, taskID *int) error
	ResumeTimer(timerID int) (*model.TimerEntry, error)
	GetTimersByTask(taskID int) ([]model.TimerEntry, error)
	GetTimersByUser(userID int) ([]model.TimerEntry, error)
}

type InformationListService interface {
	GetListByProjectID(projectID int) ([]model.InformationList, error)
	CreateInformationList(projectID int, key, value string) error
	UpdateInformationList(listID int, key, value *string) error
	DeleteInformationList(listID int) error
}

type Service struct {
	ProjectService         ProjectService
	ProjectUsersService    ProjectUsersService
	TaskService            TaskService
	TimerService           TimerService
	InformationListService InformationListService
}

type DependenciesService struct {
	UserServiceApi  string
	UserServicePort string

	ProjectRepository         repository.ProjectRepository
	ProjectUsersRepository    repository.ProjectUsersRepository
	TaskRepository            repository.TaskRepository
	TimerRepository           repository.TimerRepository
	InformationListRepository repository.InformationListRepository
}

func NewService(deps *DependenciesService) *Service {
	return &Service{
		ProjectService:         newProjectService(deps.ProjectRepository, deps.ProjectUsersRepository, deps.UserServiceApi, deps.UserServicePort),
		ProjectUsersService:    newProjectUserService(deps.ProjectUsersRepository),
		TaskService:            newTaskService(deps.TaskRepository, deps.ProjectRepository, deps.UserServiceApi, deps.UserServicePort),
		TimerService:           newTimerService(deps.TimerRepository, deps.UserServiceApi, deps.UserServicePort),
		InformationListService: newInformationList(deps.InformationListRepository),
	}
}
