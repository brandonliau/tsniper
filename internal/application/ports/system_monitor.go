package ports

type SystemMonitor interface {
	StartTime() int64
	RutgersSISLatency() int64
}
