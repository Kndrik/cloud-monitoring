package monitor

import (
	"context"
	"log/slog"
	"strconv"
	"sync"

	"github.com/Kndrik/cloud-monitoring/internal/data"
)

type Scheduler struct {
	logger *slog.Logger
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	models *data.Models
}

func NewScheduler(models *data.Models, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		models: models,
		logger: logger,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	instances, err := s.models.Instances.GetAll()
	if err != nil {
		return err
	}

	for _, instance := range instances {
		w := Worker{Instance: instance, Logger: s.logger}
		s.wg.Go(func() {
			w.Run(s.ctx)
		})

	}

	s.logger.Info("started workers", "amount", strconv.Itoa(len(instances)))

	return nil
}

func (s *Scheduler) Stop() {
	s.logger.Info("shutting down all workers")
	s.cancel()
	s.wg.Wait()
	s.logger.Info("all workers have been shut down")
}
