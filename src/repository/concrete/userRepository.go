package repositoryConcrete

import (
	"context"
	"errors"
	"server-watcher-app/src/models"
	enumModels "server-watcher-app/src/models/enum"
	repositoryAbstarct "server-watcher-app/src/repository/abstract"

	"gorm.io/gorm"
)

type UserRepository interface {
	repositoryAbstarct.BaseRepository[models.User]
	IsUserExists(ctx context.Context, username string) (bool, error)
	GetByUserName(ctx context.Context, username string) (*models.User, error)
	IsUserAdmin(ctx context.Context, userId uint) (bool, error)
	GetAdminUser(ctx context.Context) (*models.User, error)
	GetUserServerStatNotifSetting(ctx context.Context, userID uint) (*models.SSTSettings, error)
}

type userRepository struct {
	repositoryAbstarct.BaseRepository[models.User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: repositoryAbstarct.NewBaseRepository[models.User](db),
	}
}

func (r *userRepository) IsUserExists(ctx context.Context, username string) (bool, error) {

	var user models.User

	err := r.Query(ctx).Where("username = ?", username).First(&user).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil

}

func (r *userRepository) GetByUserName(ctx context.Context, username string) (*models.User, error) {

	var user models.User

	err := r.Query(ctx).
		Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *userRepository) IsUserAdmin(ctx context.Context, userId uint) (bool, error) {

	var user models.User

	err := r.Query(ctx).
		Select("user_role").
		Where("id = ?", userId).
		First(&user).Error

	if err != nil {
		return false, err
	}

	return user.UserRole == enumModels.Admin, nil

}

func (r *userRepository) GetAdminUser(ctx context.Context) (*models.User, error) {

	var user models.User

	err := r.Query(ctx).
		Where("user_role = ?", enumModels.Admin).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *userRepository) GetUserServerStatNotifSetting(ctx context.Context, userID uint) (*models.SSTSettings, error) {

	var settings models.SSTSettings

	err := r.Query(ctx).
		Model(&models.SSTSettings{}).
		Where("user_id = ?", userID).
		First(&settings).
		Error

	if err != nil {
		return &models.SSTSettings{
			CPUThreshold:         80.0,
			RAMThreshold:         80.0,
			DISKThreshold:        80.0,
			NotificationCooldown: 5,
		}, nil

	}

	return &settings, nil

}
