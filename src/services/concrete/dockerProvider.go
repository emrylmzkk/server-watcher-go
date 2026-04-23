package servicesConcrete

import (
	"context"
	"server-watcher-app/src/models"
	enumModels "server-watcher-app/src/models/enum"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type dockerProvider struct {
	cli *client.Client
}

func NewDockerProvider() (servicesAbstarct.ProcessProvider, error) {
	// Host üzerindeki docker.sock'a bağlanır
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &dockerProvider{cli: cli}, nil
}

func (d *dockerProvider) GetProviderType() enumModels.ProcessType {
	return enumModels.Docker
}

func (d *dockerProvider) ListProcesses(ctx context.Context) ([]models.MonitoredEntity, error) {
	containers, err := d.cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var entities []models.MonitoredEntity
	for _, c := range containers {

		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		entities = append(entities, models.MonitoredEntity{
			ExternalID: c.ID,
			Name:       name,
			Type:       enumModels.Docker,
			Status:     c.State,
		})
	}
	return entities, nil
}

func (d *dockerProvider) StopProcess(ctx context.Context, id string) error {
	return d.cli.ContainerStop(ctx, id, container.StopOptions{})
}

func (d *dockerProvider) StartProcess(ctx context.Context, id string) error {
	return d.cli.ContainerStart(ctx, id, container.StartOptions{})
}
