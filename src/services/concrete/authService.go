package servicesConcrete

import (
	"context"
	"errors"
	"log"
	"server-watcher-app/src/generic"
	"server-watcher-app/src/models"
	modelsDTOs "server-watcher-app/src/models/dtos"
	enumModels "server-watcher-app/src/models/enum"
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

func (s *authService) Register(ctx context.Context, userID uint, dto *modelsDTOs.RegisterRequestDTO) (bool, error) {

	isAdmin, err := s.userRepository.IsUserAdmin(ctx, userID)

	if err != nil || !isAdmin {
		return false, errors.New("Only the admin can register user")
	}

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

func (s *authService) TakeUserFCMToken(ctx context.Context, userID int, dto *modelsDTOs.UserFCMTokenRequestDTO) (bool, error) {

	//var user models.User

	user, err := s.userRepository.GetByID(ctx, userID)

	if err != nil {
		return false, errors.New("user not found by information")
	}

	user.FCMToken = &dto.UserDeviceFCMToken

	err = s.userRepository.Update(ctx, user)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *authService) RemoveUserFCMToken(ctx context.Context, userID int) (bool, error) {

	//var user models.User

	user, err := s.userRepository.GetByID(ctx, userID)

	if err != nil {
		return false, errors.New("user not found by information")
	}

	user.FCMToken = nil

	err = s.userRepository.Update(ctx, user)

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

func (s *authService) CreateAdminUser(ctx context.Context, uName string, password string) error {

	isExist, err := s.userRepository.IsUserExists(ctx, uName)

	if err != nil {
		return nil
	}

	if isExist == false {

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			return err
		}

		user := models.User{
			Username: uName,
			Password: string(hashed),
			Name:     "Admin",
			Surname:  "Admin",
			UserRole: 95,
		}

		err = s.userRepository.Create(ctx, &user)

		if err != nil {
			return err
		}
		log.Println("[AdminAuthProc] Admin user created successfuly")
		return nil
	}

	return nil

}

func (s *authService) GetCurrentUserInformation(ctx context.Context, userId int) (*modelsDTOs.UserResponseDTO, error) {

	user, err := s.userRepository.GetByID(ctx, userId)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &modelsDTOs.UserResponseDTO{
		ID:       int(user.ID),
		Username: user.Username,
		Name:     user.Name,
		Surname:  user.Surname,
		UserRole: user.UserRole,
	}, nil

}

func (s *authService) GetAllUser(ctx context.Context, userID uint) (*[]modelsDTOs.UserResponseDTO, error) {

	var users []modelsDTOs.UserResponseDTO

	isAdmin, err := s.userRepository.IsUserAdmin(ctx, userID)

	if err != nil || !isAdmin {
		return &users, nil
	}

	dbUsers, err := s.userRepository.GetAll(ctx)

	if err != nil {
		return &users, err
	}

	for _, user := range dbUsers {
		users = append(users, modelsDTOs.UserResponseDTO{
			ID:       int(user.ID),
			Username: user.Username,
			Name:     user.Name,
			Surname:  user.Surname,
			UserRole: user.UserRole,
		})
	}

	return &users, nil

}

func (s *authService) DeleteUser(ctx context.Context, id int, userID uint) (bool, error) {

	user, err := s.userRepository.GetByID(ctx, int(userID))

	if err != nil {
		return false, nil
	}

	if user.UserRole != enumModels.Admin {
		return false, errors.New("Only the admin can delete users")
	}

	err = s.userRepository.Delete(ctx, id)

	if err != nil {
		return false, err
	}

	return true, nil

}
