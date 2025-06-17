package timer

import (
	"project-service/internal/model"
	"time"
)

type TimerGetResponse struct {
	ID             int
	StartAt        time.Time
	StopAt         *time.Time
	DurationSecond time.Duration
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Active         bool
}

func NewTimerGetResponse(timer *model.TimerEntry) *TimerGetResponse {
	return &TimerGetResponse{
		ID:             timer.ID,
		StartAt:        timer.StartAt,
		StopAt:         timer.StopAt,
		DurationSecond: timer.DurationSecond,
		CreatedAt:      timer.CreatedAt,
		UpdatedAt:      timer.UpdatedAt,
		Active:         timer.Active,
	}
}
