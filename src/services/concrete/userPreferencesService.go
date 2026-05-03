package servicesConcrete

import (
	"context"
	"server-watcher-app/src/models"
	modelsDTOs "server-watcher-app/src/models/dtos"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	servicesAbstarct "server-watcher-app/src/services/abstract"
)

type userPreferencesService struct {
	repository repositoryConcrete.UserPreferencesRepository
}

func NewUserPreferencesService(repository repositoryConcrete.UserPreferencesRepository) servicesAbstarct.IUserPreferencesService {
	return &userPreferencesService{
		repository: repository,
	}
}

func (s *userPreferencesService) CreateServerStatPreference(ctx context.Context, userID uint, dto *modelsDTOs.ServerStatusSettingDTO) (bool, error) {

	existing, err := s.repository.GetByUserId(ctx, userID)

	if err != nil {
		return false, err
	}

	if existing != nil {
		existing.CPUThreshold = dto.CPUThreshold
		existing.DISKThreshold = dto.DISKThreshold
		existing.RAMThreshold = dto.RAMThreshold
		existing.NotificationCooldown = dto.NotificationCooldown

		err = s.repository.Update(ctx, existing)

		if err != nil {
			return false, err
		}

		return true, nil
	} else {

		newSSPreference := models.SSTSettings{
			UserID:               userID,
			CPUThreshold:         dto.CPUThreshold,
			DISKThreshold:        dto.DISKThreshold,
			RAMThreshold:         dto.RAMThreshold,
			NotificationCooldown: dto.NotificationCooldown,
		}

		err = s.repository.Create(ctx, &newSSPreference)

	}

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *userPreferencesService) ClearServerStatPreference(ctx context.Context, userID uint) (bool, error) {

	existing, err := s.repository.GetByUserId(ctx, userID)

	if err != nil {
		return false, nil
	}

	err = s.repository.Delete(ctx, int(existing.ID))

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *userPreferencesService) GetUserSSTSettings(ctx context.Context, userID uint) (*modelsDTOs.ServerStatusSettingDTO, error) {

	sstSettings, err := s.repository.GetByUserId(ctx, userID)

	if err != nil {
		return nil, err
	}

	if sstSettings == nil {

		return &modelsDTOs.ServerStatusSettingDTO{
			CPUThreshold:         80.0,
			DISKThreshold:        80.0,
			RAMThreshold:         80.0,
			NotificationCooldown: 5,
			IsDefault:            true,
		}, nil
	}

	newSSTSettings := modelsDTOs.ServerStatusSettingDTO{
		CPUThreshold:         sstSettings.CPUThreshold,
		DISKThreshold:        sstSettings.DISKThreshold,
		RAMThreshold:         sstSettings.RAMThreshold,
		NotificationCooldown: sstSettings.NotificationCooldown,
	}

	return &newSSTSettings, nil

}
