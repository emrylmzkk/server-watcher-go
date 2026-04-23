package src

import (
	"log"
	"server-watcher-app/src/handler"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	servicesConcrete "server-watcher-app/src/services/concrete"

	"gorm.io/gorm"
)

type AppContainer struct {
	ProjectHandler *handler.ProjectsHandler
}

func NewAppContainer(db *gorm.DB) *AppContainer {

	projectRepo := repositoryConcrete.NewSqliteRepository(db)

	dockerProv, err := servicesConcrete.NewDockerProvider()

	if err != nil {
		log.Printf("Docker provider could not be started", err)
	}

	pm2Prov := servicesConcrete.NewPM2Provider()

	providers := []servicesAbstarct.ProcessProvider{
		pm2Prov,
		dockerProv,
	}

	wathcerService := servicesConcrete.NewWatcherService(projectRepo, providers)

	return &AppContainer{
		ProjectHandler: handler.NewProjectsHandler(wathcerService),
	}

}
