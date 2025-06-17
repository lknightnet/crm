package model

type TaskExecutor struct {
	ID         int `gorm:"primary_key"`
	ExecutorID int
	TaskID     int
}

type TaskExecutorWithoutName struct {
	ExecutorID int
}

type TaskExecutorWithName struct {
	ExecutorID   int    `json:"executor_id"`
	ExecutorName string `json:"executor_name"`
}
