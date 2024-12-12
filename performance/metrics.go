package performance

import (
	"log"
	"sync"
	"time"
)

type Metrics struct {
	StartTime  time.Time
	EndTime    time.Time
	Operations int
	mu         sync.Mutex
}

func StartMetrics() *Metrics {
	return &Metrics{
		StartTime:  time.Now(),
		Operations: 0,
	}
}

func (m *Metrics) IncrementOperation() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Operations++
}

func (m *Metrics) Stop() {
	m.EndTime = time.Now()
}

func (m *Metrics) Report() {
	duration := m.EndTime.Sub(m.StartTime)
	throughput := float64(m.Operations) / duration.Seconds()

	log.Println("====== PERFORMANCE METRICS ======")
	log.Printf("Total Time: %s\n", duration)
	log.Printf("Operations Completed: %d\n", m.Operations)
	log.Printf("Throughput: %.2f ops/sec\n", throughput)
	log.Println("=================================")
}
