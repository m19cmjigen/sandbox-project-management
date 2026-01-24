package domain

import "time"

// Project はJiraプロジェクト情報を表すエンティティ
type Project struct {
	ID              int64     `db:"id" json:"id"`
	JiraProjectID   string    `db:"jira_project_id" json:"jira_project_id"`
	Key             string    `db:"key" json:"key"`
	Name            string    `db:"name" json:"name"`
	LeadAccountID   *string   `db:"lead_account_id" json:"lead_account_id"`
	LeadEmail       *string   `db:"lead_email" json:"lead_email"`
	OrganizationID  *int64    `db:"organization_id" json:"organization_id"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

// ProjectWithStats はプロジェクトと統計情報
type ProjectWithStats struct {
	Project
	TotalIssues   int `json:"total_issues"`
	RedIssues     int `json:"red_issues"`
	YellowIssues  int `json:"yellow_issues"`
	GreenIssues   int `json:"green_issues"`
	OpenIssues    int `json:"open_issues"`
	DoneIssues    int `json:"done_issues"`
}
