package background

import (
	"context"
	"log"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"time"
)

type SyncWorker struct {
	service  servicesAbstarct.WatcherService
	interval time.Duration
}

func NewSyncWorker(service servicesAbstarct.WatcherService, interval time.Duration) *SyncWorker {
	return &SyncWorker{
		service:  service,
		interval: interval,
	}
}

func (w *SyncWorker) Start(ctx context.Context) {

	ticker := time.NewTicker(w.interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := w.service.SyncAll(ctx)
				log.Printf("Sync worker successfully completed the process")
				if err != nil {
					log.Printf("Sync worker service error")
				}

			case <-ctx.Done():
				log.Println("SyncWorker stopped successfully")
				return
			}
		}
	}()

}
