package timer

type TimerErrorResponse struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}
