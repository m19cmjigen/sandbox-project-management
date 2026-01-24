package repository

import (
	"context"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
)

// SyncLogRepository は同期ログリポジトリのインターフェース
type SyncLogRepository interface {
	// FindAll は全ての同期ログを取得
	FindAll(ctx context.Context, limit int) ([]domain.SyncLog, error)

	// FindByID はIDで同期ログを取得
	FindByID(ctx context.Context, id int64) (*domain.SyncLog, error)

	// FindByType は同期タイプでログを取得
	FindByType(ctx context.Context, syncType domain.SyncType, limit int) ([]domain.SyncLog, error)

	// FindByStatus はステータスでログを取得
	FindByStatus(ctx context.Context, status domain.SyncStatus) ([]domain.SyncLog, error)

	// FindLatest は最新の同期ログを取得
	FindLatest(ctx context.Context) (*domain.SyncLog, error)

	// FindLatestByType はタイプ別の最新同期ログを取得
	FindLatestByType(ctx context.Context, syncType domain.SyncType) (*domain.SyncLog, error)

	// Create は新しい同期ログを作成
	Create(ctx context.Context, log *domain.SyncLog) error

	// Update は同期ログを更新
	Update(ctx context.Context, log *domain.SyncLog) error

	// UpdateStatus はステータスを更新
	UpdateStatus(ctx context.Context, id int64, status domain.SyncStatus, errorMessage *string) error
}
