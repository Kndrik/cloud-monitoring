package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/Kndrik/cloud-monitoring/internal/data"
)

type Worker struct {
	Instance     *data.Instance
	Logger       *slog.Logger
	MetricsModel *data.MetricsModel
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.Instance.RefreshRate)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metrics := GenerateFakeMetrics(w.Instance)
				w.Logger.Info("[worker] inserting metrics")
				err := w.MetricsModel.Insert(metrics)
				if err != nil {
					w.Logger.Error("failed to insert metric", "instance", w.Instance.Name, "error", err)
				}
			}
		}
	}()
}
