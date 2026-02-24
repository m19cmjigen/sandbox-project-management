// Package secrets provides secure loading of application credentials.
//
// In production, credentials are managed in AWS Secrets Manager and injected
// into ECS task environment variables automatically via the task definition's
// "secrets" field. The application itself only reads standard environment variables,
// keeping the codebase free from AWS SDK dependencies for secret retrieval.
package secrets

import (
	"fmt"
	"os"
	"strings"
)

// JiraCredentials holds the authentication information required to access the Jira Cloud API.
type JiraCredentials struct {
	// BaseURL is the base URL of the Jira Cloud instance (e.g. https://your-org.atlassian.net).
	BaseURL string
	// Email is the Atlassian account email used for Basic authentication.
	Email string
	// APIToken is the Jira API token associated with the email account.
	APIToken string
}

// LoadJira reads Jira API credentials from environment variables and validates them.
//
// Expected environment variables:
//
//	JIRA_BASE_URL  — Jira Cloud base URL
//	JIRA_EMAIL     — Atlassian account email
//	JIRA_API_TOKEN — Jira API token
//
// In production (ECS), these variables are populated via AWS Secrets Manager
// using the task definition "secrets" field. See docs/secrets-management.md.
func LoadJira() (JiraCredentials, error) {
	creds := JiraCredentials{
		BaseURL:  strings.TrimRight(os.Getenv("JIRA_BASE_URL"), "/"),
		Email:    os.Getenv("JIRA_EMAIL"),
		APIToken: os.Getenv("JIRA_API_TOKEN"),
	}

	var missing []string
	if creds.BaseURL == "" {
		missing = append(missing, "JIRA_BASE_URL")
	}
	if creds.Email == "" {
		missing = append(missing, "JIRA_EMAIL")
	}
	if creds.APIToken == "" {
		missing = append(missing, "JIRA_API_TOKEN")
	}

	if len(missing) > 0 {
		return JiraCredentials{}, fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return creds, nil
}
