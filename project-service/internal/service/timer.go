package service

import (
	"errors"
	"log"
	"project-service/internal/model"
	"project-service/internal/repository"
	"project-service/internal/service/customServiceError"
	"project-service/pkg/tg"
)

type timerService struct {
	UserServiceApi  string
	UserServicePort string
	TimerRepository repository.TimerRepository
}

func (t *timerService) StartTimer(token string, taskID *int) (*model.TimerEntry, error) {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err.Error(), "api/timer/start")
			return nil, customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err.Error(), "api/timer/start")
			return nil, customServiceError.ErrTokenNotFound
		}
		tg.SendError(err.Error(), "api/timer/get/like")
		log.Println("GetUserByToken, api/timer/start, error:", err)
		return nil, customServiceError.ErrUnknownError
	}

	timer := &model.TimerEntry{
		UserID: user.ID,
	}

	if taskID != nil {
		timer.TaskID = taskID
	}

	err = t.TimerRepository.StartTimerEntry(timer)
	if err != nil {
		log.Println("StartTimerEntry, api/timer/start, error:", err)
		return nil, customServiceError.ErrUnknownError
	}
	return timer, nil
}

func (t *timerService) StopTimer(token string, taskID *int) error {
	user, err := GetUserByToken(token, t.UserServiceApi, t.UserServicePort)
	if err != nil {
		if errors.Is(err, customServiceError.ErrTokenExpired) {
			tg.SendError(err.Error(), "api/timer/stop")
			return customServiceError.ErrTokenExpired
		}
		if errors.Is(err, customServiceError.ErrTokenNotFound) {
			tg.SendError(err.Error(), "api/timer/stop")
			return customServiceError.ErrTokenNotFound
		}
		tg.SendError(err.Error(), "api/timer/stop")
		log.Println("GetUserByToken, api/timer/stop, error:", err)
		return customServiceError.ErrUnknownError
	}

	return t.TimerRepository.StopTimerEntry(user.ID, taskID)
}
func (t *timerService) ResumeTimer(timerID int) (*model.TimerEntry, error) {
	err := t.TimerRepository.ResumeTimerEntry(timerID)
	if err != nil {
		return nil, err
	}
	// Получим обновлённый таймер
	timers, err := t.TimerRepository.GetTimersByUserID(timerID)
	if err != nil {
		return nil, err
	}
	// Найдём таймер с timerID
	for _, t := range timers {
		if t.ID == timerID {
			return &t, nil
		}
	}
	return nil, errors.New("timer not found after resume")
}

func (t *timerService) GetTimersByTask(taskID int) ([]model.TimerEntry, error) {
	return t.TimerRepository.GetTimersByTaskID(taskID)
}

func (t *timerService) GetTimersByUser(userID int) ([]model.TimerEntry, error) {
	return t.TimerRepository.GetTimersByUserID(userID)
}

func newTimerService(timerRepository repository.TimerRepository, userServiceApi, userServicePort string) *timerService {
	return &timerService{
		TimerRepository: timerRepository,
		UserServiceApi:  userServiceApi,
		UserServicePort: userServicePort,
	}
}
