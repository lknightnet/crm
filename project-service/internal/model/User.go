package model

import "time"

type User struct {
	ID            int
	Email         string
	Password      string
	Name          string
	DayOfBirthday *time.Time
	UserImage     *string
	CreatedAt     time.Time
	Alias         *string
	UUID          string
	Active        bool
	ClusterID     int
}
