package repository

import (
	"errors"
	"gorm.io/gorm"
	"project-service/internal/model"
	"project-service/pkg/database"
	"time"
)

type timerRepository struct {
	db *database.PostgreSQL
}

func (t *timerRepository) StartTimerEntry(timerEntry *model.TimerEntry) error {
	var activeTimer model.TimerEntry

	query := t.db.DB.Where("user_id = ? AND active = ?", timerEntry.UserID, true)

	if timerEntry.TaskID == nil {
		// Если запускается таймер без задачи — ищем активный с задачей
		query = query.Where("task_id IS NOT NULL")
	} else {
		// Если запускается таймер с задачей — ищем активный без задачи
		query = query.Where("task_id IS NULL")
	}

	err := query.First(&activeTimer).Error
	if err == nil {
		now := time.Now()
		elapsed := now.Sub(*activeTimer.StartedAt)
		activeTimer.DurationSecond += elapsed
		activeTimer.StoppedAt = &now
		activeTimer.Active = false
		activeTimer.UpdatedAt = now

		if err := t.db.DB.Save(&activeTimer).Error; err != nil {
			return err
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Запускаем новый таймер
	now := time.Now()
	timerEntry.StartedAt = &now
	timerEntry.StoppedAt = nil
	timerEntry.DurationSecond = 0
	timerEntry.Active = true

	return t.db.DB.Create(timerEntry).Error
}

func (t *timerRepository) StopTimerEntry(userID int, taskID *int) error {
	var timer model.TimerEntry
	query := t.db.DB.Where("user_id = ? AND active = ?", userID, true)
	if taskID == nil {
		query = query.Where("task_id IS NULL")
	} else {
		query = query.Where("task_id = ?", *taskID)
	}

	if err := query.First(&timer).Error; err != nil {
		return err
	}

	now := time.Now()
	elapsed := now.Sub(*timer.StartedAt)
	timer.DurationSecond += elapsed
	timer.StoppedAt = &now
	timer.UpdatedAt = now
	timer.Active = false

	return t.db.DB.Save(&timer).Error
}

func (t *timerRepository) ResumeTimerEntry(timerID int) error {
	var timer model.TimerEntry
	if err := t.db.DB.First(&timer, timerID).Error; err != nil {
		return err
	}

	var activeTimer model.TimerEntry
	query := t.db.DB.Where("user_id = ? AND active = ?", timer.UserID, true)

	if timer.TaskID == nil {
		query = query.Where("task_id IS NOT NULL")
	} else {
		query = query.Where("task_id IS NULL")
	}

	err := query.First(&activeTimer).Error
	if err == nil {
		now := time.Now()
		elapsed := now.Sub(*activeTimer.StartedAt)
		activeTimer.DurationSecond += elapsed
		activeTimer.StoppedAt = &now
		activeTimer.Active = false
		activeTimer.UpdatedAt = now
		if err := t.db.DB.Save(&activeTimer).Error; err != nil {
			return err
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	now := time.Now()
	timer.StartedAt = &now
	timer.StoppedAt = nil
	timer.Active = true
	timer.UpdatedAt = now

	return t.db.DB.Save(&timer).Error
}

func (t *timerRepository) GetTimersByTaskID(taskID int) ([]model.TimerEntry, error) {
	var timers []model.TimerEntry
	if err := t.db.DB.Where("task_id = ?", taskID).Find(&timers).Error; err != nil {
		return nil, err
	}
	return timers, nil
}

func (t *timerRepository) GetTimersByUserID(userID int) ([]model.TimerEntry, error) {
	var timers []model.TimerEntry
	if err := t.db.DB.Where("user_id = ?", userID).Find(&timers).Error; err != nil {
		return nil, err
	}
	return timers, nil
}

func newTimerRepository(db *database.PostgreSQL) *timerRepository {
	return &timerRepository{db: db}
}
