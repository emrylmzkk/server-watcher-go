package modelsDTOs

import enumModels "server-watcher-app/src/models/enum"

type CreatePm2ProjectRequestDTO struct {
	ExternalID          string                        `json:"external_id" binding:"required"`
	Name                string                        `json:"name" binding:"required"`
	ProjectPath         string                        `json:"project_path" binding:"required"`
	ProjectStartCommand string                        `json:"project_start_command" binding:"required"`
	ProjectRuntimeType  enumModels.ProjectRuntimeType `json:"project_runtime_tpye" binding:"required"`
}

type Pm2ProjectResponseDTO struct {
	ID                  int                           `json:"id"`
	ExternalID          string                        `json:"external_id" binding:"required"`
	Name                string                        `json:"name" binding:"required"`
	Type                enumModels.ProcessType        `json:"project_type"`
	Status              enumModels.ProjectStatus      `json:"status"`
	ProjectPath         string                        `json:"project_path" binding:"required"`
	ProjectStartCommand string                        `json:"project_start_command" binding:"required"`
	ProjectRuntimeType  enumModels.ProjectRuntimeType `json:"project_runtime_tpye" binding:"required"`
}

// type Pm2ProjectActionRequestDTO struct {
// 	ExternalID string `json:"external_id" binding:"required"`
// 	Action     string `json:"action" binding:"required"`
// }

type UpdatePm2ProjectRequestDTO struct {
	Name                string                         `json:"name"`
	ProjectPath         *string                        `json:"project_path"`
	ProjectStartCommand *string                        `json:"project_start_command"`
	ProjectRuntimeType  *enumModels.ProjectRuntimeType `json:"project_runtime_tpye"`
}
