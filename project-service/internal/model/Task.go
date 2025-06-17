package model

import "time"

type Task struct {
	ID          int `gorm:"primary_key"`
	Name        string
	Description *string
	Priority    int
	Deadline    *time.Time
	CreatedID   int //кто создал
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ProjectID   int
	Status      string
}

type TaskWithExecutor struct {
	ID          int
	Name        string
	Description *string
	Priority    int
	Deadline    *time.Time
	CreatedID   int //кто создал
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ProjectID   int
	ProjectName string
	Status      string
	Executors   []TaskExecutorWithoutName
}

type TaskWithExecutorWithName struct {
	ID          int
	Name        string
	Description *string
	Priority    int
	Deadline    *time.Time
	CreatedID   int //кто создал
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ProjectID   int
	ProjectName string
	Status      string
	Executors   []TaskExecutorWithName
}

const (
	PRIORITY_LOW = iota
	PRIORITY_MEDIUM
	PRIORITY_HIGH
	PRIORITY_CRITICAL
)

const (
	STATUS_NOT_STARTED = "Не начата"
	STATUS_STOPPED     = "Остановлена"
	STATUS_IN_PROGRESS = "Выполняется"
	STATUS_COMPLETED   = "Заверешена"
)
