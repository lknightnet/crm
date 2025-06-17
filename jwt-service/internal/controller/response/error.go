package response

type ErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}
