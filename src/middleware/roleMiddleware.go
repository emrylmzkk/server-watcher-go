package middleware

import (
	"log"
	"server-watcher-app/src/generic"
	"server-watcher-app/src/models"
	enumModels "server-watcher-app/src/models/enum"

	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(allowedRoles ...enumModels.UserRole) fiber.Handler {

	return func(c *fiber.Ctx) error {

		user, ok := c.Locals("currentUser").(*models.User)

		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(
				generic.NewErrorResponse("Forbidden", ""),
			)
		}

		for _, role := range allowedRoles {

			log.Println("Tercih edilen rol -->", allowedRoles)

			log.Println("Senin rol -->", role)
			if user.UserRole == role {
				log.Println("Kosul saglaniyor...")
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(
			generic.NewErrorResponse("Forbidden", "Insufficient permissions"),
		)

	}

}
