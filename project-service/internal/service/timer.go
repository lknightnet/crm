package service

import (
	"errors"
	"project-service/internal/model"
	"project-service/internal/repository"
	"project-service/internal/repository/customRepositoryError"
	"project-service/internal/service/customServiceError"
	"project-service/pkg/tg"
	"time"
)

type timerService struct {
	UserServiceApi  string
	UserServicePort string
	TimerRepository repository.TimerRepository
}

func (t *timerService) StartTimerEntry(token string, taskID int) error {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/timer/start/:id")
			return customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/timer/start/:id")
			return customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/timer/start/:id")
		return customServiceError.ErrUnknownError
	}

	timers, err := t.TimerRepository.GetTimersByUserID(user.ID)
	if err != nil {
		tg.SendError(err, "api/timer/start/:id")
		return customServiceError.ErrUnknownError
	}

	for _, timer := range timers {
		if timer.Active {
			return customServiceError.ErrAnotherTimerRunning
		}
	}

	timerEntry := &model.TimerEntry{
		TaskID:    taskID,
		UserID:    user.ID,
		StartAt:   time.Now(),
		CreatedAt: time.Now(),
		Active:    true,
	}

	err = t.TimerRepository.StartTimerEntry(timerEntry)
	if err != nil {
		tg.SendError(err, "api/timer/start/:id")
		return customServiceError.ErrUnknownError
	}
	return nil
}

func (t *timerService) StopTimerEntry(token string) error {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/timer/stop")
			return customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/timer/stop")
			return customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/timer/stop")
		return customServiceError.ErrUnknownError
	}

	err = t.TimerRepository.StopTimerEntry(user.ID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTimerNotFound) {
			tg.SendError(err, "api/timer/stop")
			return customServiceError.ErrTimerNotFound
		}
		tg.SendError(err, "api/timer/stop")
		return customServiceError.ErrUnknownError
	}
	return nil
}

func (t *timerService) GetTimersByTaskID(token string, taskID int) ([]model.TimerEntry, error) {
	_, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/timer/get/task/:id")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/timer/get/task/:id")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/timer/get/task/:id")
		return nil, customServiceError.ErrUnknownError
	}

	timers, err := t.TimerRepository.GetTimersByTaskID(taskID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTimerNotFound) {
			return nil, customServiceError.ErrTimerNotFound
		}
		tg.SendError(err, "api/timer/get/task/:id")
		return nil, customServiceError.ErrUnknownError
	}
	return timers, err
}

func (t *timerService) GetTimersByUserID(token string) ([]model.TimerEntry, error) {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err, "api/timer/get/user")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err, "api/timer/get/user")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err, "api/timer/get/user")
		return nil, customServiceError.ErrUnknownError
	}

	timers, err := t.TimerRepository.GetTimersByUserID(user.ID)
	if err != nil {
		if errors.Is(err, customRepositoryError.ErrTimerNotFound) {
			return nil, customServiceError.ErrTimerNotFound
		}
		tg.SendError(err, "api/timer/get/user")
		return nil, customServiceError.ErrUnknownError
	}
	return timers, err
}

func newTimerService(timerRepository repository.TimerRepository, userServiceApi, userServicePort string) *timerService {
	return &timerService{
		TimerRepository: timerRepository,
		UserServiceApi:  userServiceApi,
		UserServicePort: userServicePort,
	}
}
