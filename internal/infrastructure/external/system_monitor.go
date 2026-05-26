package external

import (
	"time"

	"tsniper/internal/application/ports"

	"tsniper/pkg/httpx"
)

var _ ports.SystemMonitor = (*SystemMonitor)(nil)

type SystemMonitor struct {
	client    *httpx.Client
	startTime int64
}

func NewSystemMonitor() *SystemMonitor {
	return &SystemMonitor{
		client: httpx.NewClient(
			httpx.WithRetryPolicy(httpx.NoRetry()),
			httpx.WithBackoffPolicy(httpx.NoBackoff()),
		),
		startTime: time.Now().Unix(),
	}
}

func (p *SystemMonitor) StartTime() int64 {
	return p.startTime
}

// todo: improve latency measurement
func (p *SystemMonitor) RutgersSISLatency() int64 {
	start := time.Now()

	resp, err := p.client.Get("https://sis.rutgers.edu/soc/")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	return duration.Milliseconds()
}

// todo: improve latency measurement
func (p *SystemMonitor) RutgersSIMLatency() int64 {
	start := time.Now()

	resp, err := p.client.Get("https://sims.rutgers.edu/webreg/")
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	return duration.Milliseconds()
}
