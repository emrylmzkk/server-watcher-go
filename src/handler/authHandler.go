package handler

import (
	"server-watcher-app/src/generic"
	modelsDTOs "server-watcher-app/src/models/dtos"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService servicesAbstarct.AuthService
}

func NewAuthHandler(authService servicesAbstarct.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {

	req, err := generic.ParseBody[modelsDTOs.RegisterRequestDTO](c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.authService.Register(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *AuthHandler) Login(c *fiber.Ctx) error {

	req, err := generic.ParseBody[modelsDTOs.LoginRequestDTO](c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.authService.Login(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {

	req, err := generic.ParseBody[modelsDTOs.RefreshTokenRequestDTO](c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.authService.RefreshToken(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *AuthHandler) GetCurrentUserInformation(c *fiber.Ctx) error {

	currentUser := generic.GetCurrentUser(c)

	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", nil))
	}

	res, err := h.authService.GetCurrentUserInformation(c.Context(), int(currentUser.ID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *AuthHandler) TakeUserFCMToken(c *fiber.Ctx) error {

	currentUser := generic.GetCurrentUser(c)

	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", nil))
	}

	req, err := generic.ParseBody[modelsDTOs.UserFCMTokenRequestDTO](c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.authService.TakeUserFCMToken(c.Context(), int(currentUser.ID), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *AuthHandler) RemoveUserFCMToken(c *fiber.Ctx) error {

	currentUser := generic.GetCurrentUser(c)

	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", nil))
	}

	res, err := h.authService.RemoveUserFCMToken(c.Context(), int(currentUser.ID))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *AuthHandler) GetAllUsers(c *fiber.Ctx) error {

	currentUser := generic.GetCurrentUser(c)

	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", nil))
	}

	res, err := h.authService.GetAllUser(c.Context(), currentUser.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}
