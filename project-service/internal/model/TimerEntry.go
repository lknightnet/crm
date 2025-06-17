package model

import "time"

type TimerEntry struct {
	ID             int
	TaskID         int
	UserID         int
	StartAt        time.Time
	StopAt         *time.Time
	DurationSecond time.Duration
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Active         bool
}
