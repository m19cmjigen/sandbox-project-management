package repository

import (
	"context"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
)

// ProjectRepository はプロジェクトリポジトリのインターフェース
type ProjectRepository interface {
	// FindAll は全てのプロジェクトを取得
	FindAll(ctx context.Context) ([]domain.Project, error)

	// FindByID はIDでプロジェクトを取得
	FindByID(ctx context.Context, id int64) (*domain.Project, error)

	// GetByKey はプロジェクトキー(例: PROJ)で取得
	GetByKey(ctx context.Context, key string) (*domain.Project, error)

	// FindByJiraProjectID はJiraプロジェクトIDで取得
	FindByJiraProjectID(ctx context.Context, jiraProjectID string) (*domain.Project, error)

	// FindByOrganizationID は組織IDでプロジェクトを取得
	FindByOrganizationID(ctx context.Context, organizationID int64) ([]domain.Project, error)

	// FindUnassigned は未分類プロジェクト（organization_id が NULL）を取得
	FindUnassigned(ctx context.Context) ([]domain.Project, error)

	// FindWithStats は統計情報付きでプロジェクトを取得
	FindWithStats(ctx context.Context, id int64) (*domain.ProjectWithStats, error)

	// FindAllWithStats は全プロジェクトを統計情報付きで取得
	FindAllWithStats(ctx context.Context) ([]domain.ProjectWithStats, error)

	// Create は新しいプロジェクトを作成
	Create(ctx context.Context, project *domain.Project) error

	// Update はプロジェクトを更新
	Update(ctx context.Context, project *domain.Project) error

	// Delete はプロジェクトを削除
	Delete(ctx context.Context, id int64) error

	// AssignToOrganization はプロジェクトを組織に紐付け
	AssignToOrganization(ctx context.Context, projectID int64, organizationID *int64) error

	// ExistsByJiraProjectID はJiraプロジェクトIDで存在チェック
	ExistsByJiraProjectID(ctx context.Context, jiraProjectID string) (bool, error)
}
