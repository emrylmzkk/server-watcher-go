package src

import (
	"log"
	"os"
	"server-watcher-app/src/background"
	genericInfluxDB "server-watcher-app/src/generic/influxDB"
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
	//ContainerStatHandler *handler.ContainerStatHandler
	PublicHandler *handler.PublicHandler
	MetricHandler *handler.MetricViewerHandler

	InitHelper           *background.InitHelper
	AuthMiddleware       fiber.Handler
	SyncWorker           *background.SyncWorker
	ContainerStatsWorker *background.ContainerStatsWorker
	InfluxDBClient       *genericInfluxDB.InfluxClient
}

func NewAppContainer(db *gorm.DB) *AppContainer {

	fluxBucketName := os.Getenv("FLUX_BUCKET_NAME")

	if fluxBucketName == "" {
		fluxBucketName = "metrics"
	}

	influxClient, err := genericInfluxDB.NewInfluxClient(genericInfluxDB.InfluxConfig{
		URL:    "http://217.76.49.19:8086",
		Token:  "my-super-token",
		Org:    "my-org",
		Bucket: fluxBucketName,
	})

	if err != nil {
		log.Fatal("InfluxDB connection failed:", err)
	}

	projectRepo := repositoryConcrete.NewSqliteRepository(db)
	pm2ProjectRepo := repositoryConcrete.NewProjectRepository(db)
	userRepository := repositoryConcrete.NewUserRepository(db)
	//containerStatRepo := repositoryConcrete.NewContainerStatLogRepo(db)

	dockerProv, err := servicesConcrete.NewDockerProvider()

	if err != nil {
		log.Printf("Docker provider could not be started: %v", err)
	}

	dockerService, err := servicesConcrete.NewDockerService(db, influxClient)

	if err != nil {
		log.Printf("Docker service could not be started: %v", err)
	}

	serverGeneralService := servicesConcrete.NewServerGeneralService()
	authService := servicesConcrete.NewAuthService(userRepository)
	//containerStatService := servicesConcrete.NewContainerStatService(containerStatRepo)

	pm2Prov := servicesConcrete.NewPM2Provider(pm2ProjectRepo)

	providers := []servicesAbstarct.ProcessProvider{
		pm2Prov,
		dockerProv,
	}

	metricViewerService := servicesConcrete.NewMetricViewerService(influxClient)

	wathcerService := servicesConcrete.NewWatcherService(projectRepo, providers)
	pm2Service := servicesConcrete.NewPm2ProjectService(pm2ProjectRepo)

	syncWorker := background.NewSyncWorker(wathcerService, 30*time.Second)

	initHelper := background.NewInitService(authService)

	return &AppContainer{
		//ProjectHandler: handler.NewProjectsHandler(wathcerService, pm2Service),
		DockerHandler:        handler.NewDockerHandler(dockerService),
		Pm2Handler:           handler.NewPm2Handler(pm2Service),
		ServerGeneralHandler: handler.NewServerGeneralHandler(serverGeneralService),
		AuthHandler:          handler.NewAuthHandler(authService),
		//ContainerStatHandler: handler.NewContainerStatHandler(containerStatService),
		MetricHandler: handler.NewMetricViewerHandler(metricViewerService),

		PublicHandler: handler.NewPublicHandler(),

		AuthMiddleware: middleware.AuthMiddleware(userRepository),

		InitHelper:           initHelper,
		SyncWorker:           syncWorker,
		ContainerStatsWorker: background.NewContainerStatsWorker(dockerService),
		InfluxDBClient:       influxClient,
	}

}
