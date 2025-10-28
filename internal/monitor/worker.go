package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/Kndrik/cloud-monitoring/internal/data"
)

type Worker struct {
	Instance *data.Instance
	Logger   *slog.Logger
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Instance.RefreshRate)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.Logger.Info("worker fetching data", "instance", w.Instance.Name)
			}
		}
	}()
}
