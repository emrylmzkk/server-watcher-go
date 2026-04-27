package background

import (
	"context"
	"log"
	"os"
	servicesAbstarct "server-watcher-app/src/services/abstract"
)

type InitHelper struct {
	authService servicesAbstarct.AuthService
}

func NewInitService(authService servicesAbstarct.AuthService) *InitHelper {
	return &InitHelper{
		authService: authService,
	}
}

func (h *InitHelper) CreateAdminUser(ctx context.Context) error {

	userName := string(os.Getenv("ADMIN_USERNAME"))
	password := string(os.Getenv("ADMIN_PASSWORD"))

	// log.Println("username -->", userName)
	// log.Println("adminpassword --> ", password)

	err := h.authService.CreateAdminUser(ctx, userName, password)

	if err != nil {
		log.Println("create admin user error", err)
		return nil
	}

	return nil

}
