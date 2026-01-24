package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

// syncLogRepository はSyncLogRepositoryのPostgreSQL実装
type syncLogRepository struct {
	db *sqlx.DB
}

// NewSyncLogRepository は新しいSyncLogRepositoryを作成
func NewSyncLogRepository(db *sqlx.DB) repository.SyncLogRepository {
	return &syncLogRepository{db: db}
}

// FindAll は全ての同期ログを取得
func (r *syncLogRepository) FindAll(ctx context.Context, limit int) ([]domain.SyncLog, error) {
	var logs []domain.SyncLog
	query := `
		SELECT id, sync_type, executed_at, completed_at, status,
		       projects_synced, issues_synced, error_message, duration_seconds
		FROM sync_logs
		ORDER BY executed_at DESC
		LIMIT $1
	`
	err := r.db.SelectContext(ctx, &logs, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find all sync logs: %w", err)
	}
	return logs, nil
}

// FindByID はIDで同期ログを取得
func (r *syncLogRepository) FindByID(ctx context.Context, id int64) (*domain.SyncLog, error) {
	var log domain.SyncLog
	query := `
		SELECT id, sync_type, executed_at, completed_at, status,
		       projects_synced, issues_synced, error_message, duration_seconds
		FROM sync_logs
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &log, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find sync log by id %d: %w", id, err)
	}
	return &log, nil
}

// FindByType は同期タイプでログを取得
func (r *syncLogRepository) FindByType(ctx context.Context, syncType domain.SyncType, limit int) ([]domain.SyncLog, error) {
	var logs []domain.SyncLog
	query := `
		SELECT id, sync_type, executed_at, completed_at, status,
		       projects_synced, issues_synced, error_message, duration_seconds
		FROM sync_logs
		WHERE sync_type = $1
		ORDER BY executed_at DESC
		LIMIT $2
	`
	err := r.db.SelectContext(ctx, &logs, query, syncType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find sync logs by type %s: %w", syncType, err)
	}
	return logs, nil
}

// FindByStatus はステータスでログを取得
func (r *syncLogRepository) FindByStatus(ctx context.Context, status domain.SyncStatus) ([]domain.SyncLog, error) {
	var logs []domain.SyncLog
	query := `
		SELECT id, sync_type, executed_at, completed_at, status,
		       projects_synced, issues_synced, error_message, duration_seconds
		FROM sync_logs
		WHERE status = $1
		ORDER BY executed_at DESC
	`
	err := r.db.SelectContext(ctx, &logs, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to find sync logs by status %s: %w", status, err)
	}
	return logs, nil
}

// FindLatest は最新の同期ログを取得
func (r *syncLogRepository) FindLatest(ctx context.Context) (*domain.SyncLog, error) {
	var log domain.SyncLog
	query := `
		SELECT id, sync_type, executed_at, completed_at, status,
		       projects_synced, issues_synced, error_message, duration_seconds
		FROM sync_logs
		ORDER BY executed_at DESC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &log, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find latest sync log: %w", err)
	}
	return &log, nil
}

// FindLatestByType はタイプ別の最新同期ログを取得
func (r *syncLogRepository) FindLatestByType(ctx context.Context, syncType domain.SyncType) (*domain.SyncLog, error) {
	var log domain.SyncLog
	query := `
		SELECT id, sync_type, executed_at, completed_at, status,
		       projects_synced, issues_synced, error_message, duration_seconds
		FROM sync_logs
		WHERE sync_type = $1
		ORDER BY executed_at DESC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &log, query, syncType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find latest sync log by type %s: %w", syncType, err)
	}
	return &log, nil
}

// Create は新しい同期ログを作成
func (r *syncLogRepository) Create(ctx context.Context, log *domain.SyncLog) error {
	query := `
		INSERT INTO sync_logs (sync_type, executed_at, status, projects_synced, issues_synced)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		log.SyncType,
		log.ExecutedAt,
		log.Status,
		log.ProjectsSynced,
		log.IssuesSynced,
	).Scan(&log.ID)

	if err != nil {
		return fmt.Errorf("failed to create sync log: %w", err)
	}
	return nil
}

// Update は同期ログを更新
func (r *syncLogRepository) Update(ctx context.Context, log *domain.SyncLog) error {
	query := `
		UPDATE sync_logs
		SET completed_at = $1, status = $2, projects_synced = $3,
		    issues_synced = $4, error_message = $5, duration_seconds = $6
		WHERE id = $7
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		log.CompletedAt,
		log.Status,
		log.ProjectsSynced,
		log.IssuesSynced,
		log.ErrorMessage,
		log.DurationSeconds,
		log.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update sync log: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sync log not found: id=%d", log.ID)
	}

	return nil
}

// UpdateStatus はステータスを更新
func (r *syncLogRepository) UpdateStatus(ctx context.Context, id int64, status domain.SyncStatus, errorMessage *string) error {
	query := `
		UPDATE sync_logs
		SET status = $1, error_message = $2, completed_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	result, err := r.db.ExecContext(ctx, query, status, errorMessage, id)
	if err != nil {
		return fmt.Errorf("failed to update sync log status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("sync log not found: id=%d", id)
	}

	return nil
}
