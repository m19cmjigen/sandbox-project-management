package jiraclient

import "fmt"

const (
	issueSearchPath = "/rest/api/3/issue/search"
	issuePageSize   = 100
)

// defaultIssueFields is the list of fields requested from Jira for each issue.
var defaultIssueFields = []string{
	"summary",
	"status",
	"priority",
	"issuetype",
	"assignee",
	"duedate",
	"updated",
	"project",
}

// IssueSearchOptions contains optional filters for SearchIssues.
type IssueSearchOptions struct {
	// JQL is the Jira Query Language expression used to filter issues.
	JQL string
	// Fields overrides the default list of fields to retrieve.
	Fields []string
}

// SearchIssues fetches all issues matching the given JQL query, handling pagination
// automatically. If opts.JQL is empty, all issues are returned.
func (c *Client) SearchIssues(opts IssueSearchOptions) ([]Issue, error) {
	fields := opts.Fields
	if len(fields) == 0 {
		fields = defaultIssueFields
	}

	var all []Issue
	startAt := 0

	for {
		req := IssueSearchRequest{
			JQL:        opts.JQL,
			StartAt:    startAt,
			MaxResults: issuePageSize,
			Fields:     fields,
		}

		var resp IssueSearchResponse
		if err := c.post(issueSearchPath, req, &resp); err != nil {
			return nil, fmt.Errorf("search issues (startAt=%d): %w", startAt, err)
		}

		all = append(all, resp.Issues...)

		if startAt+len(resp.Issues) >= resp.Total || len(resp.Issues) == 0 {
			break
		}
		startAt += len(resp.Issues)
	}

	return all, nil
}

// SearchIssuesUpdatedAfter returns all issues updated after the given RFC3339 timestamp.
// This is used for delta sync to retrieve only recently changed issues.
func (c *Client) SearchIssuesUpdatedAfter(since string) ([]Issue, error) {
	// JQL: updated >= "YYYY/MM/DD HH:MM" の形式を使用
	jql := fmt.Sprintf(`updated >= "%s" ORDER BY updated ASC`, since)
	return c.SearchIssues(IssueSearchOptions{JQL: jql})
}
