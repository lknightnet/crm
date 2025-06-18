package timer

import (
	"project-service/internal/model"
	"time"
)

type TimerErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

type TimerResponse struct {
	ID             int           `json:"id"`
	TaskID         *int          `json:"task_id"`
	UserID         int           `json:"user_id"`
	StartedAt      *time.Time    `json:"started_at"`
	StoppedAt      *time.Time    `json:"stopped_at"`
	DurationSecond time.Duration `json:"duration_second"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	Active         bool          `json:"active"`
}

func NewTimerResponse(timer *model.TimerEntry) *TimerResponse {
	return &TimerResponse{
		ID:             timer.ID,
		TaskID:         timer.TaskID,
		UserID:         timer.UserID,
		StartedAt:      timer.StartedAt,
		StoppedAt:      timer.StoppedAt,
		DurationSecond: timer.DurationSecond,
		CreatedAt:      timer.CreatedAt,
		UpdatedAt:      timer.UpdatedAt,
		Active:         timer.Active,
	}
}
