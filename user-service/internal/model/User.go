package model

import "time"

type User struct {
	ID            int `gorm:"primary_key"`
	Email         string
	Password      string
	Name          string
	DayOfBirthday *time.Time
	UserImage     *string
	CreatedAt     time.Time
	UUID          string
	Active        bool
	LastLogin     time.Time
}
