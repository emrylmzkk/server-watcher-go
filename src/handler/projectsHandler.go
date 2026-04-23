package handler

import (
	"server-watcher-app/src/generic"
	modelsDTOs "server-watcher-app/src/models/dtos"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"github.com/gofiber/fiber/v2"
)

type ProjectsHandler struct {
	wathcerService servicesAbstarct.WatcherService
	pm2Service     servicesAbstarct.Pm2ProjectService
}

func NewProjectsHandler(wathcerService servicesAbstarct.WatcherService, pm2Service servicesAbstarct.Pm2ProjectService) *ProjectsHandler {
	return &ProjectsHandler{
		wathcerService: wathcerService,
		pm2Service:     pm2Service,
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

func (h *ProjectsHandler) ControlProcess(c *fiber.Ctx) error {

	req, err := generic.ParseBody[modelsDTOs.ActionOnProject](c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.wathcerService.ControlProcess(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *ProjectsHandler) AddNewPm2Project(c *fiber.Ctx) error {

	req, err := generic.ParseBody[modelsDTOs.CreatePm2ProjectRequestDTO](c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.pm2Service.AddNewProject(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

// func (h *ProjectsHandler) DeletePm2Project(c *fiber.Ctx) error {

// 	req, err := generic.ParseParam[int](c, "id")

// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
// 	}

// 	res, err := h.pm2Service.DeletePm2Project(c.Context(), req)

// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
// 	}

// 	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

// }

// func (h *ProjectsHandler) GetAllPm2Projects(c *fiber.Ctx) error {

// 	res, err := h.pm2Service.GetPm2Projects(c.Context())

// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
// 	}

// 	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

// }
