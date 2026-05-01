package modelsDTOs

import enumModels "server-watcher-app/src/models/enum"

type RegisterRequestDTO struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
}

type LoginRequestDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserResponseDTO struct {
	ID       int                 `json:"id"`
	Name     string              `json:"name"`
	Surname  string              `json:"surname"`
	UserRole enumModels.UserRole `json:"user_role"`
}

type RefreshTokenRequestDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UserFCMTokenRequestDTO struct {
	UserDeviceFCMToken string `json:"fcm_token" validate:"required"`
}
