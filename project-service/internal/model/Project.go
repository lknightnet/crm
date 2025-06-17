package model

type Project struct {
	ID          int `gorm:"primary_key"`
	Name        string
	Description *string
	CreatedID   int //кто создал
	Visibility  bool
}

type ProjectWithCreatedUsername struct {
	ID          int
	Name        string
	Description *string
	CreatedName string
	CreatedID   int //кто создал
	Visibility  bool
}

type ProjectWithProjectUsers struct {
	ID           int
	Name         string
	Description  *string
	CreatedID    int //кто создал
	Visibility   bool
	ProjectUsers []ProjectUsers
}
