package ingester

import (
	"math"
	"time"
)

// IngestionCapacityModel calculates event processing capacity and backpressure runway.
type IngestionCapacityModel struct {
	MaxWorkers         int
	EventsPerSecWorker float64
	BufferCapacity     int
}

// NewIngestionCapacityModel initializes a capacity model with verified parameters.
func NewIngestionCapacityModel(maxWorkers int, eventsPerSecWorker float64, bufferCapacity int) *IngestionCapacityModel {
	return &IngestionCapacityModel{
		MaxWorkers:         maxWorkers,
		EventsPerSecWorker: eventsPerSecWorker,
		BufferCapacity:     bufferCapacity,
	}
}

// EstimateRunway calculates the seconds before the ingestion buffer fills at current arrival rate.
func (m *IngestionCapacityModel) EstimateRunway(arrivalRate float64, currentBuffered int) (time.Duration, bool) {
	processingCapacity := float64(m.MaxWorkers) * m.EventsPerSecWorker
	netArrival := arrivalRate - processingCapacity

	if netArrival <= 0 {
		return time.Duration(math.MaxInt64), false // Steady state: no backlog accumulation
	}

	remainingBuffer := float64(m.BufferCapacity - currentBuffered)
	if remainingBuffer <= 0 {
		return 0, true
	}

	seconds := remainingBuffer / netArrival
	return time.Duration(seconds * float64(time.Second)), true
}
