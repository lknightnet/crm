package model

import "time"

type TimerEntry struct {
	ID             int `gorm:"primary_key"`
	TaskID         *int
	UserID         int
	StartedAt      *time.Time
	StoppedAt      *time.Time
	DurationSecond time.Duration
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Active         bool
}
