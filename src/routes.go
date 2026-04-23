package src

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *AppContainer) {

	api := app.Group("/api/v1")

	pm2 := api.Group("/pm2")
	pm2.Post("/create", container.Pm2Handler.CreatePm2Project)
	pm2.Post("/start/:id", container.Pm2Handler.StartPm2Project)
	pm2.Post("/stop/:id", container.Pm2Handler.StopPm2Project)
	pm2.Delete("/:id", container.Pm2Handler.DeletePm2Project)
	pm2.Get("/projects", container.Pm2Handler.GetAllPm2Projects)
	pm2.Put("/:id", container.Pm2Handler.UpdatePm2Project)
	pm2.Post("/reset", container.Pm2Handler.ResetPm2Process)

	docker := api.Group("/docker")
	docker.Get("/containers", container.DockerHandler.GetAllContainers)
	docker.Post("/start/:id", container.DockerHandler.StartContainer)
	docker.Post("/stop/:id", container.DockerHandler.StopContainer)

}
