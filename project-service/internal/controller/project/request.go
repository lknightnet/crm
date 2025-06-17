package project

type ProjectCreateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Users       []int   `json:"project_users"`
}

type ProjectUpdateRequest struct {
	ProjectID   int     `json:"project_id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Users       *[]int  `json:"project_users"`
}
