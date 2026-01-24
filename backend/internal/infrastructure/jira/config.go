package jira

import (
	"fmt"
	"os"
)

// Config holds Jira API configuration
type Config struct {
	BaseURL  string
	Email    string
	APIToken string
}

// LoadConfig loads Jira configuration from environment variables
func LoadConfig() (*Config, error) {
	baseURL := os.Getenv("JIRA_BASE_URL")
	email := os.Getenv("JIRA_EMAIL")
	apiToken := os.Getenv("JIRA_API_TOKEN")

	if baseURL == "" {
		return nil, fmt.Errorf("JIRA_BASE_URL environment variable is required")
	}
	if email == "" {
		return nil, fmt.Errorf("JIRA_EMAIL environment variable is required")
	}
	if apiToken == "" {
		return nil, fmt.Errorf("JIRA_API_TOKEN environment variable is required")
	}

	return &Config{
		BaseURL:  baseURL,
		Email:    email,
		APIToken: apiToken,
	}, nil
}
