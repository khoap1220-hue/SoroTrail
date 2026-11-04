package store

import (
	"context"
	"strings"
	"testing"
)

func TestAttachQueryTag(t *testing.T) {
	attr := QueryAttribution{
		RoutePattern: "/v1/events/:id",
		WorkerName:   "api-handler",
		TenantID:     "tenant-9941",
	}

	sql := "SELECT id, contract_id FROM events WHERE id = $1"
	tagged := AttachQueryTag(sql, attr)

	expectedPrefix := "/* route=/v1/events/:id worker=api-handler tenant=tenant-9941 */ "
	if !strings.HasPrefix(tagged, expectedPrefix) {
		t.Errorf("tagged query prefix mismatch:\ngot:  %s\nwant: %s", tagged, expectedPrefix)
	}
}

func TestQueryAttributionSanitization(t *testing.T) {
	// Assert comment injection characters are stripped
	malicious := QueryAttribution{
		RoutePattern: "/v1/events/* injection */",
		WorkerName:   "worker; DROP TABLE events;",
		TenantID:     "safe-tenant",
	}

	tagged := AttachQueryTag("SELECT 1", malicious)
	if strings.Contains(tagged, "/* injection") {
		t.Errorf("detected comment injection in tagged query: %s", tagged)
	}
}

func TestQueryAttributionContextRoundtrip(t *testing.T) {
	orig := QueryAttribution{
		RoutePattern: "/v1/stats",
		WorkerName:   "pruner",
		TenantID:     "system",
	}

	ctx := WithQueryAttribution(context.Background(), orig)
	recovered, ok := QueryAttributionFromContext(ctx)
	if !ok {
		t.Fatal("failed retrieving attribution from context")
	}

	if recovered != orig {
		t.Errorf("recovered attribution mismatch: got %+v, want %+v", recovered, orig)
	}
}
