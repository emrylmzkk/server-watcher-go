package servicesConcrete

import (
	"context"
	"encoding/json"
	"server-watcher-app/src/models"
	modelsDTOs "server-watcher-app/src/models/dtos"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"gorm.io/gorm"
)

type dockerService struct {
	cli   *client.Client
	cache *models.DockerCache
	db    *gorm.DB
}

func NewDockerService(db *gorm.DB) (servicesAbstarct.DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &dockerService{
		cli:   cli,
		cache: &models.DockerCache{},
		db:    db,
	}, nil
}

func (s *dockerService) GetContainers(ctx context.Context) ([]modelsDTOs.DockerContainerResponseDTO, error) {
	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var response []modelsDTOs.DockerContainerResponseDTO
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		response = append(response, modelsDTOs.DockerContainerResponseDTO{
			Name:      name,
			IsRunning: c.State == "running",
			Uptime:    c.Status, // Status usually contains "Up X minutes" or "Exited (0) X minutes ago"
		})
	}

	return response, nil
}

func (s *dockerService) StartContainer(ctx context.Context, containerID string) error {
	return s.cli.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (s *dockerService) StopContainer(ctx context.Context, containerID string) error {
	return s.cli.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (s *dockerService) collectStatsParallel(ctx context.Context) ([]modelsDTOs.DockerStats, error) {

	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var responses []modelsDTOs.DockerStats
	var wg sync.WaitGroup
	var mu sync.Mutex

	sem := make(chan struct{}, 5)

	for _, item := range containers {

		wg.Add(1)
		sem <- struct{}{}

		go func(c container.Summary) {
			defer wg.Done()
			defer func() { <-sem }()

			name := ""
			if len(c.Names) > 0 {
				name = strings.TrimPrefix(c.Names[0], "/")
			}

			dto := modelsDTOs.DockerStats{
				Name:      name,
				IsRunning: c.State == "running",
				Uptime:    c.Status,
			}

			if c.State == "running" {

				stats, err := s.cli.ContainerStats(ctx, c.ID, false)
				if err == nil {
					defer stats.Body.Close()

					var raw modelsDTOs.DockerStatsRaw

					if err := json.NewDecoder(stats.Body).Decode(&raw); err == nil {

						cpuDelta := float64(raw.CPUStats.CPUUsage.TotalUsage - raw.PreCPUStats.CPUUsage.TotalUsage)
						systemDelta := float64(raw.CPUStats.SystemUsage - raw.PreCPUStats.SystemUsage)

						if systemDelta > 0 && cpuDelta > 0 {
							dto.CPU = (cpuDelta / systemDelta) * 100
						}

						dto.MemoryMB = float64(raw.MemoryStats.Usage) / 1024 / 1024

						for _, net := range raw.Networks {
							dto.NetworkRX += net.RxBytes
							dto.NetworkTX += net.TxBytes
						}
					}
				}
			}

			mu.Lock()
			responses = append(responses, dto)
			mu.Unlock()

		}(item)
	}

	wg.Wait()
	return responses, nil
}

func (s *dockerService) StartStatsCollector(ctx context.Context) {

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		dbTicker := time.NewTicker(3 * time.Minute) // Dakikada bir DB'ye kaydet

		for {
			select {
			case <-ticker.C:

				stats, err := s.collectStatsParallel(ctx)
				if err != nil {
					continue
				}

				s.cache.Mu.Lock()
				s.cache.Data = stats
				s.cache.Mu.Unlock()

			case <-dbTicker.C:
				// Cache'deki güncel veriyi DB'ye logla
				s.saveStatsToDB()

			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *dockerService) saveStatsToDB() {
	s.cache.Mu.RLock()
	data := s.cache.Data
	s.cache.Mu.RUnlock()

	if len(data) == 0 {
		return
	}

	var logs []models.ContainerStatLog
	for _, stat := range data {
		if stat.IsRunning {
			logs = append(logs, models.ContainerStatLog{
				ContainerName: stat.Name,
				CPU:           stat.CPU,
				MemoryMB:      stat.MemoryMB,
				NetworkRX:     stat.NetworkRX,
				NetworkTX:     stat.NetworkTX,
			})
		}
	}

	if len(logs) > 0 {
		// Batch insert for performance
		s.db.Create(&logs)
	}
}

func (s *dockerService) GetCachedStats() []modelsDTOs.DockerStats {
	s.cache.Mu.RLock()
	defer s.cache.Mu.RUnlock()
	return s.cache.Data
}
