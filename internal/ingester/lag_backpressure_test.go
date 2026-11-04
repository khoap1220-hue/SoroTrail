package ingester

import (
	"testing"
	"time"
)

func TestLagBackpressureThrottling(t *testing.T) {
	ctrl := NewLagBackpressureController(100)

	// 0 lag: no throttle
	if d := ctrl.UpdateLag(1000, 1000); d != 0 {
		t.Errorf("expected 0 throttle, got %v", d)
	}

	// 150 lag (> 100 max): max throttle 500ms
	if d := ctrl.UpdateLag(1150, 1000); d != 500*time.Millisecond {
		t.Errorf("expected 500ms throttle, got %v", d)
	}
}
