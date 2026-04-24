package middleware

import (
	"server-watcher-app/src/generic"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(userRepo repositoryConcrete.UserRepository) fiber.Handler {
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

		user, err := userRepo.GetByID(c.Context(), int(claims.UserID))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(generic.NewErrorResponse("Unauthorized", "User not found"))
		}

		// Sonraki handler'larda kullanmak için locals'a yaz
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("currentUser", user)

		return c.Next()
	}
}
