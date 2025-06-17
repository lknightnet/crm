package lists

type ListRequest struct {
	ProjectID int     `json:"project_id"`
	ListID    int     `json:"list_id"`
	Key       *string `json:"key"`
	Value     *string `json:"value"`
}
