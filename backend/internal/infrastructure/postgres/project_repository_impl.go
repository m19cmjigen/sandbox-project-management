package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

// projectRepository はProjectRepositoryのPostgreSQL実装
type projectRepository struct {
	db *sqlx.DB
}

// NewProjectRepository は新しいProjectRepositoryを作成
func NewProjectRepository(db *sqlx.DB) repository.ProjectRepository {
	return &projectRepository{db: db}
}

// FindAll は全てのプロジェクトを取得
func (r *projectRepository) FindAll(ctx context.Context) ([]domain.Project, error) {
	var projects []domain.Project
	query := `
		SELECT id, jira_project_id, key, name, lead_account_id, lead_email,
		       organization_id, created_at, updated_at
		FROM projects
		ORDER BY name
	`
	err := r.db.SelectContext(ctx, &projects, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all projects: %w", err)
	}
	return projects, nil
}

// FindByID はIDでプロジェクトを取得
func (r *projectRepository) FindByID(ctx context.Context, id int64) (*domain.Project, error) {
	var project domain.Project
	query := `
		SELECT id, jira_project_id, key, name, lead_account_id, lead_email,
		       organization_id, created_at, updated_at
		FROM projects
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &project, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project by id %d: %w", id, err)
	}
	return &project, nil
}

// GetByKey はプロジェクトキー(例: PROJ)で取得
func (r *projectRepository) GetByKey(ctx context.Context, key string) (*domain.Project, error) {
	var project domain.Project
	query := `
		SELECT id, jira_project_id, key, name, lead_account_id, lead_email,
		       organization_id, created_at, updated_at
		FROM projects
		WHERE key = $1
	`
	err := r.db.GetContext(ctx, &project, query, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project by key %s: %w", key, err)
	}
	return &project, nil
}

// FindByJiraProjectID はJiraプロジェクトIDで取得
func (r *projectRepository) FindByJiraProjectID(ctx context.Context, jiraProjectID string) (*domain.Project, error) {
	var project domain.Project
	query := `
		SELECT id, jira_project_id, key, name, lead_account_id, lead_email,
		       organization_id, created_at, updated_at
		FROM projects
		WHERE jira_project_id = $1
	`
	err := r.db.GetContext(ctx, &project, query, jiraProjectID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project by jira_project_id %s: %w", jiraProjectID, err)
	}
	return &project, nil
}

// FindByOrganizationID は組織IDでプロジェクトを取得
func (r *projectRepository) FindByOrganizationID(ctx context.Context, organizationID int64) ([]domain.Project, error) {
	var projects []domain.Project
	query := `
		SELECT id, jira_project_id, key, name, lead_account_id, lead_email,
		       organization_id, created_at, updated_at
		FROM projects
		WHERE organization_id = $1
		ORDER BY name
	`
	err := r.db.SelectContext(ctx, &projects, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find projects by organization_id %d: %w", organizationID, err)
	}
	return projects, nil
}

// FindUnassigned は未分類プロジェクト（organization_id が NULL）を取得
func (r *projectRepository) FindUnassigned(ctx context.Context) ([]domain.Project, error) {
	var projects []domain.Project
	query := `
		SELECT id, jira_project_id, key, name, lead_account_id, lead_email,
		       organization_id, created_at, updated_at
		FROM projects
		WHERE organization_id IS NULL
		ORDER BY name
	`
	err := r.db.SelectContext(ctx, &projects, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find unassigned projects: %w", err)
	}
	return projects, nil
}

// FindWithStats は統計情報付きでプロジェクトを取得
func (r *projectRepository) FindWithStats(ctx context.Context, id int64) (*domain.ProjectWithStats, error) {
	var stats domain.ProjectWithStats
	query := `
		SELECT
			p.id, p.jira_project_id, p.key, p.name, p.lead_account_id,
			p.lead_email, p.organization_id, p.created_at, p.updated_at,
			COALESCE(s.total_issues, 0) as total_issues,
			COALESCE(s.red_issues, 0) as red_issues,
			COALESCE(s.yellow_issues, 0) as yellow_issues,
			COALESCE(s.green_issues, 0) as green_issues,
			COALESCE(s.open_issues, 0) as open_issues,
			COALESCE(s.done_issues, 0) as done_issues
		FROM projects p
		LEFT JOIN project_delay_summary s ON p.id = s.project_id
		WHERE p.id = $1
	`
	err := r.db.GetContext(ctx, &stats, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find project with stats by id %d: %w", id, err)
	}
	return &stats, nil
}

// FindAllWithStats は全プロジェクトを統計情報付きで取得
func (r *projectRepository) FindAllWithStats(ctx context.Context) ([]domain.ProjectWithStats, error) {
	var stats []domain.ProjectWithStats
	query := `
		SELECT
			p.id, p.jira_project_id, p.key, p.name, p.lead_account_id,
			p.lead_email, p.organization_id, p.created_at, p.updated_at,
			COALESCE(s.total_issues, 0) as total_issues,
			COALESCE(s.red_issues, 0) as red_issues,
			COALESCE(s.yellow_issues, 0) as yellow_issues,
			COALESCE(s.green_issues, 0) as green_issues,
			COALESCE(s.open_issues, 0) as open_issues,
			COALESCE(s.done_issues, 0) as done_issues
		FROM projects p
		LEFT JOIN project_delay_summary s ON p.id = s.project_id
		ORDER BY p.name
	`
	err := r.db.SelectContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all projects with stats: %w", err)
	}
	return stats, nil
}

// Create は新しいプロジェクトを作成
func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
	query := `
		INSERT INTO projects (jira_project_id, key, name, lead_account_id, lead_email, organization_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		project.JiraProjectID,
		project.Key,
		project.Name,
		project.LeadAccountID,
		project.LeadEmail,
		project.OrganizationID,
	).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

// Update はプロジェクトを更新
func (r *projectRepository) Update(ctx context.Context, project *domain.Project) error {
	query := `
		UPDATE projects
		SET jira_project_id = $1, key = $2, name = $3,
		    lead_account_id = $4, lead_email = $5, organization_id = $6
		WHERE id = $7
		RETURNING updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		project.JiraProjectID,
		project.Key,
		project.Name,
		project.LeadAccountID,
		project.LeadEmail,
		project.OrganizationID,
		project.ID,
	).Scan(&project.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("project not found: id=%d", project.ID)
		}
		return fmt.Errorf("failed to update project: %w", err)
	}
	return nil
}

// Delete はプロジェクトを削除
func (r *projectRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM projects WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found: id=%d", id)
	}

	return nil
}

// AssignToOrganization はプロジェクトを組織に紐付け
func (r *projectRepository) AssignToOrganization(ctx context.Context, projectID int64, organizationID *int64) error {
	query := `
		UPDATE projects
		SET organization_id = $1
		WHERE id = $2
	`
	result, err := r.db.ExecContext(ctx, query, organizationID, projectID)
	if err != nil {
		return fmt.Errorf("failed to assign project to organization: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found: id=%d", projectID)
	}

	return nil
}

// ExistsByJiraProjectID はJiraプロジェクトIDで存在チェック
func (r *projectRepository) ExistsByJiraProjectID(ctx context.Context, jiraProjectID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM projects WHERE jira_project_id = $1)`
	err := r.db.GetContext(ctx, &exists, query, jiraProjectID)
	if err != nil {
		return false, fmt.Errorf("failed to check project exists: %w", err)
	}
	return exists, nil
}
