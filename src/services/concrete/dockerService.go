package servicesConcrete

import (
	"context"
	modelsDTOs "server-watcher-app/src/models/dtos"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type dockerService struct {
	cli *client.Client
}

func NewDockerService() (servicesAbstarct.DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &dockerService{cli: cli}, nil
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
