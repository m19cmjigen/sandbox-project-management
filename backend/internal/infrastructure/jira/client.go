package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a Jira Cloud API client
type Client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new Jira API client
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// JiraProject represents a Jira project
type JiraProject struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// JiraIssue represents a Jira issue
type JiraIssue struct {
	ID     string          `json:"id"`
	Key    string          `json:"key"`
	Fields JiraIssueFields `json:"fields"`
}

// JiraIssueFields represents fields in a Jira issue
type JiraIssueFields struct {
	Summary     string           `json:"summary"`
	Status      JiraStatus       `json:"status"`
	Priority    JiraPriority     `json:"priority"`
	Assignee    *JiraUser        `json:"assignee"`
	Created     string           `json:"created"`
	Updated     string           `json:"updated"`
	DueDate     string           `json:"duedate"`
	Description string           `json:"description"`
	Project     JiraProject      `json:"project"`
	Parent      *JiraIssueParent `json:"parent,omitempty"`
}

// JiraStatus represents a Jira issue status
type JiraStatus struct {
	Name string `json:"name"`
}

// JiraPriority represents a Jira issue priority
type JiraPriority struct {
	Name string `json:"name"`
}

// JiraUser represents a Jira user
type JiraUser struct {
	DisplayName string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// JiraIssueParent represents a parent issue
type JiraIssueParent struct {
	ID     string          `json:"id"`
	Key    string          `json:"key"`
	Fields JiraIssueFields `json:"fields"`
}

// SearchResponse represents Jira search API response
type SearchResponse struct {
	Issues     []JiraIssue `json:"issues"`
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
}

// doRequest performs an authenticated HTTP request to Jira API
func (c *Client) doRequest(ctx context.Context, method, endpoint string) (*http.Response, error) {
	url := c.config.BaseURL + endpoint

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication
	req.SetBasicAuth(c.config.Email, c.config.APIToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("jira API error (status %d): %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// GetProjects retrieves all Jira projects
func (c *Client) GetProjects(ctx context.Context) ([]JiraProject, error) {
	resp, err := c.doRequest(ctx, "GET", "/rest/api/3/project")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects []JiraProject
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to decode projects: %w", err)
	}

	return projects, nil
}

// SearchIssues searches for issues using JQL
func (c *Client) SearchIssues(ctx context.Context, jql string, startAt, maxResults int) (*SearchResponse, error) {
	endpoint := fmt.Sprintf(
		"/rest/api/3/search?jql=%s&startAt=%d&maxResults=%d&fields=summary,status,priority,assignee,created,updated,duedate,description,project,parent",
		jql, startAt, maxResults,
	)

	resp, err := c.doRequest(ctx, "GET", endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	return &searchResp, nil
}

// GetAllIssuesForProject retrieves all issues for a specific project
func (c *Client) GetAllIssuesForProject(ctx context.Context, projectKey string) ([]JiraIssue, error) {
	jql := fmt.Sprintf("project = %s ORDER BY created DESC", projectKey)
	maxResults := 100
	startAt := 0
	allIssues := []JiraIssue{}

	for {
		searchResp, err := c.SearchIssues(ctx, jql, startAt, maxResults)
		if err != nil {
			return nil, err
		}

		allIssues = append(allIssues, searchResp.Issues...)

		// Check if we've retrieved all issues
		if startAt+len(searchResp.Issues) >= searchResp.Total {
			break
		}

		startAt += maxResults
	}

	return allIssues, nil
}
