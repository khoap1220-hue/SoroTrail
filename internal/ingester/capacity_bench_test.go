package ingester

import (
	"testing"
	"time"
)

func TestIngestionCapacityModelRunway(t *testing.T) {
	// 5 workers @ 100 events/sec = 500 events/sec capacity. Buffer = 10,000.
	model := NewIngestionCapacityModel(5, 100.0, 10000)

	// Under normal load (400 arrivals/sec < 500 capacity), runway is unbounded
	_, backlogged := model.EstimateRunway(400.0, 1000)
	if backlogged {
		t.Error("expected no backlog under capacity")
	}

	// Under surge load (1500 arrivals/sec > 500 capacity), net arrival = +1000/sec
	// Remaining buffer: 10000 - 2000 = 8000 -> 8 seconds runway
	runway, backlogged := model.EstimateRunway(1500.0, 2000)
	if !backlogged {
		t.Error("expected backlog under surge")
	}

	expectedRunway := 8 * time.Second
	if runway != expectedRunway {
		t.Errorf("got runway %v, want %v", runway, expectedRunway)
	}
}
