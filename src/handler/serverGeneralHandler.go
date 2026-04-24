package handler

import (
	"server-watcher-app/src/generic"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"github.com/gofiber/fiber/v2"
)

type ServerGeneralHandler struct {
	serverGeneralService servicesAbstarct.IServerGeneralService
}

func NewServerGeneralHandler(serverGeneralService servicesAbstarct.IServerGeneralService) *ServerGeneralHandler {
	return &ServerGeneralHandler{
		serverGeneralService: serverGeneralService,
	}
}

func (h *ServerGeneralHandler) GetCPUPercent(c *fiber.Ctx) error {

	res, err := h.serverGeneralService.GetCPUPercent(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *ServerGeneralHandler) GetRamStats(c *fiber.Ctx) error {

	used, total, percent, err := h.serverGeneralService.GetRamStats(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(map[string]float64{
		"used":    used,
		"total":   total,
		"percent": percent,
	}))

}

func (h *ServerGeneralHandler) GetDiskStats(c *fiber.Ctx) error {

	used, total, percent, err := h.serverGeneralService.GetDiskStats(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(map[string]float64{
		"used":    used,
		"total":   total,
		"percent": percent,
	}))

}

func (h *ServerGeneralHandler) GetGeneralStats(c *fiber.Ctx) error {

	stats, err := h.serverGeneralService.GetSystemStats(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(stats))

}
