package normalizer

import (
	"testing"
	"time"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/jiraclient"
)

// testNow は固定の "今日" として 2026-02-24 00:00:00 JST を使用する。
var testNow = time.Date(2026, 2, 24, 0, 0, 0, 0, time.FixedZone("JST", 9*60*60))

func strPtr(s string) *string { return &s }

// ----------------------------------------------------------------
// NormalizeStatusCategory
// ----------------------------------------------------------------

func TestNormalizeStatusCategory(t *testing.T) {
	cases := []struct {
		key      string
		expected string
	}{
		{"new", "To Do"},
		{"indeterminate", "In Progress"},
		{"done", "Done"},
		{"", "To Do"},          // 未知 → フォールバック
		{"unknown", "To Do"},   // 未知 → フォールバック
	}
	for _, tc := range cases {
		got := NormalizeStatusCategory(tc.key)
		if got != tc.expected {
			t.Errorf("NormalizeStatusCategory(%q) = %q, want %q", tc.key, got, tc.expected)
		}
	}
}

// ----------------------------------------------------------------
// CalcDelayStatus — 境界値テスト
// testNow = 2026-02-24
// ----------------------------------------------------------------

func TestCalcDelayStatus_Done(t *testing.T) {
	// 完了済みはすべて GREEN
	cases := []*string{nil, strPtr(""), strPtr("2026-02-20"), strPtr("2026-02-24"), strPtr("2030-01-01")}
	for _, d := range cases {
		got := CalcDelayStatus("Done", d, testNow)
		if got != "GREEN" {
			t.Errorf("Done with dueDate=%v: expected GREEN, got %s", d, got)
		}
	}
}

func TestCalcDelayStatus_NilDueDate(t *testing.T) {
	// 未完了かつ期限なし → YELLOW
	got := CalcDelayStatus("To Do", nil, testNow)
	if got != "YELLOW" {
		t.Errorf("expected YELLOW for nil dueDate, got %s", got)
	}
}

func TestCalcDelayStatus_EmptyDueDate(t *testing.T) {
	// 空文字も nil と同様
	got := CalcDelayStatus("In Progress", strPtr(""), testNow)
	if got != "YELLOW" {
		t.Errorf("expected YELLOW for empty dueDate, got %s", got)
	}
}

func TestCalcDelayStatus_Overdue(t *testing.T) {
	// testNow = 2026-02-24, 期限 2026-02-23 (昨日) → RED
	got := CalcDelayStatus("To Do", strPtr("2026-02-23"), testNow)
	if got != "RED" {
		t.Errorf("expected RED for yesterday's due date, got %s", got)
	}
}

func TestCalcDelayStatus_OverdueFarPast(t *testing.T) {
	// 期限が遠い過去でも RED
	got := CalcDelayStatus("In Progress", strPtr("2025-01-01"), testNow)
	if got != "RED" {
		t.Errorf("expected RED for far past due date, got %s", got)
	}
}

func TestCalcDelayStatus_DueDateIsToday(t *testing.T) {
	// 期限が今日 (0 days away) → YELLOW
	got := CalcDelayStatus("To Do", strPtr("2026-02-24"), testNow)
	if got != "YELLOW" {
		t.Errorf("expected YELLOW for today's due date, got %s", got)
	}
}

func TestCalcDelayStatus_DueDateIn1Day(t *testing.T) {
	// 期限が明日 → YELLOW
	got := CalcDelayStatus("To Do", strPtr("2026-02-25"), testNow)
	if got != "YELLOW" {
		t.Errorf("expected YELLOW for tomorrow's due date, got %s", got)
	}
}

func TestCalcDelayStatus_DueDateIn3Days(t *testing.T) {
	// 期限が今日+3日 (2026-02-27) → YELLOW (境界値: 3日以内の上限)
	got := CalcDelayStatus("To Do", strPtr("2026-02-27"), testNow)
	if got != "YELLOW" {
		t.Errorf("expected YELLOW for due date 3 days from now, got %s", got)
	}
}

func TestCalcDelayStatus_DueDateIn4Days(t *testing.T) {
	// 期限が今日+4日 (2026-02-28) → GREEN (境界値: 4日以上は GREEN)
	got := CalcDelayStatus("To Do", strPtr("2026-02-28"), testNow)
	if got != "GREEN" {
		t.Errorf("expected GREEN for due date 4 days from now, got %s", got)
	}
}

func TestCalcDelayStatus_DueDateFarFuture(t *testing.T) {
	// 期限が遠い未来 → GREEN
	got := CalcDelayStatus("In Progress", strPtr("2030-12-31"), testNow)
	if got != "GREEN" {
		t.Errorf("expected GREEN for far future due date, got %s", got)
	}
}

func TestCalcDelayStatus_InvalidDate(t *testing.T) {
	// パース不能な日付は期限未設定と同様に YELLOW
	got := CalcDelayStatus("To Do", strPtr("not-a-date"), testNow)
	if got != "YELLOW" {
		t.Errorf("expected YELLOW for invalid date, got %s", got)
	}
}

// ----------------------------------------------------------------
// ConvertProject
// ----------------------------------------------------------------

func TestConvertProject_WithLead(t *testing.T) {
	p := jiraclient.Project{
		ID:   "10000",
		Key:  "PROJ",
		Name: "Test Project",
		Lead: &jiraclient.User{
			AccountID:    "acc-123",
			EmailAddress: "lead@example.com",
		},
	}
	got := ConvertProject(p)
	if got.JiraProjectID != "10000" {
		t.Errorf("JiraProjectID: got %s", got.JiraProjectID)
	}
	if got.Key != "PROJ" {
		t.Errorf("Key: got %s", got.Key)
	}
	if got.LeadAccountID != "acc-123" {
		t.Errorf("LeadAccountID: got %s", got.LeadAccountID)
	}
	if got.LeadEmail != "lead@example.com" {
		t.Errorf("LeadEmail: got %s", got.LeadEmail)
	}
}

func TestConvertProject_NoLead(t *testing.T) {
	p := jiraclient.Project{ID: "1", Key: "P", Name: "N", Lead: nil}
	got := ConvertProject(p)
	if got.LeadAccountID != "" || got.LeadEmail != "" {
		t.Errorf("expected empty lead fields, got AccountID=%q Email=%q", got.LeadAccountID, got.LeadEmail)
	}
}

// ----------------------------------------------------------------
// ConvertIssue
// ----------------------------------------------------------------

func makeTestIssue() jiraclient.Issue {
	return jiraclient.Issue{
		ID:  "10001",
		Key: "PROJ-1",
		Fields: jiraclient.IssueFields{
			Summary: "Fix the bug",
			Status: jiraclient.IssueStatus{
				Name: "In Progress",
				StatusCategory: jiraclient.IssueStatusCategory{
					Key:  "indeterminate",
					Name: "In Progress",
				},
			},
			Priority:  &jiraclient.IssuePriority{Name: "High"},
			IssueType: jiraclient.IssueType{Name: "Bug"},
			Assignee: &jiraclient.User{
				AccountID:   "acc-456",
				DisplayName: "Taro Yamada",
			},
			DueDate: "2026-02-28",
			Updated: "2026-02-20T10:00:00+09:00",
			Project: jiraclient.IssueProject{ID: "10000", Key: "PROJ"},
		},
	}
}

func TestConvertIssue_BasicFields(t *testing.T) {
	issue := makeTestIssue()
	got := ConvertIssue(issue, testNow)

	if got.JiraIssueID != "10001" {
		t.Errorf("JiraIssueID: got %s", got.JiraIssueID)
	}
	if got.JiraIssueKey != "PROJ-1" {
		t.Errorf("JiraIssueKey: got %s", got.JiraIssueKey)
	}
	if got.JiraProjectID != "10000" {
		t.Errorf("JiraProjectID: got %s", got.JiraProjectID)
	}
	if got.Summary != "Fix the bug" {
		t.Errorf("Summary: got %s", got.Summary)
	}
	if got.StatusCategory != "In Progress" {
		t.Errorf("StatusCategory: got %s", got.StatusCategory)
	}
	if got.Priority != "High" {
		t.Errorf("Priority: got %s", got.Priority)
	}
	if got.IssueType != "Bug" {
		t.Errorf("IssueType: got %s", got.IssueType)
	}
	if got.AssigneeName != "Taro Yamada" {
		t.Errorf("AssigneeName: got %s", got.AssigneeName)
	}
	if got.DueDate == nil || *got.DueDate != "2026-02-28" {
		t.Errorf("DueDate: got %v", got.DueDate)
	}
	if got.DelayStatus != "GREEN" { // 2026-02-28 is 4 days from testNow → GREEN
		t.Errorf("DelayStatus: expected GREEN, got %s", got.DelayStatus)
	}
}

func TestConvertIssue_NilAssigneeAndPriority(t *testing.T) {
	issue := makeTestIssue()
	issue.Fields.Assignee = nil
	issue.Fields.Priority = nil
	got := ConvertIssue(issue, testNow)

	if got.AssigneeName != "" || got.AssigneeAccountID != "" {
		t.Errorf("expected empty assignee fields")
	}
	if got.Priority != "" {
		t.Errorf("expected empty priority, got %s", got.Priority)
	}
}

func TestConvertIssue_EmptyDueDate(t *testing.T) {
	issue := makeTestIssue()
	issue.Fields.DueDate = ""
	got := ConvertIssue(issue, testNow)

	if got.DueDate != nil {
		t.Errorf("expected nil DueDate, got %v", got.DueDate)
	}
	if got.DelayStatus != "YELLOW" {
		t.Errorf("expected YELLOW for no due date, got %s", got.DelayStatus)
	}
}

func TestConvertIssue_LastUpdatedAt(t *testing.T) {
	issue := makeTestIssue()
	issue.Fields.Updated = "2026-02-20T10:00:00+09:00"
	got := ConvertIssue(issue, testNow)

	expected := time.Date(2026, 2, 20, 1, 0, 0, 0, time.UTC) // +09:00 → UTC
	if !got.LastUpdatedAt.Equal(expected) {
		t.Errorf("LastUpdatedAt: expected %v, got %v", expected, got.LastUpdatedAt)
	}
}

func TestConvertIssue_InvalidUpdatedAt(t *testing.T) {
	issue := makeTestIssue()
	issue.Fields.Updated = "not-a-timestamp"
	got := ConvertIssue(issue, testNow)

	// パース失敗時はゼロ値になる（クラッシュしない）
	if !got.LastUpdatedAt.IsZero() {
		t.Errorf("expected zero LastUpdatedAt for invalid timestamp")
	}
}
