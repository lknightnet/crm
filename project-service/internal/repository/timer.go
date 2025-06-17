package repository

import (
	"errors"
	"gorm.io/gorm"
	"project-service/internal/model"
	"project-service/internal/repository/customRepositoryError"
	"project-service/pkg/database"
	"time"
)

type timerRepository struct {
	db *database.PostgreSQL
}

func (t *timerRepository) StartTimerEntry(timerEntry *model.TimerEntry) error {
	if err := t.db.DB.Create(timerEntry).Error; err != nil {
		return err
	}
	return nil
}

func (t *timerRepository) StopTimerEntry(userID int) error {
	var entry model.TimerEntry

	err := t.db.DB.Where("user_id = ? AND active =", userID, true).First(&entry).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return customRepositoryError.ErrTimerNotFound
		}
		return err
	}

	now := time.Now()
	duration := time.Duration(now.Sub(entry.StartAt).Seconds())

	entry.StopAt = &now
	entry.DurationSecond = duration
	entry.Active = false
	entry.UpdatedAt = time.Now()

	if err := t.db.DB.Save(&entry).Error; err != nil {
		return err
	}

	return nil
}

func (t *timerRepository) GetTimersByTaskID(taskID int) ([]model.TimerEntry, error) {
	var timers []model.TimerEntry
	err := t.db.DB.
		Where("task_id = ?", taskID).
		Find(&timers).Error

	if err != nil {
		return nil, err
	}

	if len(timers) == 0 {
		return nil, customRepositoryError.ErrTimerNotFound
	}

	return timers, nil
}

func (t *timerRepository) GetTimersByUserID(userID int) ([]model.TimerEntry, error) {
	var timers []model.TimerEntry
	err := t.db.DB.
		Where("user_id = ?", userID).
		Find(&timers).Error

	if err != nil {
		return nil, err
	}

	if len(timers) == 0 {
		return nil, customRepositoryError.ErrTimerNotFound
	}

	return timers, nil
}

func newTimerRepository(db *database.PostgreSQL) *timerRepository {
	return &timerRepository{db: db}
}
