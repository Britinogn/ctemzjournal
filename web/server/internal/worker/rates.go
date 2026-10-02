package worker

import (
	"context"
	"log"
	"time"

	"github.com/britinogn/ctemzjournal/internal/rates"
)

// RunRates refreshes the rates cache on a timer until ctx is done.
// A failed refresh never stops the loop: the cache keeps serving the
// last values with stale: true (docx §9).
func RunRates(ctx context.Context, svc *rates.Service, interval time.Duration) error {
	if interval < time.Minute {
		interval = 20 * time.Minute
	}
	if err := svc.Refresh(ctx); err != nil {
		log.Printf("rates: initial refresh failed (serving stale): %v", err)
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := svc.Refresh(ctx); err != nil {
				log.Printf("rates: refresh failed (serving stale): %v", err)
			}
		}
	}
}
