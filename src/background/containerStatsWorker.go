package background

import (
	"context"
	servicesAbstarct "server-watcher-app/src/services/abstract"
)

type ContainerStatsWorker struct {
	dockerService servicesAbstarct.DockerService
}

func NewContainerStatsWorker(dockerService servicesAbstarct.DockerService) *ContainerStatsWorker {
	return &ContainerStatsWorker{
		dockerService: dockerService,
	}
}

func (w *ContainerStatsWorker) Start(ctx context.Context) {

	w.dockerService.StartStatsCollector(ctx)
}

