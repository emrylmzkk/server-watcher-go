package src

import (
	"log"
	"server-watcher-app/src/background"
	"server-watcher-app/src/handler"
	"server-watcher-app/src/middleware"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	servicesConcrete "server-watcher-app/src/services/concrete"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AppContainer struct {
	DockerHandler        *handler.DockerHandler
	Pm2Handler           *handler.Pm2Handler
	ServerGeneralHandler *handler.ServerGeneralHandler
	AuthHandler          *handler.AuthHandler

	AuthMiddleware       fiber.Handler
	SyncWorker           *background.SyncWorker
	ContainerStatsWorker *background.ContainerStatsWorker
}

func NewAppContainer(db *gorm.DB) *AppContainer {

	projectRepo := repositoryConcrete.NewSqliteRepository(db)
	pm2ProjectRepo := repositoryConcrete.NewProjectRepository(db)
	userRepository := repositoryConcrete.NewUserRepository(db)

	dockerProv, err := servicesConcrete.NewDockerProvider()

	if err != nil {
		log.Printf("Docker provider could not be started: %v", err)
	}

	dockerService, err := servicesConcrete.NewDockerService(db)

	if err != nil {
		log.Printf("Docker service could not be started: %v", err)
	}

	serverGeneralService := servicesConcrete.NewServerGeneralService()
	authService := servicesConcrete.NewAuthService(userRepository)

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
		DockerHandler:        handler.NewDockerHandler(dockerService),
		Pm2Handler:           handler.NewPm2Handler(pm2Service),
		ServerGeneralHandler: handler.NewServerGeneralHandler(serverGeneralService),
		AuthHandler:          handler.NewAuthHandler(authService),
		AuthMiddleware:       middleware.AuthMiddleware(userRepository),

		SyncWorker:           syncWorker,
		ContainerStatsWorker: background.NewContainerStatsWorker(dockerService),
	}

}
