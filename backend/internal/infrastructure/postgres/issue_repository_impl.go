package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/repository"
)

// issueRepository はIssueRepositoryのPostgreSQL実装
type issueRepository struct {
	db *sqlx.DB
}

// NewIssueRepository は新しいIssueRepositoryを作成
func NewIssueRepository(db *sqlx.DB) repository.IssueRepository {
	return &issueRepository{db: db}
}

// FindAll は全てのチケットを取得
func (r *issueRepository) FindAll(ctx context.Context) ([]domain.Issue, error) {
	var issues []domain.Issue
	query := `
		SELECT id, jira_issue_id, jira_issue_key, project_id, summary, status,
		       status_category, due_date, assignee_name, assignee_account_id,
		       delay_status, priority, issue_type, last_updated_at, created_at, updated_at
		FROM issues
		ORDER BY last_updated_at DESC
	`
	err := r.db.SelectContext(ctx, &issues, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find all issues: %w", err)
	}
	return issues, nil
}

// FindByID はIDでチケットを取得
func (r *issueRepository) FindByID(ctx context.Context, id int64) (*domain.Issue, error) {
	var issue domain.Issue
	query := `
		SELECT id, jira_issue_id, jira_issue_key, project_id, summary, status,
		       status_category, due_date, assignee_name, assignee_account_id,
		       delay_status, priority, issue_type, last_updated_at, created_at, updated_at
		FROM issues
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &issue, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find issue by id %d: %w", id, err)
	}
	return &issue, nil
}

// FindByJiraIssueID はJiraチケットIDで取得
func (r *issueRepository) FindByJiraIssueID(ctx context.Context, jiraIssueID string) (*domain.Issue, error) {
	var issue domain.Issue
	query := `
		SELECT id, jira_issue_id, jira_issue_key, project_id, summary, status,
		       status_category, due_date, assignee_name, assignee_account_id,
		       delay_status, priority, issue_type, last_updated_at, created_at, updated_at
		FROM issues
		WHERE jira_issue_id = $1
	`
	err := r.db.GetContext(ctx, &issue, query, jiraIssueID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find issue by jira_issue_id %s: %w", jiraIssueID, err)
	}
	return &issue, nil
}

// FindByProjectID はプロジェクトIDでチケットを取得
func (r *issueRepository) FindByProjectID(ctx context.Context, projectID int64) ([]domain.Issue, error) {
	var issues []domain.Issue
	query := `
		SELECT id, jira_issue_id, jira_issue_key, project_id, summary, status,
		       status_category, due_date, assignee_name, assignee_account_id,
		       delay_status, priority, issue_type, last_updated_at, created_at, updated_at
		FROM issues
		WHERE project_id = $1
		ORDER BY last_updated_at DESC
	`
	err := r.db.SelectContext(ctx, &issues, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to find issues by project_id %d: %w", projectID, err)
	}
	return issues, nil
}

// FindByFilter はフィルタ条件でチケットを取得
func (r *issueRepository) FindByFilter(ctx context.Context, filter domain.IssueFilter) ([]domain.Issue, error) {
	var issues []domain.Issue
	var conditions []string
	var args []interface{}
	argIndex := 1

	// ベースクエリ
	query := `
		SELECT id, jira_issue_id, jira_issue_key, project_id, summary, status,
		       status_category, due_date, assignee_name, assignee_account_id,
		       delay_status, priority, issue_type, last_updated_at, created_at, updated_at
		FROM issues
	`

	// フィルタ条件を構築
	if filter.ProjectID != nil {
		conditions = append(conditions, fmt.Sprintf("project_id = $%d", argIndex))
		args = append(args, *filter.ProjectID)
		argIndex++
	}

	if filter.DelayStatus != nil {
		conditions = append(conditions, fmt.Sprintf("delay_status = $%d", argIndex))
		args = append(args, *filter.DelayStatus)
		argIndex++
	}

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.AssigneeID != nil {
		conditions = append(conditions, fmt.Sprintf("assignee_account_id = $%d", argIndex))
		args = append(args, *filter.AssigneeID)
		argIndex++
	}

	if filter.DueDateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("due_date >= $%d", argIndex))
		args = append(args, *filter.DueDateFrom)
		argIndex++
	}

	if filter.DueDateTo != nil {
		conditions = append(conditions, fmt.Sprintf("due_date <= $%d", argIndex))
		args = append(args, *filter.DueDateTo)
		argIndex++
	}

	// WHERE句を追加
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// ORDER BY
	query += " ORDER BY last_updated_at DESC"

	// LIMIT/OFFSET
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filter.Limit)
		argIndex++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	err := r.db.SelectContext(ctx, &issues, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find issues by filter: %w", err)
	}
	return issues, nil
}

// FindByDelayStatus は遅延ステータスでチケットを取得
func (r *issueRepository) FindByDelayStatus(ctx context.Context, status domain.DelayStatus) ([]domain.Issue, error) {
	var issues []domain.Issue
	query := `
		SELECT id, jira_issue_id, jira_issue_key, project_id, summary, status,
		       status_category, due_date, assignee_name, assignee_account_id,
		       delay_status, priority, issue_type, last_updated_at, created_at, updated_at
		FROM issues
		WHERE delay_status = $1
		ORDER BY due_date NULLS LAST, last_updated_at DESC
	`
	err := r.db.SelectContext(ctx, &issues, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to find issues by delay_status %s: %w", status, err)
	}
	return issues, nil
}

// CountByProjectID はプロジェクトIDでチケット数を取得
func (r *issueRepository) CountByProjectID(ctx context.Context, projectID int64) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM issues WHERE project_id = $1`
	err := r.db.GetContext(ctx, &count, query, projectID)
	if err != nil {
		return 0, fmt.Errorf("failed to count issues by project_id %d: %w", projectID, err)
	}
	return count, nil
}

// CountByDelayStatus は遅延ステータスでチケット数を取得
func (r *issueRepository) CountByDelayStatus(ctx context.Context, projectID int64, status domain.DelayStatus) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM issues WHERE project_id = $1 AND delay_status = $2`
	err := r.db.GetContext(ctx, &count, query, projectID, status)
	if err != nil {
		return 0, fmt.Errorf("failed to count issues by delay_status: %w", err)
	}
	return count, nil
}

// Create は新しいチケットを作成
func (r *issueRepository) Create(ctx context.Context, issue *domain.Issue) error {
	query := `
		INSERT INTO issues (
			jira_issue_id, jira_issue_key, project_id, summary, status,
			status_category, due_date, assignee_name, assignee_account_id,
			priority, issue_type, last_updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, delay_status, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		issue.JiraIssueID,
		issue.JiraIssueKey,
		issue.ProjectID,
		issue.Summary,
		issue.Status,
		issue.StatusCategory,
		issue.DueDate,
		issue.AssigneeName,
		issue.AssigneeAccountID,
		issue.Priority,
		issue.IssueType,
		issue.LastUpdatedAt,
	).Scan(&issue.ID, &issue.DelayStatus, &issue.CreatedAt, &issue.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}
	return nil
}

// Update はチケットを更新
func (r *issueRepository) Update(ctx context.Context, issue *domain.Issue) error {
	query := `
		UPDATE issues
		SET jira_issue_key = $1, project_id = $2, summary = $3, status = $4,
		    status_category = $5, due_date = $6, assignee_name = $7, assignee_account_id = $8,
		    priority = $9, issue_type = $10, last_updated_at = $11
		WHERE jira_issue_id = $12
		RETURNING id, delay_status, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		issue.JiraIssueKey,
		issue.ProjectID,
		issue.Summary,
		issue.Status,
		issue.StatusCategory,
		issue.DueDate,
		issue.AssigneeName,
		issue.AssigneeAccountID,
		issue.Priority,
		issue.IssueType,
		issue.LastUpdatedAt,
		issue.JiraIssueID,
	).Scan(&issue.ID, &issue.DelayStatus, &issue.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("issue not found: jira_issue_id=%s", issue.JiraIssueID)
		}
		return fmt.Errorf("failed to update issue: %w", err)
	}
	return nil
}

// Delete はチケットを削除
func (r *issueRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM issues WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("issue not found: id=%d", id)
	}

	return nil
}

// BulkUpsert は複数のチケットを一括登録/更新
func (r *issueRepository) BulkUpsert(ctx context.Context, issues []domain.Issue) error {
	if len(issues) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO issues (
			jira_issue_id, jira_issue_key, project_id, summary, status,
			status_category, due_date, assignee_name, assignee_account_id,
			priority, issue_type, last_updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (jira_issue_id) DO UPDATE SET
			jira_issue_key = EXCLUDED.jira_issue_key,
			project_id = EXCLUDED.project_id,
			summary = EXCLUDED.summary,
			status = EXCLUDED.status,
			status_category = EXCLUDED.status_category,
			due_date = EXCLUDED.due_date,
			assignee_name = EXCLUDED.assignee_name,
			assignee_account_id = EXCLUDED.assignee_account_id,
			priority = EXCLUDED.priority,
			issue_type = EXCLUDED.issue_type,
			last_updated_at = EXCLUDED.last_updated_at
	`

	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, issue := range issues {
		_, err := stmt.ExecContext(
			ctx,
			issue.JiraIssueID,
			issue.JiraIssueKey,
			issue.ProjectID,
			issue.Summary,
			issue.Status,
			issue.StatusCategory,
			issue.DueDate,
			issue.AssigneeName,
			issue.AssigneeAccountID,
			issue.Priority,
			issue.IssueType,
			issue.LastUpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to upsert issue %s: %w", issue.JiraIssueID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ExistsByJiraIssueID はJiraチケットIDで存在チェック
func (r *issueRepository) ExistsByJiraIssueID(ctx context.Context, jiraIssueID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM issues WHERE jira_issue_id = $1)`
	err := r.db.GetContext(ctx, &exists, query, jiraIssueID)
	if err != nil {
		return false, fmt.Errorf("failed to check issue exists: %w", err)
	}
	return exists, nil
}
