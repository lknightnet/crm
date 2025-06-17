package model

import "time"

type ProjectUsers struct {
	ID        int `gorm:"primary_key"`
	UserID    int
	ProjectID int
	CreatedAt time.Time
}

type ProjectUsersDTO struct {
	UserID   int
	UserName string
}
