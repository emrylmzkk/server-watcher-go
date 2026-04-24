package servicesConcrete

import (
	"context"
	"log"
	modelsDTOs "server-watcher-app/src/models/dtos"
	repositoryAbstarct "server-watcher-app/src/repository/abstract"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"time"
)

type wathcerService struct {
	repo      repositoryAbstarct.MonitoredEntityRepository
	providers []servicesAbstarct.ProcessProvider
}

func NewWatcherService(repo repositoryAbstarct.MonitoredEntityRepository, providers []servicesAbstarct.ProcessProvider) servicesAbstarct.WatcherService {
	return &wathcerService{
		repo:      repo,
		providers: providers,
	}
}

func (s *wathcerService) SyncAll(ctx context.Context) error {
	for _, provider := range s.providers {
		processes, err := provider.ListProcesses(ctx)
		if err != nil {
			log.Printf("Provider %v hatası: %v", provider.GetProviderType(), err)

			continue
		}

		for _, proc := range processes {
			proc.LastCheck = time.Now()
			if err := s.repo.CreateOrUpdate(&proc); err != nil {
				log.Printf("DB Update hatası: %v", err)
			}
		}
	}
	return nil
}

func (s *wathcerService) ControlProcess(ctx context.Context, dto *modelsDTOs.ActionOnProject) (bool, error) {

	entity, err := s.repo.GetByExternalID(dto.ExternalID)
	if err != nil {
		return false, err
	}

	// İlgili provider'ı seç ve aksiyonu al
	for _, p := range s.providers {
		if p.GetProviderType() == entity.Type {
			if dto.Action == "start" {
				return true, p.StartProcess(ctx, dto.ExternalID)
			} else if dto.Action == "stop" {
				return true, p.StopProcess(ctx, dto.ExternalID)
			}
		}
	}
	return false, nil
}
