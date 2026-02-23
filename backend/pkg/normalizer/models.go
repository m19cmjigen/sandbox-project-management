package normalizer

import "time"

// DBProject is the normalized project record ready to be upserted into the DB.
type DBProject struct {
	JiraProjectID string
	Key           string
	Name          string
	LeadAccountID string // empty string if no lead
	LeadEmail     string // empty string if no lead
}

// DBIssue is the normalized issue record ready to be upserted into the DB.
type DBIssue struct {
	JiraIssueID       string
	JiraIssueKey      string
	JiraProjectID     string // used to resolve project_id FK
	Summary           string
	Status            string
	StatusCategory    string    // "To Do" | "In Progress" | "Done"
	DueDate           *string   // nil when not set; "YYYY-MM-DD" when set
	AssigneeName      string    // empty string if unassigned
	AssigneeAccountID string    // empty string if unassigned
	DelayStatus       string    // "RED" | "YELLOW" | "GREEN"
	Priority          string    // empty string if not set
	IssueType         string
	LastUpdatedAt     time.Time
}
