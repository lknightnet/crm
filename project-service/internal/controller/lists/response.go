package lists

type ListErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}
