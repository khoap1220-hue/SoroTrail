package ingester

import (
	"sync"
	"time"
)

// LagBackpressureController dynamically throttles ingestion rate when replication lag rises.
type LagBackpressureController struct {
	mu             sync.RWMutex
	maxLagThreshold int64
	currentLag      int64
	throttleDelay   time.Duration
}

func NewLagBackpressureController(maxLag int64) *LagBackpressureController {
	return &LagBackpressureController{
		maxLagThreshold: maxLag,
		throttleDelay:   0,
	}
}

func (c *LagBackpressureController) UpdateLag(headLedger, ingestedLedger int64) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.currentLag = headLedger - ingestedLedger
	if c.currentLag <= 0 {
		c.throttleDelay = 0
	} else if c.currentLag > c.maxLagThreshold {
		c.throttleDelay = 500 * time.Millisecond
	} else {
		ratio := float64(c.currentLag) / float64(c.maxLagThreshold)
		c.throttleDelay = time.Duration(ratio * float64(200*time.Millisecond))
	}
	return c.throttleDelay
}
