// Package normalizer converts Jira API responses into normalized DB records.
// It implements the same status normalization and delay calculation logic
// as the PostgreSQL triggers defined in the initial schema migration.
package normalizer

import (
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/jiraclient"
)

// jst is the Japan Standard Time location used for date boundary calculations.
// Jira due dates are plain dates (no time component), so we compare against
// the current date in JST to avoid timezone-related off-by-one errors.
var jst = mustLoadLocation("Asia/Tokyo")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		// UTC にフォールバック（テスト環境でtzdata がない場合）
		return time.UTC
	}
	return loc
}

// NormalizeStatusCategory maps a Jira statusCategory key to one of the three
// canonical values used in the DB schema: "To Do", "In Progress", or "Done".
func NormalizeStatusCategory(key string) string {
	switch key {
	case "done":
		return "Done"
	case "indeterminate":
		return "In Progress"
	default:
		// "new" およびその他の未知の値はすべて "To Do" に分類する
		return "To Do"
	}
}

// CalcDelayStatus computes the delay status of an issue using the same logic as
// the calculate_delay_status PostgreSQL trigger function.
//
//   - RED    : not Done AND due_date is in the past
//   - YELLOW : not Done AND (due_date is nil OR due_date is within 3 days from now)
//   - GREEN  : Done OR (due_date exists and is more than 3 days away)
//
// now should be the current time in the desired timezone (typically JST).
func CalcDelayStatus(statusCategory string, dueDate *string, now time.Time) string {
	if statusCategory == "Done" {
		return "GREEN"
	}

	if dueDate == nil || *dueDate == "" {
		// 期限未設定は注意扱い
		return "YELLOW"
	}

	due, err := time.ParseInLocation("2006-01-02", *dueDate, now.Location())
	if err != nil {
		// パースできない日付は期限未設定と同様に扱う
		return "YELLOW"
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	threeDay := today.AddDate(0, 0, 3)

	switch {
	case due.Before(today):
		return "RED"
	case !due.After(threeDay): // today <= due <= today+3days
		return "YELLOW"
	default:
		return "GREEN"
	}
}

// ConvertProject converts a jiraclient.Project to a DBProject.
func ConvertProject(p jiraclient.Project) DBProject {
	dp := DBProject{
		JiraProjectID: p.ID,
		Key:           p.Key,
		Name:          p.Name,
	}
	if p.Lead != nil {
		dp.LeadAccountID = p.Lead.AccountID
		dp.LeadEmail = p.Lead.EmailAddress
	}
	return dp
}

// ConvertIssue converts a jiraclient.Issue to a DBIssue.
// now is used for delay status calculation and should typically be time.Now().In(jst).
func ConvertIssue(issue jiraclient.Issue, now time.Time) DBIssue {
	statusCategory := NormalizeStatusCategory(issue.Fields.Status.StatusCategory.Key)

	var dueDate *string
	if issue.Fields.DueDate != "" {
		d := issue.Fields.DueDate
		dueDate = &d
	}

	di := DBIssue{
		JiraIssueID:    issue.ID,
		JiraIssueKey:   issue.Key,
		JiraProjectID:  issue.Fields.Project.ID,
		Summary:        issue.Fields.Summary,
		Status:         issue.Fields.Status.Name,
		StatusCategory: statusCategory,
		DueDate:        dueDate,
		IssueType:      issue.Fields.IssueType.Name,
		DelayStatus:    CalcDelayStatus(statusCategory, dueDate, now),
	}

	if issue.Fields.Assignee != nil {
		di.AssigneeName = issue.Fields.Assignee.DisplayName
		di.AssigneeAccountID = issue.Fields.Assignee.AccountID
	}

	if issue.Fields.Priority != nil {
		di.Priority = issue.Fields.Priority.Name
	}

	if issue.Fields.Updated != "" {
		if t, err := time.Parse(time.RFC3339, issue.Fields.Updated); err == nil {
			di.LastUpdatedAt = t
		}
	}

	return di
}

// Now returns the current time in JST. Use this as the now argument when calling
// ConvertIssue or CalcDelayStatus in production code.
func Now() time.Time {
	return time.Now().In(jst)
}
