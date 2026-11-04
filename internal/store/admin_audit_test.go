package store

import (
	"context"
	"testing"
)

func TestAdminAuditLogger(t *testing.T) {
	logger := NewAdminAuditLogger()
	ctx := context.Background()

	rec := logger.Record(ctx, "admin@stellar.org", ActionTenantSuspend, "tenant-12", "non-payment")
	if rec.ID != 1 || rec.Action != ActionTenantSuspend {
		t.Fatalf("unexpected record: %+v", rec)
	}

	recent := logger.ListRecent(10)
	if len(recent) != 1 {
		t.Errorf("expected 1 recent record, got %d", len(recent))
	}
}
