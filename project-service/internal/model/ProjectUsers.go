package model

import "time"

type ProjectUsers struct {
	ID        int `gorm:"primary_key"`
	UserID    int
	ProjectID int
	CreatedAt time.Time
}

type ProjectUsersDTO struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
}
