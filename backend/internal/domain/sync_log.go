package domain

import (
	"database/sql"
	"time"
)

// SyncType は同期タイプの型
type SyncType string

const (
	SyncTypeFull  SyncType = "FULL"
	SyncTypeDelta SyncType = "DELTA"
)

// SyncStatus は同期ステータスの型
type SyncStatus string

const (
	SyncStatusRunning SyncStatus = "RUNNING"
	SyncStatusSuccess SyncStatus = "SUCCESS"
	SyncStatusFailure SyncStatus = "FAILURE"
)

// SyncLog はバッチ実行ログを表すエンティティ
type SyncLog struct {
	ID              int64          `db:"id" json:"id"`
	SyncType        SyncType       `db:"sync_type" json:"sync_type"`
	ExecutedAt      time.Time      `db:"executed_at" json:"executed_at"`
	CompletedAt     sql.NullTime   `db:"completed_at" json:"completed_at"`
	Status          SyncStatus     `db:"status" json:"status"`
	ProjectsSynced  int            `db:"projects_synced" json:"projects_synced"`
	IssuesSynced    int            `db:"issues_synced" json:"issues_synced"`
	ErrorMessage    *string        `db:"error_message" json:"error_message"`
	DurationSeconds *int           `db:"duration_seconds" json:"duration_seconds"`
}
