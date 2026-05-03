package src

import (
	"log"
	"server-watcher-app/src/background"
	"server-watcher-app/src/generic"
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
	DockerHandler          *handler.DockerHandler
	Pm2Handler             *handler.Pm2Handler
	ServerGeneralHandler   *handler.ServerGeneralHandler
	AuthHandler            *handler.AuthHandler
	ContainerStatHandler   *handler.ContainerStatHandler
	UserPreferencesHandler *handler.UserPreferencesHandler

	PublicHandler *handler.PublicHandler

	InitHelper *background.InitHelper

	AuthMiddleware       fiber.Handler
	SyncWorker           *background.SyncWorker
	ContainerStatsWorker *background.ContainerStatsWorker
	SystemNotification   *background.SystemStatsWorker
}

func NewAppContainer(db *gorm.DB) *AppContainer {

	fcmClient := generic.InitFirebaseMessaging()

	projectRepo := repositoryConcrete.NewSqliteRepository(db)
	pm2ProjectRepo := repositoryConcrete.NewProjectRepository(db)
	userRepository := repositoryConcrete.NewUserRepository(db)
	containerStatRepo := repositoryConcrete.NewContainerStatLogRepo(db)
	userPreferencesRepo := repositoryConcrete.NewUserPreferencesRepository(db)

	dockerProv, err := servicesConcrete.NewDockerProvider()

	if err != nil {
		log.Printf("Docker provider could not be started: %v", err)
	}

	dockerService, err := servicesConcrete.NewDockerService(db)

	if err != nil {
		log.Printf("Docker service could not be started: %v", err)
	}

	notificationService := servicesConcrete.NewNotificationService(fcmClient)

	serverGeneralService := servicesConcrete.NewServerGeneralService(notificationService, userRepository)
	authService := servicesConcrete.NewAuthService(userRepository)
	containerStatService := servicesConcrete.NewContainerStatService(containerStatRepo)
	userPreferencesService := servicesConcrete.NewUserPreferencesService(userPreferencesRepo)

	pm2Prov := servicesConcrete.NewPM2Provider(pm2ProjectRepo)

	providers := []servicesAbstarct.ProcessProvider{
		pm2Prov,
		dockerProv,
	}

	wathcerService := servicesConcrete.NewWatcherService(projectRepo, providers)
	pm2Service := servicesConcrete.NewPm2ProjectService(pm2ProjectRepo)

	syncWorker := background.NewSyncWorker(wathcerService, 30*time.Second)

	initHelper := background.NewInitService(authService, pm2Service)

	return &AppContainer{
		//ProjectHandler: handler.NewProjectsHandler(wathcerService, pm2Service),
		DockerHandler:          handler.NewDockerHandler(dockerService),
		Pm2Handler:             handler.NewPm2Handler(pm2Service),
		ServerGeneralHandler:   handler.NewServerGeneralHandler(serverGeneralService),
		AuthHandler:            handler.NewAuthHandler(authService),
		ContainerStatHandler:   handler.NewContainerStatHandler(containerStatService),
		UserPreferencesHandler: handler.NewUserPreferencesHandler(userPreferencesService),

		PublicHandler: handler.NewPublicHandler(),

		AuthMiddleware: middleware.AuthMiddleware(userRepository),

		InitHelper:           initHelper,
		SyncWorker:           syncWorker,
		ContainerStatsWorker: background.NewContainerStatsWorker(dockerService),
		SystemNotification:   background.NewSystemStatsWatcher(serverGeneralService),
	}

}
