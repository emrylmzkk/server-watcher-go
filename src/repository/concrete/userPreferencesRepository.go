package repositoryConcrete

import (
	"context"
	"errors"
	"server-watcher-app/src/models"
	repositoryAbstarct "server-watcher-app/src/repository/abstract"

	"gorm.io/gorm"
)

type UserPreferencesRepository interface {
	repositoryAbstarct.BaseRepository[models.SSTSettings]

	GetByUserId(ctx context.Context, userID uint) (*models.SSTSettings, error)
}

type userPreferencesRepository struct {
	repositoryAbstarct.BaseRepository[models.SSTSettings]
}

func NewUserPreferencesRepository(db *gorm.DB) UserPreferencesRepository {
	return &userPreferencesRepository{
		BaseRepository: repositoryAbstarct.NewBaseRepository[models.SSTSettings](db),
	}
}

func (r *userPreferencesRepository) GetByUserId(ctx context.Context, userID uint) (*models.SSTSettings, error) {

	var setting models.SSTSettings

	err := r.Query(ctx).
		Where("user_id = ?", userID).
		First(&setting).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &setting, nil

}
