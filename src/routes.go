package src

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *AppContainer) {

	api := app.Group("/api/v1")

	projects := api.Group("/projects")
	projects.Post("/sync", container.ProjectHandler.SyncNow)
	projects.Post("/createpm2", container.ProjectHandler.AddNewPm2Project)
	projects.Post("/action", container.ProjectHandler.ControlProcess)
	projects.Delete("/:id", container.ProjectHandler.DeletePm2Project)
	projects.Get("/pm2-projects", container.ProjectHandler.GetAllPm2Projects)

}
