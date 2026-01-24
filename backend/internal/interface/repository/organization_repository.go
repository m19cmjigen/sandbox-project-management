package repository

import (
	"context"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
)

// OrganizationRepository は組織リポジトリのインターフェース
type OrganizationRepository interface {
	// FindAll は全ての組織を取得
	FindAll(ctx context.Context) ([]domain.Organization, error)

	// FindByID はIDで組織を取得
	FindByID(ctx context.Context, id int64) (*domain.Organization, error)

	// FindByParentID は親組織IDで子組織を取得
	FindByParentID(ctx context.Context, parentID int64) ([]domain.Organization, error)

	// FindRoots はルート組織（parent_id が NULL）を取得
	FindRoots(ctx context.Context) ([]domain.Organization, error)

	// FindByPath はパスで組織を検索（階層検索用）
	FindByPath(ctx context.Context, pathPrefix string) ([]domain.Organization, error)

	// Create は新しい組織を作成
	Create(ctx context.Context, org *domain.Organization) error

	// Update は組織を更新
	Update(ctx context.Context, org *domain.Organization) error

	// Delete は組織を削除
	Delete(ctx context.Context, id int64) error

	// ExistsByID は組織が存在するかチェック
	ExistsByID(ctx context.Context, id int64) (bool, error)

	// HasChildren は子組織が存在するかチェック
	HasChildren(ctx context.Context, id int64) (bool, error)
}
