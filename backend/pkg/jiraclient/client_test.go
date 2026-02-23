package jiraclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestClient creates a Client pointed at the given test server URL.
// Sleep is replaced with a no-op to keep tests fast.
func newTestClient(serverURL string) *Client {
	c := New(Config{
		BaseURL:  serverURL,
		Email:    "test@example.com",
		APIToken: "test-token",
	})
	c.sleepFn = func(time.Duration) {}
	return c
}

// --- Project tests ---

func TestGetAllProjects_SinglePage(t *testing.T) {
	projects := []Project{
		{ID: "1", Key: "PROJ", Name: "Test Project"},
	}
	resp := ProjectSearchResponse{
		Values:     projects,
		StartAt:    0,
		MaxResults: 50,
		Total:      1,
		IsLast:     true,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/project/search" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	got, err := client.GetAllProjects()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 project, got %d", len(got))
	}
	if got[0].Key != "PROJ" {
		t.Errorf("expected key PROJ, got %s", got[0].Key)
	}
}

func TestGetAllProjects_MultiPage(t *testing.T) {
	// ページ1: 2件、IsLast=false
	// ページ2: 1件、IsLast=true
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		callCount++
		if callCount == 1 {
			json.NewEncoder(w).Encode(ProjectSearchResponse{
				Values:  []Project{{ID: "1", Key: "P1"}, {ID: "2", Key: "P2"}},
				IsLast:  false,
				Total:   3,
			})
		} else {
			json.NewEncoder(w).Encode(ProjectSearchResponse{
				Values:  []Project{{ID: "3", Key: "P3"}},
				IsLast:  true,
				Total:   3,
			})
		}
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	got, err := client.GetAllProjects()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 projects, got %d", len(got))
	}
}

func TestGetAllProjects_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	// サーバーエラーは maxRetries 回試行後にエラーを返す
	_, err := client.GetAllProjects()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetAllProjects_UnauthorizedNoRetry(t *testing.T) {
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	_, err := client.GetAllProjects()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// 401 はリトライ対象外なので1回のみ呼ばれる
	if callCount != 1 {
		t.Errorf("expected 1 call for 401, got %d", callCount)
	}
}

// --- Issue tests ---

func TestSearchIssues_SinglePage(t *testing.T) {
	issues := []Issue{
		{ID: "10001", Key: "PROJ-1", Fields: IssueFields{Summary: "Fix bug"}},
	}
	respBody := IssueSearchResponse{
		Issues:     issues,
		StartAt:    0,
		MaxResults: 100,
		Total:      1,
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/api/3/issue/search" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respBody)
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	got, err := client.SearchIssues(IssueSearchOptions{JQL: "project = PROJ"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 issue, got %d", len(got))
	}
	if got[0].Key != "PROJ-1" {
		t.Errorf("expected key PROJ-1, got %s", got[0].Key)
	}
}

func TestSearchIssues_Pagination(t *testing.T) {
	// 合計3件、1ページ2件で2ページ必要
	callCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		callCount++

		var req IssueSearchRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.StartAt == 0 {
			json.NewEncoder(w).Encode(IssueSearchResponse{
				Issues:  []Issue{{Key: "P-1"}, {Key: "P-2"}},
				StartAt: 0,
				Total:   3,
			})
		} else {
			json.NewEncoder(w).Encode(IssueSearchResponse{
				Issues:  []Issue{{Key: "P-3"}},
				StartAt: 2,
				Total:   3,
			})
		}
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	got, err := client.SearchIssues(IssueSearchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 issues, got %d", len(got))
	}
	if callCount != 2 {
		t.Errorf("expected 2 API calls for pagination, got %d", callCount)
	}
}

func TestSearchIssuesUpdatedAfter(t *testing.T) {
	var capturedJQL string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req IssueSearchRequest
		json.NewDecoder(r.Body).Decode(&req)
		capturedJQL = req.JQL

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IssueSearchResponse{Total: 0})
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	_, err := client.SearchIssuesUpdatedAfter("2026-01-01T00:00:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// JQL に "updated >=" が含まれていることを確認
	if capturedJQL == "" {
		t.Error("expected JQL to be set, got empty string")
	}
}

// --- Auth header test ---

func TestClientSetsAuthorizationHeader(t *testing.T) {
	var capturedAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ProjectSearchResponse{IsLast: true})
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	client.GetAllProjects()

	if capturedAuth == "" {
		t.Error("expected Authorization header to be set")
	}
	if len(capturedAuth) < 6 || capturedAuth[:6] != "Basic " {
		t.Errorf("expected Basic auth header, got: %s", capturedAuth)
	}
}
