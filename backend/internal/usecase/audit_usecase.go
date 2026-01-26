package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
)

// AuditUsecase defines audit logging use cases
type AuditUsecase interface {
	// LogAction logs a user action
	LogAction(ctx context.Context, req domain.AuditLogCreateRequest) error

	// GetAuditLog retrieves an audit log by ID
	GetAuditLog(ctx context.Context, id int64) (*domain.AuditLog, error)

	// ListAuditLogs retrieves audit logs with filtering
	ListAuditLogs(ctx context.Context, filter *domain.AuditLogFilter) ([]*domain.AuditLog, int64, error)

	// CleanupOldLogs deletes audit logs older than the specified retention period
	CleanupOldLogs(ctx context.Context, retentionDays int) (int64, error)
}

type auditUsecase struct {
	auditLogRepo domain.AuditLogRepository
}

// NewAuditUsecase creates a new audit usecase
func NewAuditUsecase(auditLogRepo domain.AuditLogRepository) AuditUsecase {
	return &auditUsecase{
		auditLogRepo: auditLogRepo,
	}
}

func (u *auditUsecase) LogAction(ctx context.Context, req domain.AuditLogCreateRequest) error {
	log := &domain.AuditLog{
		Action:       req.Action,
		ResourceType: req.ResourceType,
		Method:       req.Method,
		Path:         req.Path,
	}

	if req.UserID != nil {
		log.UserID = sql.NullInt64{Int64: *req.UserID, Valid: true}
	}

	if req.Username != nil {
		log.Username = sql.NullString{String: *req.Username, Valid: true}
	}

	if req.ResourceID != nil {
		log.ResourceID = sql.NullString{String: *req.ResourceID, Valid: true}
	}

	if req.IPAddress != nil {
		log.IPAddress = sql.NullString{String: *req.IPAddress, Valid: true}
	}

	if req.UserAgent != nil {
		log.UserAgent = sql.NullString{String: *req.UserAgent, Valid: true}
	}

	if req.RequestBody != nil {
		log.RequestBody = sql.NullString{String: *req.RequestBody, Valid: true}
	}

	if req.ResponseStatus != nil {
		log.ResponseStatus = sql.NullInt32{Int32: int32(*req.ResponseStatus), Valid: true}
	}

	if req.ResponseBody != nil {
		log.ResponseBody = sql.NullString{String: *req.ResponseBody, Valid: true}
	}

	if req.ErrorMessage != nil {
		log.ErrorMessage = sql.NullString{String: *req.ErrorMessage, Valid: true}
	}

	if req.DurationMs != nil {
		log.DurationMs = sql.NullInt32{Int32: int32(*req.DurationMs), Valid: true}
	}

	if err := u.auditLogRepo.Create(ctx, log); err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

func (u *auditUsecase) GetAuditLog(ctx context.Context, id int64) (*domain.AuditLog, error) {
	log, err := u.auditLogRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}
	if log == nil {
		return nil, fmt.Errorf("audit log not found")
	}
	return log, nil
}

func (u *auditUsecase) ListAuditLogs(ctx context.Context, filter *domain.AuditLogFilter) ([]*domain.AuditLog, int64, error) {
	// Set default limit if not specified
	if filter == nil {
		filter = &domain.AuditLogFilter{
			Limit: 100,
		}
	} else if filter.Limit == 0 {
		filter.Limit = 100
	}

	// Get total count
	count, err := u.auditLogRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Get logs
	logs, err := u.auditLogRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}

	return logs, count, nil
}

func (u *auditUsecase) CleanupOldLogs(ctx context.Context, retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		return 0, fmt.Errorf("retention days must be positive")
	}

	cutoffDate := time.Now().AddDate(0, 0, -retentionDays)
	deleted, err := u.auditLogRepo.DeleteOlderThan(ctx, &cutoffDate)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old audit logs: %w", err)
	}

	return deleted, nil
}
