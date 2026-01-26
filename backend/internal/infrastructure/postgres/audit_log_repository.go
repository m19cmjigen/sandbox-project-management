package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
)

type auditLogRepository struct {
	db *sqlx.DB
}

// NewAuditLogRepository creates a new PostgreSQL audit log repository
func NewAuditLogRepository(db *sqlx.DB) domain.AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			user_id, username, action, resource_type, resource_id,
			method, path, ip_address, user_agent, request_body,
			response_status, response_body, error_message, duration_ms
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		log.UserID,
		log.Username,
		log.Action,
		log.ResourceType,
		log.ResourceID,
		log.Method,
		log.Path,
		log.IPAddress,
		log.UserAgent,
		log.RequestBody,
		log.ResponseStatus,
		log.ResponseBody,
		log.ErrorMessage,
		log.DurationMs,
	).Scan(&log.ID, &log.CreatedAt)
}

func (r *auditLogRepository) GetByID(ctx context.Context, id int64) (*domain.AuditLog, error) {
	var log domain.AuditLog
	query := `
		SELECT id, user_id, username, action, resource_type, resource_id,
		       method, path, ip_address, user_agent, request_body,
		       response_status, response_body, error_message, duration_ms, created_at
		FROM audit_logs
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &log, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log by ID: %w", err)
	}

	return &log, nil
}

func (r *auditLogRepository) List(ctx context.Context, filter *domain.AuditLogFilter) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, user_id, username, action, resource_type, resource_id,
		       method, path, ip_address, user_agent, request_body,
		       response_status, response_body, error_message, duration_ms, created_at
		FROM audit_logs
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.UserID != nil {
			query += fmt.Sprintf(" AND user_id = $%d", argCount)
			args = append(args, *filter.UserID)
			argCount++
		}

		if filter.Username != nil {
			query += fmt.Sprintf(" AND LOWER(username) = LOWER($%d)", argCount)
			args = append(args, *filter.Username)
			argCount++
		}

		if filter.Action != nil {
			query += fmt.Sprintf(" AND action = $%d", argCount)
			args = append(args, *filter.Action)
			argCount++
		}

		if filter.ResourceType != nil {
			query += fmt.Sprintf(" AND resource_type = $%d", argCount)
			args = append(args, *filter.ResourceType)
			argCount++
		}

		if filter.ResourceID != nil {
			query += fmt.Sprintf(" AND resource_id = $%d", argCount)
			args = append(args, *filter.ResourceID)
			argCount++
		}

		if filter.Method != nil {
			query += fmt.Sprintf(" AND method = $%d", argCount)
			args = append(args, *filter.Method)
			argCount++
		}

		if filter.StartDate != nil {
			query += fmt.Sprintf(" AND created_at >= $%d", argCount)
			args = append(args, *filter.StartDate)
			argCount++
		}

		if filter.EndDate != nil {
			query += fmt.Sprintf(" AND created_at <= $%d", argCount)
			args = append(args, *filter.EndDate)
			argCount++
		}
	}

	query += " ORDER BY created_at DESC"

	if filter != nil {
		if filter.Limit > 0 {
			query += fmt.Sprintf(" LIMIT $%d", argCount)
			args = append(args, filter.Limit)
			argCount++
		}

		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET $%d", argCount)
			args = append(args, filter.Offset)
			argCount++
		}
	}

	var logs []*domain.AuditLog
	err := r.db.SelectContext(ctx, &logs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}

	return logs, nil
}

func (r *auditLogRepository) Count(ctx context.Context, filter *domain.AuditLogFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM audit_logs WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.UserID != nil {
			query += fmt.Sprintf(" AND user_id = $%d", argCount)
			args = append(args, *filter.UserID)
			argCount++
		}

		if filter.Username != nil {
			query += fmt.Sprintf(" AND LOWER(username) = LOWER($%d)", argCount)
			args = append(args, *filter.Username)
			argCount++
		}

		if filter.Action != nil {
			query += fmt.Sprintf(" AND action = $%d", argCount)
			args = append(args, *filter.Action)
			argCount++
		}

		if filter.ResourceType != nil {
			query += fmt.Sprintf(" AND resource_type = $%d", argCount)
			args = append(args, *filter.ResourceType)
			argCount++
		}

		if filter.ResourceID != nil {
			query += fmt.Sprintf(" AND resource_id = $%d", argCount)
			args = append(args, *filter.ResourceID)
			argCount++
		}

		if filter.Method != nil {
			query += fmt.Sprintf(" AND method = $%d", argCount)
			args = append(args, *filter.Method)
			argCount++
		}

		if filter.StartDate != nil {
			query += fmt.Sprintf(" AND created_at >= $%d", argCount)
			args = append(args, *filter.StartDate)
			argCount++
		}

		if filter.EndDate != nil {
			query += fmt.Sprintf(" AND created_at <= $%d", argCount)
			args = append(args, *filter.EndDate)
			argCount++
		}
	}

	var count int64
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return count, nil
}

func (r *auditLogRepository) DeleteOlderThan(ctx context.Context, olderThan *time.Time) (int64, error) {
	if olderThan == nil {
		return 0, fmt.Errorf("olderThan date is required")
	}

	query := `DELETE FROM audit_logs WHERE created_at < $1`

	result, err := r.db.ExecContext(ctx, query, *olderThan)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old audit logs: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
