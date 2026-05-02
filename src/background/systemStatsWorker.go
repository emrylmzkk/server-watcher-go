package background

import (
	"context"
	"log"
	servicesAbstarct "server-watcher-app/src/services/abstract"
	"time"
)

type SystemStatsWorker struct {
	statsService servicesAbstarct.IServerGeneralService
}

func NewSystemStatsWatcher(statsService servicesAbstarct.IServerGeneralService) *SystemStatsWorker {
	return &SystemStatsWorker{
		statsService: statsService,
	}
}

func (w *SystemStatsWorker) Start(ctx context.Context) {

	go func() {

		ticker := time.NewTicker(20 * time.Second)

		log.Println("System notification worker calisiyor...")

		for {
			select {
			case <-ticker.C:

				_, err := w.statsService.GetStatsForNotification(ctx)

				if err != nil {
					continue
				}

			case <-ctx.Done():
				return

			}
		}

	}()

}
