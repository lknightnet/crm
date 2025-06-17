package task

import "time"

type TaskCreateRequest struct {
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	LevelPriority int        `json:"level_priority"`
	Deadline      *time.Time `json:"deadline"`
	ProjectID     int        `json:"project_id"`
	Executors     []int      `json:"executors"`
}

type TaskUpdateRequest struct {
	TaskID        int        `json:"task_id"`
	Name          *string    `json:"name"`
	Description   *string    `json:"description"`
	LevelPriority *int       `json:"level_priority"`
	Deadline      *time.Time `json:"deadline"`
	ProjectID     *int       `json:"project_id"`
	Executors     *[]int     `json:"executors"`
}
