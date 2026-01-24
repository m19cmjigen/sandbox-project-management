package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

// organizationRepository はOrganizationRepositoryのPostgreSQL実装
type organizationRepository struct {
	db *sqlx.DB
}

// NewOrganizationRepository は新しいOrganizationRepositoryを作成
func NewOrganizationRepository(db *sqlx.DB) repository.OrganizationRepository {
	return &organizationRepository{db: db}
}

// FindAll は全ての組織を取得
func (r *organizationRepository) FindAll(ctx context.Context) ([]domain.Organization, error) {
	var orgs []domain.Organization
	query := `
		SELECT id, name, parent_id, path, level, created_at, updated_at
		FROM organizations
		ORDER BY path
	`
	err := r.db.SelectContext(ctx, &orgs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all organizations: %w", err)
	}
	return orgs, nil
}

// FindByID はIDで組織を取得
func (r *organizationRepository) FindByID(ctx context.Context, id int64) (*domain.Organization, error) {
	var org domain.Organization
	query := `
		SELECT id, name, parent_id, path, level, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &org, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find organization by id %d: %w", id, err)
	}
	return &org, nil
}

// FindByParentID は親組織IDで子組織を取得
func (r *organizationRepository) FindByParentID(ctx context.Context, parentID int64) ([]domain.Organization, error) {
	var orgs []domain.Organization
	query := `
		SELECT id, name, parent_id, path, level, created_at, updated_at
		FROM organizations
		WHERE parent_id = $1
		ORDER BY name
	`
	err := r.db.SelectContext(ctx, &orgs, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to find organizations by parent_id %d: %w", parentID, err)
	}
	return orgs, nil
}

// FindRoots はルート組織（parent_id が NULL）を取得
func (r *organizationRepository) FindRoots(ctx context.Context) ([]domain.Organization, error) {
	var orgs []domain.Organization
	query := `
		SELECT id, name, parent_id, path, level, created_at, updated_at
		FROM organizations
		WHERE parent_id IS NULL
		ORDER BY name
	`
	err := r.db.SelectContext(ctx, &orgs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find root organizations: %w", err)
	}
	return orgs, nil
}

// FindByPath はパスで組織を検索（階層検索用）
func (r *organizationRepository) FindByPath(ctx context.Context, pathPrefix string) ([]domain.Organization, error) {
	var orgs []domain.Organization
	query := `
		SELECT id, name, parent_id, path, level, created_at, updated_at
		FROM organizations
		WHERE path LIKE $1
		ORDER BY path
	`
	err := r.db.SelectContext(ctx, &orgs, query, pathPrefix+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to find organizations by path %s: %w", pathPrefix, err)
	}
	return orgs, nil
}

// Create は新しい組織を作成
func (r *organizationRepository) Create(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (name, parent_id, path, level)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		org.Name,
		org.ParentID,
		org.Path,
		org.Level,
	).Scan(&org.ID, &org.CreatedAt, &org.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create organization: %w", err)
	}
	return nil
}

// Update は組織を更新
func (r *organizationRepository) Update(ctx context.Context, org *domain.Organization) error {
	query := `
		UPDATE organizations
		SET name = $1, parent_id = $2, path = $3, level = $4
		WHERE id = $5
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		org.Name,
		org.ParentID,
		org.Path,
		org.Level,
		org.ID,
	).Scan(&org.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("organization not found: id=%d", org.ID)
		}
		return fmt.Errorf("failed to update organization: %w", err)
	}
	return nil
}

// Delete は組織を削除
func (r *organizationRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM organizations WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("organization not found: id=%d", id)
	}

	return nil
}

// ExistsByID は組織が存在するかチェック
func (r *organizationRepository) ExistsByID(ctx context.Context, id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM organizations WHERE id = $1)`
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check organization exists: %w", err)
	}
	return exists, nil
}

// HasChildren は子組織が存在するかチェック
func (r *organizationRepository) HasChildren(ctx context.Context, id int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM organizations WHERE parent_id = $1)`
	err := r.db.GetContext(ctx, &exists, query, id)
	if err != nil {
		return false, fmt.Errorf("failed to check has children: %w", err)
	}
	return exists, nil
}
