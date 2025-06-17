package response

import "user-service/internal/model"

type GetUserByUserNameResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewUsersByUserNameResponse(users []model.User) []GetUserByUserNameResponse {
	usersWithUsername := make([]GetUserByUserNameResponse, len(users))
	for i, user := range users {
		usersWithUsername[i] = GetUserByUserNameResponse{
			ID:   user.ID,
			Name: user.Name,
		}
	}
	return usersWithUsername
}
