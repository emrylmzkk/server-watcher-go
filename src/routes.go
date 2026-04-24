package src

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *AppContainer) {

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", container.AuthHandler.Register)
	auth.Post("/login", container.AuthHandler.Login)
	auth.Post("/refresh-token", container.AuthHandler.RefreshToken)
	auth.Get("/me", container.AuthMiddleware, container.AuthHandler.GetCurrentUserInformation)

	pm2 := api.Group("/pm2", container.AuthMiddleware)
	pm2.Post("/create", container.Pm2Handler.CreatePm2Project)
	pm2.Post("/start/:id", container.Pm2Handler.StartPm2Project)
	pm2.Post("/stop/:id", container.Pm2Handler.StopPm2Project)
	pm2.Delete("/:id", container.Pm2Handler.DeletePm2Project)
	pm2.Get("/projects", container.Pm2Handler.GetAllPm2Projects)
	pm2.Put("/:id", container.Pm2Handler.UpdatePm2Project)
	pm2.Post("/reset", container.Pm2Handler.ResetPm2Process)
	pm2.Get("/inside-list", container.Pm2Handler.GetPm2InsideList)

	docker := api.Group("/docker", container.AuthMiddleware)
	docker.Get("/containers", container.DockerHandler.GetAllContainers)
	docker.Post("/start/:id", container.DockerHandler.StartContainer)
	docker.Post("/stop/:id", container.DockerHandler.StopContainer)

	serverGeneral := api.Group("/server-general", container.AuthMiddleware)
	serverGeneral.Get("/cpu", container.ServerGeneralHandler.GetCPUPercent)
	serverGeneral.Get("/ram", container.ServerGeneralHandler.GetRamStats)
	serverGeneral.Get("/disk", container.ServerGeneralHandler.GetDiskStats)
	serverGeneral.Get("/stats", container.ServerGeneralHandler.GetGeneralStats)

}
