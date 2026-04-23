package servicesConcrete

import (
	"context"
	"log"
	"os/exec"
	"server-watcher-app/src/models"
	modelsDTOs "server-watcher-app/src/models/dtos"
	enumModels "server-watcher-app/src/models/enum"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"strings"
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

func (s *pm2ProjectService) StartPm2Project(ctx context.Context, id int) (bool, error) {

	project, err := s.projectRepo.GetPm2ProjectById(ctx, id)

	if err != nil {
		return false, err
	}

	// if project.ProjectStartCommand == nil || project.ProjectPath == nil {
	// 	return false, fmt.Errorf("project start command or path is missing")
	// }

	// PM2'de mükerrer (duplicate) oluşmaması için önce varsa siliyoruz
	log.Printf("[PM2] Deleting existing process: %s", project.Name)
	deleteOut, _ := exec.CommandContext(ctx, "pm2", "delete", project.Name).CombinedOutput()
	log.Printf("[PM2] Delete output: %s", string(deleteOut))

	log.Printf("[PM2] Starting process: %s with command: %s", project.Name, *project.ProjectStartCommand)
	
	startCommand := *project.ProjectStartCommand
	// Eğer komut yanlışlıkla pm2 start ile başlıyorsa temizle
	startCommand = strings.TrimPrefix(startCommand, "pm2 start ")
	// Eğer içinde --name varsa o kısmı da temizlemeye çalışalım (basitçe)
	if idx := strings.Index(startCommand, " --name"); idx != -1 {
		startCommand = startCommand[:idx]
	}
	startCommand = strings.Trim(startCommand, "'\" ")

	cmd := exec.CommandContext(ctx, "pm2", "start", startCommand, "--name", project.Name)
	cmd.Dir = *project.ProjectPath

	startOut, err := cmd.CombinedOutput()
	log.Printf("[PM2] Start output: %s", string(startOut))

	if err != nil {
		log.Printf("[PM2] Start error: %v", err)
		return false, err
	}

	// UNIQUE hatasını önlemek için sadece belirli alanları güncelliyoruz
	err = s.projectRepo.Query(ctx).Model(&models.MonitoredEntity{}).Where("id = ?", project.ID).Updates(map[string]interface{}{
		"status":     string(enumModels.Running),
		"last_check": time.Now(),
	}).Error

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *pm2ProjectService) StopPm2Project(ctx context.Context, id int) (bool, error) {

	project, err := s.projectRepo.GetByID(ctx, id)

	if err != nil {
		return false, err
	}

	err = exec.CommandContext(ctx, "pm2", "stop", project.Name).Run()
	if err != nil {
		return false, err
	}

	// Sadece belirli alanları güncelliyoruz
	err = s.projectRepo.Query(ctx).Model(&models.MonitoredEntity{}).Where("id = ?", project.ID).Updates(map[string]interface{}{
		"status":     string(enumModels.Exited),
		"last_check": time.Now(),
	}).Error

	if err != nil {
		return false, err
	}

	return true, nil

}
