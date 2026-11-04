package store

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type attributionContextKey struct{}

// QueryAttribution identifies the originator of a database query with bounded cardinality.
type QueryAttribution struct {
	RoutePattern string // e.g., "/v1/events/:id", normalized without path parameters
	WorkerName   string // e.g., "ingester", "pruner", "archiver", "reconciler"
	TenantID     string // Opaque hash or UUID (never raw user names or contract code)
}

// WithQueryAttribution embeds attribution metadata into the request/worker context.
func WithQueryAttribution(ctx context.Context, attr QueryAttribution) context.Context {
	return context.WithValue(ctx, attributionContextKey{}, attr)
}

// QueryAttributionFromContext extracts attribution metadata from context, or returns default.
func QueryAttributionFromContext(ctx context.Context) (QueryAttribution, bool) {
	attr, ok := ctx.Value(attributionContextKey{}).(QueryAttribution)
	return attr, ok
}

var safeIdentRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-/:.]+$`)

// SanitizeTag ensures a tag value contains only safe characters and cannot break SQL comments.
func SanitizeTag(val string) string {
	val = strings.ReplaceAll(val, "/*", "")
	val = strings.ReplaceAll(val, "*/", "")
	val = strings.TrimSpace(val)
	if len(val) > 64 {
		val = val[:64]
	}
	if !safeIdentRegex.MatchString(val) {
		return "unknown"
	}
	return val
}

// AttachQueryTag prepends an attribution comment tag to the SQL query statement.
// Example format: /* op=route route=/v1/events/:id tenant=t-8f92 worker=none */ SELECT ...
func AttachQueryTag(query string, attr QueryAttribution) string {
	route := SanitizeTag(attr.RoutePattern)
	if route == "" {
		route = "none"
	}
	worker := SanitizeTag(attr.WorkerName)
	if worker == "" {
		worker = "none"
	}
	tenant := SanitizeTag(attr.TenantID)
	if tenant == "" {
		tenant = "none"
	}

	comment := fmt.Sprintf("/* route=%s worker=%s tenant=%s */ ", route, worker, tenant)
	return comment + query
}
