package servicesAbstarct

import (
	"context"
	modelsDTOs "server-watcher-app/src/models/dtos"
)

type AuthService interface {
	Register(ctx context.Context, dto *modelsDTOs.RegisterRequestDTO) (bool, error)
	Login(ctx context.Context, dto *modelsDTOs.LoginRequestDTO) (*modelsDTOs.AuthResponseDTO, error)
	RefreshToken(ctx context.Context, dto *modelsDTOs.RefreshTokenRequestDTO) (*modelsDTOs.AuthResponseDTO, error)
	GetCurrentUserInformation(ctx context.Context, userId int) (*modelsDTOs.UserResponseDTO, error)
	CreateAdminUser(ctx context.Context, uName string, password string) error
	TakeUserFCMToken(ctx context.Context, userID int, dto *modelsDTOs.UserFCMTokenRequestDTO) (bool, error)
	RemoveUserFCMToken(ctx context.Context, userID int) (bool, error)
	GetAllUser(ctx context.Context, userID uint) (*[]modelsDTOs.UserResponseDTO, error)
}
