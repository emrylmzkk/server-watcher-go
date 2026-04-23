package handler

import (
	"server-watcher-app/src/generic"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"github.com/gofiber/fiber/v2"
)

type ProjectsHandler struct {
	wathcerService servicesAbstarct.WatcherService
}

func NewProjectsHandler(wathcerService servicesAbstarct.WatcherService) *ProjectsHandler {
	return &ProjectsHandler{
		wathcerService: wathcerService,
	}
}

func (h *ProjectsHandler) SyncNow(c *fiber.Ctx) error {

	err := h.wathcerService.SyncAll(c.Context())

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	response := map[string]interface{}{
		"status": "synchronized",
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(response))

}
