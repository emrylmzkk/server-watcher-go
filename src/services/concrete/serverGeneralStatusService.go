package servicesConcrete

import (
	"context"
	"errors"
	"fmt"
	"server-watcher-app/src/generic"
	"server-watcher-app/src/models"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"strconv"
	"strings"
)

type serverGeneralService struct {
}

func NewServerGeneralService() servicesAbstarct.IServerGeneralService {
	return &serverGeneralService{}
}

func (s *serverGeneralService) GetCPUPercent(ctx context.Context) (float64, error) {

	//cmd := "top -bn1 | grep 'Cpu(s)'"

	cmd := generic.NewCmd(ctx, "sh", "-c", "top -bn1 | grep 'Cpu(s)'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0.0, errors.New(string(output))
	}

	parts := strings.Split(string(output), ",")
	for _, part := range parts {
		if strings.Contains(part, "id") {
			trimmed := strings.TrimSpace(part)
			valStr := strings.Split(trimmed, " ")[0]
			idle, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return 0, err
			}
			return float64(int((100.0-idle)*10) / 10), nil // Yuvarlama
		}
	}

	return 0.0, errors.New("cpu percentage not found")

}

func (s *serverGeneralService) GetRamStats(ctx context.Context) (used, total, percent float64, err error) {

	cmd := generic.NewCmd(ctx, "sh", "-c", "free -m")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return 0, 0, 0, errors.New(string(output))
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return 0, 0, 0, fmt.Errorf("RAM parse error")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 3 {
		return 0, 0, 0, fmt.Errorf("RAM parse error")
	}

	totalMB, _ := strconv.ParseFloat(fields[1], 64)
	usedMB, _ := strconv.ParseFloat(fields[2], 64)

	total = totalMB / 1024 // MB to GB
	used = usedMB / 1024
	percent = (used / total) * 100

	return used, total, percent, nil

}

func (s *serverGeneralService) GetDiskStats(ctx context.Context) (used, total, percent float64, err error) {

	cmd := generic.NewCmd(ctx, "sh", "-c", "df -h /")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 2 {
		return 0, 0, 0, fmt.Errorf("Disk parse error")
	}

	fields := strings.Fields(lines[1])
	if len(fields) < 5 {
		return 0, 0, 0, fmt.Errorf("Disk parse error")
	}

	// G ve % işaretlerini temizleme
	total, _ = strconv.ParseFloat(strings.ReplaceAll(fields[1], "G", ""), 64)
	used, _ = strconv.ParseFloat(strings.ReplaceAll(fields[2], "G", ""), 64)
	percent, _ = strconv.ParseFloat(strings.ReplaceAll(fields[4], "%", ""), 64)

	return used, total, percent, nil

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
