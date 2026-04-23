package src

import (
	"log"
	"server-watcher-app/src/background"
	"server-watcher-app/src/handler"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	servicesConcrete "server-watcher-app/src/services/concrete"
	"time"

	"gorm.io/gorm"
)

type AppContainer struct {
	DockerHandler *handler.DockerHandler
	Pm2Handler    *handler.Pm2Handler

	SyncWorker *background.SyncWorker
}

func NewAppContainer(db *gorm.DB) *AppContainer {

	projectRepo := repositoryConcrete.NewSqliteRepository(db)
	pm2ProjectRepo := repositoryConcrete.NewProjectRepository(db)

	dockerProv, err := servicesConcrete.NewDockerProvider()

	if err != nil {
		log.Printf("Docker provider could not be started: %v", err)
	}

	dockerService, err := servicesConcrete.NewDockerService()
	if err != nil {
		log.Printf("Docker service could not be started: %v", err)
	}

	pm2Prov := servicesConcrete.NewPM2Provider(pm2ProjectRepo)

	providers := []servicesAbstarct.ProcessProvider{
		pm2Prov,
		dockerProv,
	}

	wathcerService := servicesConcrete.NewWatcherService(projectRepo, providers)
	pm2Service := servicesConcrete.NewPm2ProjectService(pm2ProjectRepo)

	syncWorker := background.NewSyncWorker(wathcerService, 30*time.Second)

	return &AppContainer{
		//ProjectHandler: handler.NewProjectsHandler(wathcerService, pm2Service),
		DockerHandler: handler.NewDockerHandler(dockerService),
		Pm2Handler:    handler.NewPm2Handler(pm2Service),

		SyncWorker: syncWorker,
	}

}
