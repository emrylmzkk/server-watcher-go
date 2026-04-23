package src

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container *AppContainer) {

	api := app.Group("/api/v1")

	projects := api.Group("/projects")
	projects.Post("/sync", container.ProjectHandler.SyncNow)

}
