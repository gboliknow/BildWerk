package transformer

import "github.com/gboliknow/bildwerk/internal/models"

func ToUserResponse(user *models.User) models.UserResponse {
	return models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}
