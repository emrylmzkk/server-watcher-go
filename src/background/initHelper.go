package background

import (
	"context"
	"log"
	"os"
	servicesAbstarct "server-watcher-app/src/services/abstract"
)

type InitHelper struct {
	authService servicesAbstarct.AuthService
	pm2Service  servicesAbstarct.Pm2ProjectService
}

func NewInitService(authService servicesAbstarct.AuthService, pm2Service servicesAbstarct.Pm2ProjectService) *InitHelper {
	return &InitHelper{
		authService: authService,
		pm2Service:  pm2Service,
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

func (h *InitHelper) CreateMockPm2Project(ctx context.Context) error {

	err := h.pm2Service.CreateExamplePm2Project(ctx)

	if err != nil {
		log.Println("create mock pm2 project error", err)
		return nil
	}

	return nil

}
