package jiraclient

// Project represents a Jira project as returned by the project search API.
type Project struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	Lead *User  `json:"lead,omitempty"`
}

// User represents a Jira user (project lead or assignee).
type User struct {
	AccountID    string `json:"accountId"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// ProjectSearchResponse is the response from GET /rest/api/3/project/search.
type ProjectSearchResponse struct {
	Values     []Project `json:"values"`
	StartAt    int       `json:"startAt"`
	MaxResults int       `json:"maxResults"`
	Total      int       `json:"total"`
	IsLast     bool      `json:"isLast"`
}

// Issue represents a Jira issue as returned by the issue search API.
type Issue struct {
	ID     string      `json:"id"`
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

// IssueFields contains the fields of a Jira issue.
type IssueFields struct {
	Summary        string          `json:"summary"`
	Status         IssueStatus     `json:"status"`
	Priority       *IssuePriority  `json:"priority,omitempty"`
	IssueType      IssueType       `json:"issuetype"`
	Assignee       *User           `json:"assignee,omitempty"`
	DueDate        string          `json:"duedate"` // "YYYY-MM-DD" or ""
	Updated        string          `json:"updated"`
	Project        IssueProject    `json:"project"`
}

// IssueStatus represents the status of a Jira issue.
type IssueStatus struct {
	Name           string             `json:"name"`
	StatusCategory IssueStatusCategory `json:"statusCategory"`
}

// IssueStatusCategory represents the category of a Jira issue status.
type IssueStatusCategory struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

// IssuePriority represents the priority of a Jira issue.
type IssuePriority struct {
	Name string `json:"name"`
}

// IssueType represents the type of a Jira issue.
type IssueType struct {
	Name string `json:"name"`
}

// IssueProject contains project information embedded in an issue.
type IssueProject struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// IssueSearchRequest is the body for POST /rest/api/3/issue/search.
type IssueSearchRequest struct {
	JQL        string   `json:"jql"`
	StartAt    int      `json:"startAt"`
	MaxResults int      `json:"maxResults"`
	Fields     []string `json:"fields"`
}

// IssueSearchResponse is the response from POST /rest/api/3/issue/search.
type IssueSearchResponse struct {
	Issues     []Issue `json:"issues"`
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
}
