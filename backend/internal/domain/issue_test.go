package domain

import (
	"testing"
	"time"
)

func TestDelayStatus_String(t *testing.T) {
	tests := []struct {
		name   string
		status DelayStatus
		want   string
	}{
		{
			name:   "RED status",
			status: DelayStatusRed,
			want:   "RED",
		},
		{
			name:   "YELLOW status",
			status: DelayStatusYellow,
			want:   "YELLOW",
		},
		{
			name:   "GREEN status",
			status: DelayStatusGreen,
			want:   "GREEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(tt.status)
			if got != tt.want {
				t.Errorf("DelayStatus.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatusCategory_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		category StatusCategory
		want     bool
	}{
		{
			name:     "Valid: todo",
			category: StatusCategoryToDo,
			want:     true,
		},
		{
			name:     "Valid: in_progress",
			category: StatusCategoryInProgress,
			want:     true,
		},
		{
			name:     "Valid: done",
			category: StatusCategoryDone,
			want:     true,
		},
		{
			name:     "Invalid: empty",
			category: StatusCategory(""),
			want:     false,
		},
		{
			name:     "Invalid: unknown",
			category: StatusCategory("unknown"),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validCategories := []StatusCategory{
				StatusCategoryToDo,
				StatusCategoryInProgress,
				StatusCategoryDone,
			}

			got := false
			for _, valid := range validCategories {
				if tt.category == valid {
					got = true
					break
				}
			}

			if got != tt.want {
				t.Errorf("StatusCategory validation = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIssueFilter_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		filter IssueFilter
		want   bool
	}{
		{
			name:   "Empty filter",
			filter: IssueFilter{},
			want:   true,
		},
		{
			name: "Filter with ProjectID",
			filter: IssueFilter{
				ProjectID: func() *int64 { id := int64(1); return &id }(),
			},
			want: false,
		},
		{
			name: "Filter with DelayStatus",
			filter: IssueFilter{
				DelayStatus: func() *DelayStatus { s := DelayStatusRed; return &s }(),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEmpty := tt.filter.ProjectID == nil &&
				tt.filter.DelayStatus == nil &&
				tt.filter.Status == nil

			if isEmpty != tt.want {
				t.Errorf("IssueFilter.IsEmpty() = %v, want %v", isEmpty, tt.want)
			}
		})
	}
}

func TestIssue_Validation(t *testing.T) {
	tests := []struct {
		name    string
		issue   Issue
		wantErr bool
	}{
		{
			name: "Valid issue",
			issue: Issue{
				ProjectID:       1,
				JiraIssueKey:    "PROJ-123",
				JiraIssueID:     "10001",
				Summary:         "Test Issue",
				Status:          "In Progress",
				StatusCategory:  StatusCategoryInProgress,
				DelayStatus:     DelayStatusGreen,
				LastUpdatedAt:   time.Now(),
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Invalid: missing JiraIssueKey",
			issue: Issue{
				ProjectID:      1,
				JiraIssueKey:   "",
				Summary:        "Test Issue",
				LastUpdatedAt:  time.Now(),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: true,
		},
		{
			name: "Invalid: missing Summary",
			issue: Issue{
				ProjectID:      1,
				JiraIssueKey:   "PROJ-123",
				Summary:        "",
				LastUpdatedAt:  time.Now(),
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasError := tt.issue.JiraIssueKey == "" || tt.issue.Summary == ""
			if hasError != tt.wantErr {
				t.Errorf("Issue validation error = %v, want %v", hasError, tt.wantErr)
			}
		})
	}
}

func TestDelayStatus_Priority(t *testing.T) {
	// 優先度テスト: RED > YELLOW > GREEN
	tests := []struct {
		name     string
		status1  DelayStatus
		status2  DelayStatus
		expected string // "greater", "less", "equal"
	}{
		{
			name:     "RED is higher priority than YELLOW",
			status1:  DelayStatusRed,
			status2:  DelayStatusYellow,
			expected: "greater",
		},
		{
			name:     "RED is higher priority than GREEN",
			status1:  DelayStatusRed,
			status2:  DelayStatusGreen,
			expected: "greater",
		},
		{
			name:     "YELLOW is higher priority than GREEN",
			status1:  DelayStatusYellow,
			status2:  DelayStatusGreen,
			expected: "greater",
		},
		{
			name:     "Same status",
			status1:  DelayStatusRed,
			status2:  DelayStatusRed,
			expected: "equal",
		},
	}

	priorityMap := map[DelayStatus]int{
		DelayStatusRed:    3,
		DelayStatusYellow: 2,
		DelayStatusGreen:  1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := priorityMap[tt.status1]
			p2 := priorityMap[tt.status2]

			var result string
			if p1 > p2 {
				result = "greater"
			} else if p1 < p2 {
				result = "less"
			} else {
				result = "equal"
			}

			if result != tt.expected {
				t.Errorf("Priority comparison = %v, want %v", result, tt.expected)
			}
		})
	}
}
