package repositoryConcrete

import (
	"context"
	"server-watcher-app/src/models"
	repositoryAbstarct "server-watcher-app/src/repository/abstract"

	"gorm.io/gorm"
)

type UserRepository interface {
	repositoryAbstarct.BaseRepository[models.User]
	IsUserExists(ctx context.Context, username string) (bool, error)
	GetByUserName(ctx context.Context, username string) (*models.User, error)
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
