package secrets

import (
	"testing"
)

func TestLoadJira_Success(t *testing.T) {
	t.Setenv("JIRA_BASE_URL", "https://example.atlassian.net/")
	t.Setenv("JIRA_EMAIL", "user@example.com")
	t.Setenv("JIRA_API_TOKEN", "token123")

	creds, err := LoadJira()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	// trailing slash should be trimmed
	if creds.BaseURL != "https://example.atlassian.net" {
		t.Errorf("expected trimmed BaseURL, got: %s", creds.BaseURL)
	}
	if creds.Email != "user@example.com" {
		t.Errorf("expected email, got: %s", creds.Email)
	}
	if creds.APIToken != "token123" {
		t.Errorf("expected token, got: %s", creds.APIToken)
	}
}

func TestLoadJira_MissingAll(t *testing.T) {
	t.Setenv("JIRA_BASE_URL", "")
	t.Setenv("JIRA_EMAIL", "")
	t.Setenv("JIRA_API_TOKEN", "")

	_, err := LoadJira()
	if err == nil {
		t.Fatal("expected error for missing credentials, got nil")
	}
}

func TestLoadJira_MissingPartial(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		email    string
		apiToken string
	}{
		{"missing BaseURL", "", "user@example.com", "token"},
		{"missing Email", "https://example.atlassian.net", "", "token"},
		{"missing APIToken", "https://example.atlassian.net", "user@example.com", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("JIRA_BASE_URL", tc.baseURL)
			t.Setenv("JIRA_EMAIL", tc.email)
			t.Setenv("JIRA_API_TOKEN", tc.apiToken)

			_, err := LoadJira()
			if err == nil {
				t.Errorf("expected error for %s, got nil", tc.name)
			}
		})
	}
}
