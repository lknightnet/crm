package request

type UpdateUser struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
