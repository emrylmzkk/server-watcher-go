package handler

import (
	"server-watcher-app/src/generic"
	modelsDTOs "server-watcher-app/src/models/dtos"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"github.com/gofiber/fiber/v2"
)

type UserPreferencesHandler struct {
	service servicesAbstarct.IUserPreferencesService
}

func NewUserPreferencesHandler(service servicesAbstarct.IUserPreferencesService) *UserPreferencesHandler {
	return &UserPreferencesHandler{
		service: service,
	}
}

func (h *UserPreferencesHandler) CreateServerStatPreference(c *fiber.Ctx) error {

	currentUser := generic.GetCurrentUser(c)

	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", nil))
	}

	req, err := generic.ParseBody[modelsDTOs.ServerStatusSettingDTO](c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.service.CreateServerStatPreference(c.Context(), currentUser.ID, req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *UserPreferencesHandler) ClearServerStatPreference(c *fiber.Ctx) error {

	currentUser := generic.GetCurrentUser(c)

	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", nil))
	}

	res, err := h.service.ClearServerStatPreference(c.Context(), currentUser.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *UserPreferencesHandler) GetUserSSTSettings(c *fiber.Ctx) error {

	currentUser := generic.GetCurrentUser(c)

	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", nil))
	}

	res, err := h.service.GetUserSSTSettings(c.Context(), currentUser.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}
