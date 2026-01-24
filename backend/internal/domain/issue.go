package domain

import (
	"database/sql"
	"time"
)

// DelayStatus は遅延ステータスの型
type DelayStatus string

const (
	DelayStatusRed    DelayStatus = "RED"
	DelayStatusYellow DelayStatus = "YELLOW"
	DelayStatusGreen  DelayStatus = "GREEN"
)

// StatusCategory はステータスカテゴリの型
type StatusCategory string

const (
	StatusCategoryToDo       StatusCategory = "To Do"
	StatusCategoryInProgress StatusCategory = "In Progress"
	StatusCategoryDone       StatusCategory = "Done"
)

// Issue はJiraチケット情報を表すエンティティ
type Issue struct {
	ID                int64          `db:"id" json:"id"`
	JiraIssueID       string         `db:"jira_issue_id" json:"jira_issue_id"`
	JiraIssueKey      string         `db:"jira_issue_key" json:"jira_issue_key"`
	ProjectID         int64          `db:"project_id" json:"project_id"`
	Summary           string         `db:"summary" json:"summary"`
	Status            string         `db:"status" json:"status"`
	StatusCategory    StatusCategory `db:"status_category" json:"status_category"`
	DueDate           sql.NullTime   `db:"due_date" json:"due_date"`
	AssigneeName      *string        `db:"assignee_name" json:"assignee_name"`
	AssigneeAccountID *string        `db:"assignee_account_id" json:"assignee_account_id"`
	DelayStatus       DelayStatus    `db:"delay_status" json:"delay_status"`
	Priority          *string        `db:"priority" json:"priority"`
	IssueType         *string        `db:"issue_type" json:"issue_type"`
	LastUpdatedAt     time.Time      `db:"last_updated_at" json:"last_updated_at"`
	CreatedAt         time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at" json:"updated_at"`
}

// IssueFilter はIssue検索のフィルタ条件
type IssueFilter struct {
	ProjectID    *int64
	DelayStatus  *DelayStatus
	Status       *string
	AssigneeID   *string
	DueDateFrom  *time.Time
	DueDateTo    *time.Time
	Limit        int
	Offset       int
}
