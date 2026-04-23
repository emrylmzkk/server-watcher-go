package servicesConcrete

import (
	"context"
	"encoding/json"
	"os/exec"
	"server-watcher-app/src/models"
	enumModels "server-watcher-app/src/models/enum"
	servicesAbstarct "server-watcher-app/src/services/abstract"
)

type pm2Provider struct{}

func NewPM2Provider() servicesAbstarct.ProcessProvider {
	return &pm2Provider{}
}

func (p *pm2Provider) GetProviderType() enumModels.ProcessType {
	return enumModels.PM2
}

// PM2'den gelen JSON'u karşılamak için geçici struct
type pm2Process struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	PmID   int    `json:"pm_id"`
}

func (p *pm2Provider) ListProcesses(ctx context.Context) ([]models.MonitoredEntity, error) {
	// pm2 jlist komutu tüm süreçleri JSON döner
	cmd := exec.CommandContext(ctx, "pm2", "jlist")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var rawProcesses []pm2Process
	if err := json.Unmarshal(output, &rawProcesses); err != nil {
		return nil, err
	}

	var entities []models.MonitoredEntity
	for _, proc := range rawProcesses {
		entities = append(entities, models.MonitoredEntity{
			ExternalID: proc.Name, // PM2'de genelde isim üzerinden yönetmek daha kolaydır
			Name:       proc.Name,
			Type:       enumModels.PM2,
			Status:     proc.Status,
		})
	}
	return entities, nil
}

func (p *pm2Provider) StopProcess(ctx context.Context, id string) error {
	return exec.CommandContext(ctx, "pm2", "stop", id).Run()
}

func (p *pm2Provider) StartProcess(ctx context.Context, id string) error {
	return exec.CommandContext(ctx, "pm2", "start", id).Run()
}
