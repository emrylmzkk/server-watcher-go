package src

import (
	"server-watcher-app/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *AppContainer) {

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", container.AuthHandler.Register)
	auth.Post("/login", container.AuthHandler.Login)
	auth.Post("/refresh-token", container.AuthHandler.RefreshToken)

	pm2 := api.Group("/pm2", middleware.AuthMiddleware())
	pm2.Post("/create", container.Pm2Handler.CreatePm2Project)
	pm2.Post("/start/:id", container.Pm2Handler.StartPm2Project)
	pm2.Post("/stop/:id", container.Pm2Handler.StopPm2Project)
	pm2.Delete("/:id", container.Pm2Handler.DeletePm2Project)
	pm2.Get("/projects", container.Pm2Handler.GetAllPm2Projects)
	pm2.Put("/:id", container.Pm2Handler.UpdatePm2Project)
	pm2.Post("/reset", container.Pm2Handler.ResetPm2Process)

	docker := api.Group("/docker", middleware.AuthMiddleware())
	docker.Get("/containers", container.DockerHandler.GetAllContainers)
	docker.Post("/start/:id", container.DockerHandler.StartContainer)
	docker.Post("/stop/:id", container.DockerHandler.StopContainer)

	serverGeneral := api.Group("/server-general", middleware.AuthMiddleware())
	serverGeneral.Get("/cpu", container.ServerGeneralHandler.GetCPUPercent)
	serverGeneral.Get("/ram", container.ServerGeneralHandler.GetRamStats)
	serverGeneral.Get("/disk", container.ServerGeneralHandler.GetDiskStats)
	serverGeneral.Get("/stats", container.ServerGeneralHandler.GetGeneralStats)

}
