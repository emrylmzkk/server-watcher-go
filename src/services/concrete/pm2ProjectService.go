package servicesConcrete

import (
	"context"
	"server-watcher-app/src/models"
	modelsDTOs "server-watcher-app/src/models/dtos"
	enumModels "server-watcher-app/src/models/enum"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"time"
)

type pm2ProjectService struct {
	projectRepo repositoryConcrete.ProjectRepository
}

func NewPm2ProjectService(projectRepo repositoryConcrete.ProjectRepository) servicesAbstarct.Pm2ProjectService {
	return &pm2ProjectService{
		projectRepo: projectRepo,
	}
}

func (s *pm2ProjectService) AddNewProject(ctx context.Context, dto *modelsDTOs.CreatePm2ProjectRequestDTO) (bool, error) {

	var existing models.MonitoredEntity
	err := s.projectRepo.Query(ctx).Where("external_id = ?", dto.ExternalID).First(&existing).Error

	if err == nil {
		existing.ProjectPath = &dto.ProjectPath
		existing.ProjectStartCommand = &dto.ProjectStartCommand
		existing.ProjectRuntimeType = &dto.ProjectRuntimeType
		existing.Name = dto.Name
		existing.Type = enumModels.PM2

		err = s.projectRepo.Update(ctx, &existing)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	newPm2Project := models.MonitoredEntity{
		ExternalID:          dto.ExternalID,
		Name:                dto.Name,
		Type:                enumModels.PM2,
		Status:              string(enumModels.Exited),
		LastCheck:           time.Now(),
		ProjectPath:         &dto.ProjectPath,
		ProjectStartCommand: &dto.ProjectStartCommand,
		ProjectRuntimeType:  &dto.ProjectRuntimeType,
	}

	err = s.projectRepo.Create(ctx, &newPm2Project)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *pm2ProjectService) DeletePm2Project(ctx context.Context, id int) (bool, error) {

	_, err := s.projectRepo.GetByID(ctx, id)

	if err != nil {
		return false, err
	}

	err = s.projectRepo.Delete(ctx, id)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *pm2ProjectService) GetPm2Projects(ctx context.Context) (*[]modelsDTOs.Pm2ProjectResponseDTO, error) {

	var responses []modelsDTOs.Pm2ProjectResponseDTO

	pm2projects, err := s.projectRepo.GetPm2Projects(ctx)

	if err != nil {
		return &responses, err
	}

	for _, response := range pm2projects {

		var projectPath string
		if response.ProjectPath != nil {
			projectPath = *response.ProjectPath
		}

		var projectStartCommand string
		if response.ProjectStartCommand != nil {
			projectStartCommand = *response.ProjectStartCommand
		}

		var projectRuntimeType enumModels.ProjectRuntimeType
		if response.ProjectRuntimeType != nil {
			projectRuntimeType = *response.ProjectRuntimeType
		}

		responses = append(responses, modelsDTOs.Pm2ProjectResponseDTO{
			ID:                  response.ID,
			ExternalID:          response.ExternalID,
			Name:                response.Name,
			Type:                response.Type,
			Status:              enumModels.ProjectStatus(response.Status),
			ProjectPath:         projectPath,
			ProjectStartCommand: projectStartCommand,
			ProjectRuntimeType:  projectRuntimeType,
		})

	}

	return &responses, nil

}
