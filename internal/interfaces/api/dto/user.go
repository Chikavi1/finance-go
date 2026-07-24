package dto

type UpdateProfileRequest struct {
	Name      string  `json:"name" validate:"required,min=2,max=255"`
	AvatarURL *string `json:"avatar_url"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}
