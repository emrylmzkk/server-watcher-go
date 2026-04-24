package generic

import (
	"context"
	"server-watcher-app/src/models"
	repositoryConcrete "server-watcher-app/src/repository/concrete"

	"github.com/gofiber/fiber/v2"
)

func GetCurrentUser(c *fiber.Ctx) *models.User {
	user, ok := c.Locals("currentUser").(*models.User)

	if !ok {
		return nil
	}

	return user
}

func IsAdmin(c context.Context, userId uint, userRepo repositoryConcrete.UserRepository) bool {

	isAdmin, err := userRepo.IsUserAdmin(c, userId)

	if err != nil {
		return false
	}

	if isAdmin {
		return true
	}

	return false
}
