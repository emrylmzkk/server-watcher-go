package servicesAbstarct

import (
	"context"
	modelsDTOs "server-watcher-app/src/models/dtos"
)

type AuthService interface {
	Register(ctx context.Context, dto *modelsDTOs.RegisterRequestDTO) (bool, error)
	Login(ctx context.Context, dto *modelsDTOs.LoginRequestDTO) (*modelsDTOs.AuthResponseDTO, error)
	RefreshToken(ctx context.Context, refreshToken string) (*modelsDTOs.AuthResponseDTO, error)
}
