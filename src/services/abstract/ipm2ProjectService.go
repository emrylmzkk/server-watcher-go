package servicesAbstarct

import (
	"context"
	modelsDTOs "server-watcher-app/src/models/dtos"
)

type Pm2ProjectService interface {
	AddNewProject(ctx context.Context, dto *modelsDTOs.CreatePm2ProjectRequestDTO) (bool, error)
	DeletePm2Project(ctx context.Context, id int) (bool, error)
	GetPm2Projects(ctx context.Context) (*[]modelsDTOs.Pm2ProjectResponseDTO, error)
	StartPm2Project(ctx context.Context, id int) (bool, error)
	StopPm2Project(ctx context.Context, id int) (bool, error)
	UpdatePm2Project(ctx context.Context, id int, dto *modelsDTOs.UpdatePm2ProjectRequestDTO) (bool, error)
	ResetPm2Process(ctx context.Context) (bool, error)
}
