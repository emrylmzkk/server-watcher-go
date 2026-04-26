package handler

import (
	"server-watcher-app/src/generic"

	"github.com/gofiber/fiber/v2"
)

type PublicHandler struct {
}

func NewPublicHandler() *PublicHandler {
	return &PublicHandler{}
}

func (h *PublicHandler) CheckServerCondition(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(generic.NewSuccessResponse("Server is already running"))

}
