package middleware

import (
	"server-watcher-app/src/generic"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", "missing or malformed token"))
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		claims, err := generic.ValidateToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", err.Error()))
		}

		// Sonraki handler'larda kullanmak için locals'a yaz
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}
