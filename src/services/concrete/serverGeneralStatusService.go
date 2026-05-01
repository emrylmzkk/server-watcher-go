package servicesConcrete

import (
	"context"
	"log"
	"math"
	"server-watcher-app/src/models"
	enumNotification "server-watcher-app/src/models/enum/notification"
	repositoryConcrete "server-watcher-app/src/repository/concrete"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type serverGeneralService struct {
	notificationsService servicesAbstarct.INotificationService
	userRepo             repositoryConcrete.UserRepository
}

func NewServerGeneralService(notificationsService servicesAbstarct.INotificationService, userRepo repositoryConcrete.UserRepository) servicesAbstarct.IServerGeneralService {
	return &serverGeneralService{
		notificationsService: notificationsService,
		userRepo:             userRepo,
	}
}

func (s *serverGeneralService) GetCPUPercent(ctx context.Context) (float64, error) {

	percentages, err := cpu.PercentWithContext(ctx, 500*time.Millisecond, false)

	if err != nil {
		return 0, err
	}

	if len(percentages) == 0 {
		return 0, nil
	}

	val := math.Round(percentages[0]*10) / 10

	if val > 10.0 {
		log.Printf("Buraya girdi")
		user, _ := s.userRepo.GetAdminUser(ctx)

		log.Printf("Gelen kullanici", user)

		if user.FCMToken != nil {
			s.notificationsService.SendToUser(user.ID, *user.FCMToken, "Yüksek CPU Uyarısı!", "Sunucu CPU kullanımı sınırı aştı", enumNotification.TypeCPUUsage)
		}

	}

	return val, nil

}

func (s *serverGeneralService) GetDiskStats(ctx context.Context) (used, total, percent float64, err error) {

	d, err := disk.UsageWithContext(ctx, "/")

	if err != nil {
		return 0, 0, 0, err
	}

	totalGB := float64(d.Total) / (1024 * 1024 * 1024)
	usedGB := float64(d.Used) / (1024 * 1024 * 1024)

	return math.Round(usedGB*100) / 100, math.Round(totalGB*100) / 100, math.Round(d.UsedPercent*10) / 10, nil

}

func (s *serverGeneralService) GetRamStats(ctx context.Context) (used, total, percent float64, err error) {

	v, err := mem.VirtualMemoryWithContext(ctx)

	if err != nil {
		return 0, 0, 0, err
	}

	totalGB := float64(v.Total) / (1024 * 1024 * 1024)
	usedGB := float64(v.Used) / (1024 * 1024 * 1024)

	return math.Round(usedGB*100) / 100, math.Round(totalGB*100) / 100, math.Round(v.UsedPercent*10) / 10, nil

}

func (s *serverGeneralService) GetSystemStats(ctx context.Context) (*models.SystemStats, error) {
	cpu, err := s.GetCPUPercent(ctx)
	if err != nil {
		return nil, err
	}

	rUsed, rTotal, rPercent, err := s.GetRamStats(ctx)
	if err != nil {
		return nil, err
	}

	dUsed, dTotal, dPercent, err := s.GetDiskStats(ctx)
	if err != nil {
		return nil, err
	}

	return &models.SystemStats{
		CPUPercent:  cpu,
		RAMUsedGB:   rUsed,
		RAMTotalGB:  rTotal,
		RAMPercent:  rPercent,
		DiskUsedGB:  dUsed,
		DiskTotalGB: dTotal,
		DiskPercent: dPercent,
	}, nil
}
