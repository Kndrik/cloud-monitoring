package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/Kndrik/cloud-monitoring/internal/data"
)

type Worker struct {
	instance     *data.Instance
	logger       *slog.Logger
	metricsModel *data.MetricsModel
}

func NewWorker(instance *data.Instance, logger *slog.Logger, metricsModel *data.MetricsModel) *Worker {
	return &Worker{
		instance:     instance,
		logger:       logger,
		metricsModel: metricsModel,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.instance.RefreshRate)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				metrics := GenerateFakeMetrics(w.instance)
				w.logger.Info("[worker] inserting metrics")
				err := w.metricsModel.Insert(metrics)
				if err != nil {
					w.logger.Error("failed to insert metric", "instance", w.instance.Name, "error", err)
				}
			}
		}
	}()
}
