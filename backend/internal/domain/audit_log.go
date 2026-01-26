package domain

import (
	"database/sql"
	"time"
)

// AuditAction represents the type of action performed
type AuditAction string

const (
	AuditActionCreate AuditAction = "CREATE"
	AuditActionUpdate AuditAction = "UPDATE"
	AuditActionDelete AuditAction = "DELETE"
	AuditActionView   AuditAction = "VIEW"
	AuditActionLogin  AuditAction = "LOGIN"
	AuditActionLogout AuditAction = "LOGOUT"
	AuditActionSync   AuditAction = "SYNC"
)

// ResourceType represents the type of resource being acted upon
type ResourceType string

const (
	ResourceTypeUser         ResourceType = "user"
	ResourceTypeOrganization ResourceType = "organization"
	ResourceTypeProject      ResourceType = "project"
	ResourceTypeIssue        ResourceType = "issue"
	ResourceTypeSyncLog      ResourceType = "sync_log"
	ResourceTypeAuth         ResourceType = "auth"
	ResourceTypeDashboard    ResourceType = "dashboard"
)

// AuditLog represents an audit trail entry
type AuditLog struct {
	ID             int64          `db:"id" json:"id"`
	UserID         sql.NullInt64  `db:"user_id" json:"user_id,omitempty"`
	Username       sql.NullString `db:"username" json:"username,omitempty"`
	Action         AuditAction    `db:"action" json:"action"`
	ResourceType   ResourceType   `db:"resource_type" json:"resource_type"`
	ResourceID     sql.NullString `db:"resource_id" json:"resource_id,omitempty"`
	Method         string         `db:"method" json:"method"`
	Path           string         `db:"path" json:"path"`
	IPAddress      sql.NullString `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent      sql.NullString `db:"user_agent" json:"user_agent,omitempty"`
	RequestBody    sql.NullString `db:"request_body" json:"request_body,omitempty"`
	ResponseStatus sql.NullInt32  `db:"response_status" json:"response_status,omitempty"`
	ResponseBody   sql.NullString `db:"response_body" json:"response_body,omitempty"`
	ErrorMessage   sql.NullString `db:"error_message" json:"error_message,omitempty"`
	DurationMs     sql.NullInt32  `db:"duration_ms" json:"duration_ms,omitempty"`
	CreatedAt      time.Time      `db:"created_at" json:"created_at"`
}

// AuditLogFilter represents filtering options for audit logs
type AuditLogFilter struct {
	UserID       *int64
	Username     *string
	Action       *AuditAction
	ResourceType *ResourceType
	ResourceID   *string
	Method       *string
	StartDate    *time.Time
	EndDate      *time.Time
	Limit        int
	Offset       int
}

// AuditLogCreateRequest represents a request to create an audit log entry
type AuditLogCreateRequest struct {
	UserID         *int64
	Username       *string
	Action         AuditAction
	ResourceType   ResourceType
	ResourceID     *string
	Method         string
	Path           string
	IPAddress      *string
	UserAgent      *string
	RequestBody    *string
	ResponseStatus *int
	ResponseBody   *string
	ErrorMessage   *string
	DurationMs     *int
}
