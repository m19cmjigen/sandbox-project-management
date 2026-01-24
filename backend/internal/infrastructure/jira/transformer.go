package jira

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/domain"
)

// TransformProject converts a Jira project to our domain project
func TransformProject(jiraProject JiraProject, organizationID int64) *domain.Project {
	orgID := &organizationID
	return &domain.Project{
		Key:            jiraProject.Key,
		Name:           jiraProject.Name,
		OrganizationID: orgID,
		JiraProjectID:  jiraProject.ID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// TransformIssue converts a Jira issue to our domain issue
func TransformIssue(jiraIssue JiraIssue, projectID int64) (*domain.Issue, error) {
	// Parse timestamps
	createdAt, err := parseJiraTimestamp(jiraIssue.Fields.Created)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created timestamp: %w", err)
	}

	updatedAt, err := parseJiraTimestamp(jiraIssue.Fields.Updated)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated timestamp: %w", err)
	}

	// Parse due date (optional)
	var dueDate *time.Time
	if jiraIssue.Fields.DueDate != "" {
		parsed, err := parseDueDate(jiraIssue.Fields.DueDate)
		if err == nil {
			dueDate = &parsed
		}
	}

	// Get assignee name
	var assigneeName *string
	if jiraIssue.Fields.Assignee != nil {
		assigneeName = &jiraIssue.Fields.Assignee.DisplayName
	}

	// Convert time.Time to sql.NullTime for DueDate
	var dueDateNull sql.NullTime
	if dueDate != nil {
		dueDateNull = sql.NullTime{
			Time:  *dueDate,
			Valid: true,
		}
	}

	// Priority as pointer
	priority := jiraIssue.Fields.Priority.Name

	issue := &domain.Issue{
		ProjectID:      projectID,
		JiraIssueKey:   jiraIssue.Key,
		JiraIssueID:    jiraIssue.ID,
		Summary:        jiraIssue.Fields.Summary,
		Status:         jiraIssue.Fields.Status.Name,
		DueDate:        dueDateNull,
		AssigneeName:   assigneeName,
		Priority:       &priority,
		LastUpdatedAt:  updatedAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		// DelayStatus will be calculated by database trigger
		DelayStatus: "GREEN",
	}

	return issue, nil
}

// parseJiraTimestamp parses Jira timestamp format
func parseJiraTimestamp(timestamp string) (time.Time, error) {
	// Jira uses ISO 8601 format: 2023-01-15T10:30:00.000+0900
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05.999Z0700",
		"2006-01-02T15:04:05.999-0700",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, timestamp); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse timestamp: %s", timestamp)
}

// parseDueDate parses Jira due date format (YYYY-MM-DD)
func parseDueDate(dueDate string) (time.Time, error) {
	return time.Parse("2006-01-02", dueDate)
}

// IsDone checks if a Jira status indicates the issue is done
func IsDone(statusName string) bool {
	doneStatuses := map[string]bool{
		"Done":     true,
		"Closed":   true,
		"完了":       true,
		"クローズ":     true,
		"Resolved": true,
		"解決済み":     true,
	}
	return doneStatuses[statusName]
}
