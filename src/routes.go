package src

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *AppContainer) {

	api := app.Group("/api/v1")

	pm2 := api.Group("/pm2")
	//pm2.Post("/sync", container.ProjectHandler.SyncNow)
	pm2.Post("/create", container.Pm2Handler.CreatePm2Project)
	//pm2.Post("/action", container.Pm2Handler.ControlPm2Project)
	pm2.Delete("/:id", container.Pm2Handler.DeletePm2Project)
	pm2.Get("/projects", container.Pm2Handler.GetAllPm2Projects)

	docker := api.Group("/docker")
	docker.Get("/containers", container.DockerHandler.GetAllContainers)
	docker.Post("/start/:id", container.DockerHandler.StartContainer)
	docker.Post("/stop/:id", container.DockerHandler.StopContainer)

}
