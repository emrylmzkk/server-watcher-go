package handler

import (
	"server-watcher-app/src/generic"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"github.com/gofiber/fiber/v2"
)

type ContainerStatHandler struct {
	statService servicesAbstarct.ContainerStatsService
}

func NewContainerStatHandler(statService servicesAbstarct.ContainerStatsService) *ContainerStatHandler {
	return &ContainerStatHandler{
		statService: statService,
	}
}

func (h *ContainerStatHandler) GetStatsByName(c *fiber.Ctx) error {

	req, err := generic.ParseParam[string](c, "containerName")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.statService.GetStatsByName(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *ContainerStatHandler) GetStatsByNameP(c *fiber.Ctx) error {

	pagination, err := generic.ParseQuery(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Bad Request", err))
	}

	res, err := h.statService.GetStatsByNameP(c.Context(), *pagination)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))

	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))
}
