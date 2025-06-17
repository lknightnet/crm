package project

type ProjectCreateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Users       []int   `json:"project_users"`
}
