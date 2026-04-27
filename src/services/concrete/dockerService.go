package servicesConcrete

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	genericInfluxDB "server-watcher-app/src/generic/influxDB"
	"server-watcher-app/src/models"
	modelsDTOs "server-watcher-app/src/models/dtos"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gorm.io/gorm"
)

type dockerService struct {
	cli          *client.Client
	cache        *models.DockerCache
	db           *gorm.DB
	influxClient *genericInfluxDB.InfluxClient
}

func NewDockerService(db *gorm.DB, influxClient *genericInfluxDB.InfluxClient) (servicesAbstarct.DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &dockerService{
		cli:          cli,
		cache:        &models.DockerCache{},
		db:           db,
		influxClient: influxClient,
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

func (s *dockerService) GetActiveContainers(ctx context.Context) ([]modelsDTOs.DockerContainerResponseDTO, error) {

	filter := filters.NewArgs()
	filter.Add("status", "running")

	containers, err := s.cli.ContainerList(ctx, container.ListOptions{
		All:     false,
		Filters: filter,
	})

	if err != nil {
		return nil, err
	}

	var response []modelsDTOs.DockerContainerResponseDTO

	for _, c := range containers {
		name := ""

		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		inspect, _ := s.cli.ContainerInspect(ctx, c.ID)

		var hostPort string

		for _, bindings := range inspect.NetworkSettings.Ports {
			for _, b := range bindings {
				if b.HostPort != "" {
					hostPort = b.HostPort
				}
			}
		}

		response = append(response, modelsDTOs.DockerContainerResponseDTO{
			Name:      name,
			IsRunning: true,
			Uptime:    c.Status,
			HostPort:  hostPort,
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

			inspect, err := s.cli.ContainerInspect(ctx, c.ID)
			if err != nil {
				return
			}

			var hostPort string

			for _, bindings := range inspect.NetworkSettings.Ports {
				for _, b := range bindings {
					if b.HostPort != "" {
						hostPort = b.HostPort
					}
				}
			}

			dto := modelsDTOs.DockerStats{
				Name:      name,
				IsRunning: c.State == "running",
				Uptime:    c.Status,
				HostPort:  hostPort,
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
		dbTicker := time.NewTicker(3 * time.Minute)
		logTicker := time.NewTicker(5 * time.Second) // her 5 saniyede log topla

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
				s.saveStatsToDB()

			case <-logTicker.C:
				s.collectAndSaveLogs(ctx)

			case <-ctx.Done():
				return
			}
		}
	}()
}

// func (s *dockerService) saveStatsToDB() {
// 	s.cache.Mu.RLock()
// 	data := s.cache.Data
// 	s.cache.Mu.RUnlock()

// 	if len(data) == 0 {
// 		return
// 	}

// 	var logs []models.ContainerStatLog
// 	for _, stat := range data {
// 		if stat.IsRunning {
// 			logs = append(logs, models.ContainerStatLog{
// 				ContainerName: stat.Name,
// 				CPU:           stat.CPU,
// 				MemoryMB:      stat.MemoryMB,
// 				NetworkRX:     stat.NetworkRX,
// 				NetworkTX:     stat.NetworkTX,
// 			})
// 		}
// 	}

// 	if len(logs) > 0 {
// 		// Batch insert for performance
// 		s.db.Create(&logs)
// 	}
// }

func (s *dockerService) saveStatsToDB() {

	s.cache.Mu.RLock()
	data := s.cache.Data
	s.cache.Mu.RUnlock()

	if len(data) == 0 {
		return
	}

	ctx := context.Background()
	now := time.Now()

	for _, stat := range data {

		if !stat.IsRunning {
			continue
		}

		point := influxdb2.NewPointWithMeasurement("container_stats").
			AddTag("container_name", stat.Name).
			AddField("cpu", stat.CPU).
			AddField("memory_mb", stat.MemoryMB).
			AddField("network_rx", stat.NetworkRX).
			AddField("network_tx", stat.NetworkTX).
			SetTime(now)

		err := s.influxClient.WriteAPI.WritePoint(ctx, point)

		if err != nil {
			log.Printf("InfluxDB yazma hatası [%s]: %v", stat.Name, err)
		}

	}

}

func (s *dockerService) GetCachedStats() []modelsDTOs.DockerStats {
	s.cache.Mu.RLock()
	defer s.cache.Mu.RUnlock()
	return s.cache.Data
}

func (s *dockerService) collectAndSaveLogs(ctx context.Context) {

	hostId := os.Getenv("HOST_ID")

	if hostId == "" {
		hostId = "unknown_host"
	}

	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: false}) // sadece running olanlar
	if err != nil {
		log.Printf("Container listesi alınamadı: %v", err)
		return
	}

	since := time.Now().Add(-5 * time.Second).Unix() // son 5 saniyenin logları

	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		out, err := s.cli.ContainerLogs(ctx, c.ID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Since:      fmt.Sprintf("%d", since),
			Timestamps: true,
		})
		if err != nil {
			log.Printf("Log alınamadı [%s]: %v", name, err)
			continue
		}

		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			line := scanner.Text()

			if len(line) < 8 {
				continue
			}

			// Docker log stream formatı: ilk 8 byte header, geri kalanı mesaj
			message := strings.TrimSpace(line[8:])
			if message == "" {
				continue
			}

			point := influxdb2.NewPointWithMeasurement("container_logs").
				AddTag("host_id", hostId).
				AddTag("metric_type", "container_log").
				AddTag("container_name", name).
				AddField("message", message).
				SetTime(time.Now())

			if err := s.influxClient.WriteAPI.WritePoint(ctx, point); err != nil {
				log.Printf("InfluxDB log yazma hatası [%s]: %v", name, err)
			}
		}

		out.Close()
	}
}
