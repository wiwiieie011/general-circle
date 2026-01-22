package dto

import "user-service/internal/models"

func ToUserResponse(u *models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      string(u.Role),
	}
}

func ToMeResponse(u *models.User) MeResponse {
	return MeResponse{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      string(u.Role),
	}
}

func ToPublicUserResponse(u *models.User) PublicUserResponse {
	return PublicUserResponse{
		ID:          u.ID,
		FullName:    u.FirstName + " " + u.LastName,
		IsOrganizer: u.Role == models.RoleOrganizer,
	}
}
