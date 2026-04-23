package handler

import (
	"server-watcher-app/src/generic"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"github.com/gofiber/fiber/v2"
)

type DockerHandler struct {
	dockerService servicesAbstarct.DockerService
}

func NewDockerHandler(dockerService servicesAbstarct.DockerService) *DockerHandler {
	return &DockerHandler{
		dockerService: dockerService,
	}
}

func (h *DockerHandler) GetAllContainers(c *fiber.Ctx) error {
	res, err := h.dockerService.GetContainers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Docker Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))
}

func (h *DockerHandler) StartContainer(c *fiber.Ctx) error {

	req, err := generic.ParseParam[string](c, "id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	err = h.dockerService.StartContainer(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Docker Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse("Container started"))
}

func (h *DockerHandler) StopContainer(c *fiber.Ctx) error {

	req, err := generic.ParseParam[string](c, "id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	err = h.dockerService.StopContainer(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Docker Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse("Container stopped"))
}
