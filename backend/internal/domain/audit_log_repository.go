package domain

import (
	"context"
	"time"
)

// AuditLogRepository defines the interface for audit log data access
type AuditLogRepository interface {
	// Create creates a new audit log entry
	Create(ctx context.Context, log *AuditLog) error

	// GetByID retrieves an audit log entry by ID
	GetByID(ctx context.Context, id int64) (*AuditLog, error)

	// List retrieves audit logs with optional filtering
	List(ctx context.Context, filter *AuditLogFilter) ([]*AuditLog, error)

	// Count counts audit logs matching the filter
	Count(ctx context.Context, filter *AuditLogFilter) (int64, error)

	// DeleteOlderThan deletes audit logs older than the specified date
	DeleteOlderThan(ctx context.Context, olderThan *time.Time) (int64, error)
}
