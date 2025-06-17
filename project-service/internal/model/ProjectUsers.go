package model

import "time"

type ProjectUsers struct {
	ID        int
	UserID    int
	ProjectID int
	CreatedAt time.Time
}
