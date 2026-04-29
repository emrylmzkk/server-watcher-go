package servicesConcrete

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"server-watcher-app/src/generic"
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

func (s *pm2ProjectService) GetPm2ProjectById(ctx context.Context, id int) (*modelsDTOs.Pm2ProjectResponseDTO, error) {

	project, err := s.projectRepo.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	response := &modelsDTOs.Pm2ProjectResponseDTO{
		ID:         project.ID,
		ExternalID: project.ExternalID,
		Name:       project.Name,
		Type:       project.Type,
		Status:     enumModels.ProjectStatus(project.Status),
		// ProjectPath:         *project.ProjectPath,
		// ProjectStartCommand: *project.ProjectStartCommand,
		// ProjectRuntimeType:  *project.ProjectRuntimeType,
	}

	if project.ProjectPath != nil {
		response.ProjectPath = *project.ProjectPath
	}

	if project.ProjectStartCommand != nil {
		response.ProjectStartCommand = *project.ProjectStartCommand
	}

	if project.ProjectRuntimeType != nil {
		response.ProjectRuntimeType = *project.ProjectRuntimeType
	}

	return response, nil

}

func (s *pm2ProjectService) UpdatePm2Project(ctx context.Context, id int, dto *modelsDTOs.UpdatePm2ProjectRequestDTO) (bool, error) {

	pm2Project, err := s.projectRepo.GetByID(ctx, id)

	if err != nil {
		return false, err
	}

	if dto.Name != "" {
		pm2Project.Name = dto.Name
	}
	if dto.ProjectPath != nil {
		pm2Project.ProjectPath = dto.ProjectPath
	}
	if dto.ProjectStartCommand != nil {
		pm2Project.ProjectStartCommand = dto.ProjectStartCommand
	}
	if dto.ProjectRuntimeType != nil {
		pm2Project.ProjectRuntimeType = dto.ProjectRuntimeType
	}

	err = s.projectRepo.Update(ctx, pm2Project)
	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *pm2ProjectService) ClearAndDeletePm2Project(ctx context.Context, id int) (bool, error) {

	project, err := s.projectRepo.GetByID(ctx, id)

	if err != nil {
		return false, err
	}

	cmd := generic.NewCmd(ctx, "pm2", "delete", project.Name)

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Printf("[PM2 Warning] Süreç PM2 listesinde bulunamadı veya silinemedi: %s", string(output))
		return false, err
	}

	err = s.projectRepo.Delete(ctx, id)

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

	cmd := generic.NewCmd(ctx, "pm2", "start", project.ExternalID)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.New(string(output))
	}

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

	project, err := s.projectRepo.GetPm2ProjectById(ctx, id)

	if err != nil {
		return false, err
	}

	// if project.Name == "" {
	// 	return false, errors.New("project name not found")
	// }

	// log.Println("project name: ", project.Name)

	cmd := generic.NewCmd(ctx, "pm2", "stop", project.Name)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.New(string(output))
	}

	project.Status = string(enumModels.Exited)
	project.LastCheck = time.Now()

	err = s.projectRepo.Update(ctx, &project)

	if err != nil {
		return false, err
	}

	return true, nil

}

func (s *pm2ProjectService) ResetPm2Process(ctx context.Context) (bool, error) {

	cmd := generic.NewCmd(ctx, "pm2", "delete", "all")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.New(string(output))
	}

	projects, err := s.projectRepo.GetPm2Projects(ctx)
	if err != nil {
		return false, err
	}

	for _, project := range projects {
		err = s.projectRepo.Query(ctx).Model(&models.MonitoredEntity{}).Where("id = ?", project.ID).Updates(map[string]interface{}{
			"status":     string(enumModels.Exited),
			"last_check": time.Now(),
		}).Error
		if err != nil {
			return false, err
		}
	}

	return true, nil

}

// func (s *pm2ProjectService) GetPm2InsideList(ctx context.Context) (*[]modelsDTOs.Pm2InsideListResponseDTO, error) {

// 	var responses []modelsDTOs.Pm2InsideListResponseDTO

// 	cmd := generic.NewCmd(ctx, "pm2 jlist | jq '.[] | {name: .name, status: .pm2_env.status, pm_id: .pm_id}'")

// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return &responses, errors.New(string(output))
// 	}

// 	return &responses, nil

// }
func (s *pm2ProjectService) GetPm2InsideList(ctx context.Context) (*[]modelsDTOs.Pm2InsideListResponseDTO, error) {

	var responses []modelsDTOs.Pm2InsideListResponseDTO

	cmd := generic.NewCmd(ctx, "pm2", "jlist")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &responses, errors.New(string(output))
	}

	//float64(107978752) / 1024 / 1024 --> memory formatlanmasi gerekiyor

	// json --> struct

	//👉 JSON → struct dönüşümü = Unmarshal
	//👉 struct → JSON dönüşümü = Marshal

	if err := json.Unmarshal(output, &responses); err != nil {
		return nil, err
	}

	return &responses, nil

}

func (s *pm2ProjectService) CreateExamplePm2Project(ctx context.Context) error {

	mockExternal := "example-project"

	isExist, err := s.projectRepo.IsProjectExist(ctx, mockExternal)

	if err != nil {
		return nil
	}

	if isExist == false {

		projectPath := "/root/test"
		projectStartCommand := "go run main.go"
		projectRuntimeType := enumModels.Go

		pm2Project := models.MonitoredEntity{
			ExternalID:          mockExternal,
			Name:                mockExternal,
			Type:                enumModels.PM2,
			Status:              string(enumModels.Exited),
			LastCheck:           time.Now(),
			ProjectPath:         &projectPath,
			ProjectStartCommand: &projectStartCommand,
			ProjectRuntimeType:  &projectRuntimeType,
		}

		err = s.projectRepo.Create(ctx, &pm2Project)

		if err != nil {
			return err
		}

		log.Println("[MockPm2Project] Mock Pm2 project created successfuly")
		return nil

	}

	return nil

}

func (s *pm2ProjectService) SyncPm2Projects(ctx context.Context) (bool, error) {

	currentPm2InsideList, err := s.GetPm2InsideList(ctx)

	if err != nil {
		return false, nil
	}

	dbPm2List, err := s.projectRepo.GetPm2Projects(ctx)

	if err != nil {
		return false, nil
	}

	if len(*currentPm2InsideList) == 0 {
		return false, errors.New("No processes were found in the PM2 system; the deletion process was halted for security reasons")
	}

	dbMap := make(map[string]models.MonitoredEntity)

	for _, pm2Project := range dbPm2List {
		dbMap[pm2Project.ExternalID] = pm2Project
	}

	for _, currentPm2 := range *currentPm2InsideList {

		status := string(enumModels.Exited)

		if currentPm2.Pm2Env.Status == "online" {
			status = string(enumModels.Running)
		}

		if existing, ok := dbMap[currentPm2.Name]; ok {

			err = s.projectRepo.Query(ctx).Model(&models.MonitoredEntity{}).
				Where("id = ?", existing.ID).
				Updates(map[string]interface{}{
					"stauts":     status,
					"last_check": time.Now(),
				}).Error

		} else {

			newProject := models.MonitoredEntity{
				ExternalID:          currentPm2.Name,
				Name:                currentPm2.Name,
				Type:                enumModels.PM2,
				Status:              status,
				LastCheck:           time.Now(),
				ProjectPath:         &currentPm2.Pm2Env.ProjectDir,
				ProjectStartCommand: &currentPm2.Pm2Env.ExecPath,
			}

			err = s.projectRepo.Create(ctx, &newProject)

		}

		if err != nil {
			log.Printf("Sync hatası (%s): %v", currentPm2.Name, err)
		}

		delete(dbMap, currentPm2.Name)

	}

	// for _, exitedProject := range dbMap {

	// 	err := s.projectRepo.Delete(ctx, exitedProject.ID)
	// 	if err != nil {
	// 		log.Printf("Silme hatası (ID: %d, Name: %s): %v", exitedProject.ID, exitedProject.Name, err)
	// 	} else {
	// 		log.Printf("[Sync] Sistemde bulunmayan proje DB'den silindi: %s", exitedProject.Name)
	// 	}

	// }

	for _, leftProject := range dbMap {
		s.projectRepo.Query(ctx).Model(&models.MonitoredEntity{}).
			Where("id = ?", leftProject.ID).
			Updates(map[string]interface{}{
				"status": string(enumModels.Deleted),
			})

	}

	return true, nil

}
