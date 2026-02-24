package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/normalizer"
)

// Repository defines the DB operations required by the sync process.
type Repository interface {
	// UpsertProjects inserts or updates projects and returns the number of rows affected.
	UpsertProjects(ctx context.Context, projects []normalizer.DBProject) (int, error)
	// UpsertIssues inserts or updates issues. projectIDMap maps jira_project_id → DB id.
	// Returns the number of rows affected.
	UpsertIssues(ctx context.Context, issues []normalizer.DBIssue, projectIDMap map[string]int64) (int, error)
	// GetProjectIDMap returns a map of jira_project_id → DB id for all known projects.
	GetProjectIDMap(ctx context.Context) (map[string]int64, error)
	// StartSyncLog creates a sync_log record in RUNNING state and returns its ID.
	StartSyncLog(ctx context.Context, syncType string) (int64, error)
	// FinishSyncLog updates the sync_log record with the final status.
	FinishSyncLog(ctx context.Context, id int64, status string, projectsSynced, issuesSynced int, errMsg string) error
}

// sqlxRepository is the PostgreSQL implementation of Repository.
type sqlxRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository backed by the given sqlx.DB.
func NewRepository(db *sqlx.DB) Repository {
	return &sqlxRepository{db: db}
}

func (r *sqlxRepository) UpsertProjects(ctx context.Context, projects []normalizer.DBProject) (int, error) {
	if len(projects) == 0 {
		return 0, nil
	}

	const q = `
		INSERT INTO projects (jira_project_id, key, name, lead_account_id, lead_email)
		VALUES (:jira_project_id, :key, :name, :lead_account_id, :lead_email)
		ON CONFLICT (jira_project_id) DO UPDATE SET
			key            = EXCLUDED.key,
			name           = EXCLUDED.name,
			lead_account_id = EXCLUDED.lead_account_id,
			lead_email     = EXCLUDED.lead_email,
			updated_at     = CURRENT_TIMESTAMP`

	type row struct {
		JiraProjectID string `db:"jira_project_id"`
		Key           string `db:"key"`
		Name          string `db:"name"`
		LeadAccountID string `db:"lead_account_id"`
		LeadEmail     string `db:"lead_email"`
	}

	rows := make([]row, len(projects))
	for i, p := range projects {
		rows[i] = row{
			JiraProjectID: p.JiraProjectID,
			Key:           p.Key,
			Name:          p.Name,
			LeadAccountID: p.LeadAccountID,
			LeadEmail:     p.LeadEmail,
		}
	}

	result, err := r.db.NamedExecContext(ctx, q, rows)
	if err != nil {
		return 0, fmt.Errorf("upsert projects: %w", err)
	}
	n, _ := result.RowsAffected()
	return int(n), nil
}

func (r *sqlxRepository) UpsertIssues(ctx context.Context, issues []normalizer.DBIssue, projectIDMap map[string]int64) (int, error) {
	if len(issues) == 0 {
		return 0, nil
	}

	const q = `
		INSERT INTO issues (
			jira_issue_id, jira_issue_key, project_id, summary,
			status, status_category, due_date,
			assignee_name, assignee_account_id,
			delay_status, priority, issue_type, last_updated_at
		) VALUES (
			:jira_issue_id, :jira_issue_key, :project_id, :summary,
			:status, :status_category, :due_date,
			:assignee_name, :assignee_account_id,
			:delay_status, :priority, :issue_type, :last_updated_at
		)
		ON CONFLICT (jira_issue_id) DO UPDATE SET
			jira_issue_key      = EXCLUDED.jira_issue_key,
			summary             = EXCLUDED.summary,
			status              = EXCLUDED.status,
			status_category     = EXCLUDED.status_category,
			due_date            = EXCLUDED.due_date,
			assignee_name       = EXCLUDED.assignee_name,
			assignee_account_id = EXCLUDED.assignee_account_id,
			delay_status        = EXCLUDED.delay_status,
			priority            = EXCLUDED.priority,
			issue_type          = EXCLUDED.issue_type,
			last_updated_at     = EXCLUDED.last_updated_at,
			updated_at          = CURRENT_TIMESTAMP`

	type row struct {
		JiraIssueID       string  `db:"jira_issue_id"`
		JiraIssueKey      string  `db:"jira_issue_key"`
		ProjectID         int64   `db:"project_id"`
		Summary           string  `db:"summary"`
		Status            string  `db:"status"`
		StatusCategory    string  `db:"status_category"`
		DueDate           *string `db:"due_date"`
		AssigneeName      string  `db:"assignee_name"`
		AssigneeAccountID string  `db:"assignee_account_id"`
		DelayStatus       string  `db:"delay_status"`
		Priority          string  `db:"priority"`
		IssueType         string  `db:"issue_type"`
		LastUpdatedAt     time.Time `db:"last_updated_at"`
	}

	var rows []row
	for _, issue := range issues {
		projectID, ok := projectIDMap[issue.JiraProjectID]
		if !ok {
			// 紐付け先プロジェクトが未知の場合はスキップ
			continue
		}
		rows = append(rows, row{
			JiraIssueID:       issue.JiraIssueID,
			JiraIssueKey:      issue.JiraIssueKey,
			ProjectID:         projectID,
			Summary:           issue.Summary,
			Status:            issue.Status,
			StatusCategory:    issue.StatusCategory,
			DueDate:           issue.DueDate,
			AssigneeName:      issue.AssigneeName,
			AssigneeAccountID: issue.AssigneeAccountID,
			DelayStatus:       issue.DelayStatus,
			Priority:          issue.Priority,
			IssueType:         issue.IssueType,
			LastUpdatedAt:     issue.LastUpdatedAt,
		})
	}

	if len(rows) == 0 {
		return 0, nil
	}

	result, err := r.db.NamedExecContext(ctx, q, rows)
	if err != nil {
		return 0, fmt.Errorf("upsert issues: %w", err)
	}
	n, _ := result.RowsAffected()
	return int(n), nil
}

func (r *sqlxRepository) GetProjectIDMap(ctx context.Context) (map[string]int64, error) {
	rows, err := r.db.QueryxContext(ctx, `SELECT id, jira_project_id FROM projects`)
	if err != nil {
		return nil, fmt.Errorf("get project id map: %w", err)
	}
	defer rows.Close()

	m := make(map[string]int64)
	for rows.Next() {
		var id int64
		var jiraID string
		if err := rows.Scan(&id, &jiraID); err != nil {
			return nil, err
		}
		m[jiraID] = id
	}
	return m, rows.Err()
}

func (r *sqlxRepository) StartSyncLog(ctx context.Context, syncType string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO sync_logs (sync_type, status) VALUES ($1, 'RUNNING') RETURNING id`,
		syncType,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("start sync log: %w", err)
	}
	return id, nil
}

func (r *sqlxRepository) FinishSyncLog(ctx context.Context, id int64, status string, projectsSynced, issuesSynced int, errMsg string) error {
	var errMsgPtr *string
	if errMsg != "" {
		errMsgPtr = &errMsg
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE sync_logs SET
			status           = $1,
			completed_at     = CURRENT_TIMESTAMP,
			projects_synced  = $2,
			issues_synced    = $3,
			error_message    = $4,
			duration_seconds = EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - executed_at))::INTEGER
		WHERE id = $5`,
		status, projectsSynced, issuesSynced, errMsgPtr, id,
	)
	if err != nil {
		return fmt.Errorf("finish sync log: %w", err)
	}
	return nil
}
