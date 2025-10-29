package monitor

import (
	"math/rand/v2"
	"time"

	"github.com/Kndrik/cloud-monitoring/internal/data"
)

func GenerateFakeMetrics(instance *data.Instance) *data.Metrics {
	return &data.Metrics{
		InstanceId:  instance.Id,
		CpuUsage:    rand.Float64() * 100,
		MemoryUsage: rand.Float64() * 100,
		Uptime:      time.Duration(rand.IntN(72)) * time.Hour,
	}
}
