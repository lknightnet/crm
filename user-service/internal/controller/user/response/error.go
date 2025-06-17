package response

type ErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

type MessageResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
