package data

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Metrics struct {
	Id          int           `json:"id"`
	InstanceId  int           `json:"instance_id"`
	CpuUsage    float32       `json:"cpu_usage"`
	MemoryUsage float32       `json:"memory_usage"`
	Uptime      time.Duration `json:"uptime"`
	RecordedAt  time.Time     `json:"-"`
}

type MetricsModel struct {
	DB *pgxpool.Pool
}

func (m *MetricsModel) Insert(metrics *Metrics) error {
	query := `
		INSERT INTO instances (instance_id, cpu_usage, memory_usage, uptime)
		VALUES ($1, $2, $3, $4)
		RETURNING id, recorded_at`

	args := []any{metrics.InstanceId, metrics.CpuUsage, metrics.MemoryUsage, metrics.Uptime}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, args...).Scan(&metrics.Id, &metrics.RecordedAt)
}
