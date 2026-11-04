package store

import (
	"context"
	"sync"
	"time"
)

// AdminActionType classifies sensitive administrative operations.
type AdminActionType string

const (
	ActionRoleChange        AdminActionType = "ROLE_CHANGE"
	ActionTenantSuspend     AdminActionType = "TENANT_SUSPEND"
	ActionWebhookRotate     AdminActionType = "WEBHOOK_ROTATE"
	ActionPrunePolicyUpdate AdminActionType = "PRUNE_POLICY_UPDATE"
)

// AdminActionRecord represents an append-only audit log entry.
type AdminActionRecord struct {
	ID        int64           `json:"id"`
	Timestamp time.Time       `json:"timestamp"`
	Actor     string          `json:"actor"`
	Action    AdminActionType `json:"action"`
	TargetID  string          `json:"target_id"`
	Details   string          `json:"details"`
}

// AdminAuditLogger records immutable administrative operations.
type AdminAuditLogger struct {
	mu      sync.RWMutex
	records []AdminActionRecord
}

func NewAdminAuditLogger() *AdminAuditLogger {
	return &AdminAuditLogger{records: make([]AdminActionRecord, 0)}
}

func (l *AdminAuditLogger) Record(ctx context.Context, actor string, action AdminActionType, target, details string) AdminActionRecord {
	l.mu.Lock()
	defer l.mu.Unlock()

	rec := AdminActionRecord{
		ID:        int64(len(l.records) + 1),
		Timestamp: time.Now().UTC(),
		Actor:     actor,
		Action:    action,
		TargetID:  target,
		Details:   details,
	}
	l.records = append(l.records, rec)
	return rec
}

func (l *AdminAuditLogger) ListRecent(limit int) []AdminActionRecord {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.records) <= limit {
		res := make([]AdminActionRecord, len(l.records))
		copy(res, l.records)
		return res
	}
	res := make([]AdminActionRecord, limit)
	copy(res, l.records[len(l.records)-limit:])
	return res
}
