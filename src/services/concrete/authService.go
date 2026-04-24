package servicesConcrete

import (
	"context"
	"errors"
	"server-watcher-app/src/generic"
	"server-watcher-app/src/models"
	modelsDTOs "server-watcher-app/src/models/dtos"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository repositoryConcrete.UserRepository
}

func NewAuthService(userRepository repositoryConcrete.UserRepository) servicesAbstarct.AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func (s *authService) Register(ctx context.Context, dto *modelsDTOs.RegisterRequestDTO) (bool, error) {

	isExist, err := s.userRepository.IsUserExists(ctx, dto.Username)

	if isExist == true || err != nil {
		return false, errors.New("username already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)

	if err != nil {
		return false, err
	}

	user := models.User{
		Username: dto.Username,
		Password: string(hashed),
		Name:     dto.Name,
		Surname:  dto.Surname,
		UserRole: 2,
	}

	err = s.userRepository.Create(ctx, &user)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *authService) Login(ctx context.Context, dto *modelsDTOs.LoginRequestDTO) (*modelsDTOs.AuthResponseDTO, error) {

	user, err := s.userRepository.GetByUserName(ctx, dto.Username)

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
	// 	return nil, errors.New("invalid username or password")
	// }

	access, refresh, err := generic.GenerateTokenPair(user)

	return &modelsDTOs.AuthResponseDTO{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil

}

func (s *authService) RefreshToken(ctx context.Context, dto *modelsDTOs.RefreshTokenRequestDTO) (*modelsDTOs.AuthResponseDTO, error) {

	claims, err := generic.ValidateToken(dto.RefreshToken)

	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return nil, errors.New("token type must be refresh")
	}

	user, err := s.userRepository.GetByID(ctx, int(claims.UserID))

	if err != nil {
		return nil, errors.New("user not found")
	}

	access, refresh, err := generic.GenerateTokenPair(user)

	return &modelsDTOs.AuthResponseDTO{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil

}
