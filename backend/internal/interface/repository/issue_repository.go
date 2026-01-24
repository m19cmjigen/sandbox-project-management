package repository

import (
	"context"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
)

// IssueRepository はチケットリポジトリのインターフェース
type IssueRepository interface {
	// FindAll は全てのチケットを取得
	FindAll(ctx context.Context) ([]domain.Issue, error)

	// FindByID はIDでチケットを取得
	FindByID(ctx context.Context, id int64) (*domain.Issue, error)

	// FindByJiraIssueID はJiraチケットIDで取得
	FindByJiraIssueID(ctx context.Context, jiraIssueID string) (*domain.Issue, error)

	// GetByJiraKey はJiraキー(例: PROJ-123)でチケットを取得
	GetByJiraKey(ctx context.Context, jiraKey string) (*domain.Issue, error)

	// FindByProjectID はプロジェクトIDでチケットを取得
	FindByProjectID(ctx context.Context, projectID int64) ([]domain.Issue, error)

	// FindByFilter はフィルタ条件でチケットを取得
	FindByFilter(ctx context.Context, filter domain.IssueFilter) ([]domain.Issue, error)

	// FindByDelayStatus は遅延ステータスでチケットを取得
	FindByDelayStatus(ctx context.Context, status domain.DelayStatus) ([]domain.Issue, error)

	// CountByProjectID はプロジェクトIDでチケット数を取得
	CountByProjectID(ctx context.Context, projectID int64) (int, error)

	// CountByDelayStatus は遅延ステータスでチケット数を取得
	CountByDelayStatus(ctx context.Context, projectID int64, status domain.DelayStatus) (int, error)

	// Create は新しいチケットを作成
	Create(ctx context.Context, issue *domain.Issue) error

	// Update はチケットを更新
	Update(ctx context.Context, issue *domain.Issue) error

	// Delete はチケットを削除
	Delete(ctx context.Context, id int64) error

	// BulkUpsert は複数のチケットを一括登録/更新
	BulkUpsert(ctx context.Context, issues []domain.Issue) error

	// ExistsByJiraIssueID はJiraチケットIDで存在チェック
	ExistsByJiraIssueID(ctx context.Context, jiraIssueID string) (bool, error)
}
