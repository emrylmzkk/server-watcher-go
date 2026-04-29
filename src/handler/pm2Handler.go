package handler

import (
	"server-watcher-app/src/generic"
	modelsDTOs "server-watcher-app/src/models/dtos"
	servicesAbstarct "server-watcher-app/src/services/abstract"

	"github.com/gofiber/fiber/v2"
)

type Pm2Handler struct {
	pm2Service servicesAbstarct.Pm2ProjectService
}

func NewPm2Handler(pm2Service servicesAbstarct.Pm2ProjectService) *Pm2Handler {
	return &Pm2Handler{
		pm2Service: pm2Service,
	}
}

func (h *Pm2Handler) CreatePm2Project(c *fiber.Ctx) error {

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

func (h *Pm2Handler) GetPm2ProjectById(c *fiber.Ctx) error {

	req, err := generic.ParseParam[int](c, "id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.pm2Service.GetPm2ProjectById(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *Pm2Handler) DeletePm2Project(c *fiber.Ctx) error {

	req, err := generic.ParseParam[int](c, "id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.pm2Service.DeletePm2Project(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *Pm2Handler) ClearAndDeletePm2Project(c *fiber.Ctx) error {

	req, err := generic.ParseParam[int](c, "id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.pm2Service.ClearAndDeletePm2Project(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *Pm2Handler) GetAllPm2Projects(c *fiber.Ctx) error {

	res, err := h.pm2Service.GetPm2Projects(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *Pm2Handler) StartPm2Project(c *fiber.Ctx) error {

	req, err := generic.ParseParam[int](c, "id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.pm2Service.StartPm2Project(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *Pm2Handler) StopPm2Project(c *fiber.Ctx) error {

	req, err := generic.ParseParam[int](c, "id")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.pm2Service.StopPm2Project(c.Context(), req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *Pm2Handler) UpdatePm2Project(c *fiber.Ctx) error {

	id, err := generic.ParseParam[int](c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	req, err := generic.ParseBody[modelsDTOs.UpdatePm2ProjectRequestDTO](c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(generic.NewErrorResponse("Invalid request body", err.Error()))
	}

	res, err := h.pm2Service.UpdatePm2Project(c.Context(), id, req)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *Pm2Handler) ResetPm2Process(c *fiber.Ctx) error {

	res, err := h.pm2Service.ResetPm2Process(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *Pm2Handler) GetPm2InsideList(c *fiber.Ctx) error {

	res, err := h.pm2Service.GetPm2InsideList(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}

func (h *Pm2Handler) SyncPm2Projects(c *fiber.Ctx) error {

	res, err := h.pm2Service.SyncPm2Projects(c.Context())

	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(generic.NewErrorResponse("Server Error", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse(res))

}
