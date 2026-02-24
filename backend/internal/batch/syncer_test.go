package batch

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/jiraclient"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/normalizer"
)

// ----------------------------------------------------------------
// Mock JiraClient
// ----------------------------------------------------------------

type mockJiraClient struct {
	projects    []jiraclient.Project
	projectsErr error
	issues      []jiraclient.Issue
	issuesErr   error
	// searchCallCount は SearchIssues が呼ばれた回数
	searchCallCount int64
}

func (m *mockJiraClient) GetAllProjects() ([]jiraclient.Project, error) {
	return m.projects, m.projectsErr
}

func (m *mockJiraClient) SearchIssues(_ jiraclient.IssueSearchOptions) ([]jiraclient.Issue, error) {
	atomic.AddInt64(&m.searchCallCount, 1)
	return m.issues, m.issuesErr
}

// ----------------------------------------------------------------
// Mock Repository
// ----------------------------------------------------------------

type mockRepository struct {
	upsertProjectsCount int
	upsertIssuesCount   int
	projectIDMap        map[string]int64
	syncLogID           int64
	finishedStatus      string
	finishedErrMsg      string

	upsertProjectsErr error
	upsertIssuesErr   error
	getProjectMapErr  error
	startLogErr       error
	finishLogErr      error
	lastSyncTime      *time.Time
	lastSyncTimeErr   error
}

func (m *mockRepository) UpsertProjects(_ context.Context, projects []normalizer.DBProject) (int, error) {
	m.upsertProjectsCount = len(projects)
	return len(projects), m.upsertProjectsErr
}

func (m *mockRepository) UpsertIssues(_ context.Context, issues []normalizer.DBIssue, _ map[string]int64) (int, error) {
	m.upsertIssuesCount = len(issues)
	return len(issues), m.upsertIssuesErr
}

func (m *mockRepository) GetProjectIDMap(_ context.Context) (map[string]int64, error) {
	return m.projectIDMap, m.getProjectMapErr
}

func (m *mockRepository) StartSyncLog(_ context.Context, _ string) (int64, error) {
	return m.syncLogID, m.startLogErr
}

func (m *mockRepository) FinishSyncLog(_ context.Context, _ int64, status string, _, _ int, errMsg string) error {
	m.finishedStatus = status
	m.finishedErrMsg = errMsg
	return m.finishLogErr
}

func (m *mockRepository) GetLastSuccessfulSyncTime(_ context.Context, _ string) (*time.Time, error) {
	return m.lastSyncTime, m.lastSyncTimeErr
}

// ----------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------

func newTestSyncer(jira JiraClient, repo Repository) *Syncer {
	log := zap.NewNop()
	return NewSyncer(jira, repo, log, 2)
}

func makeProject(id, key string) jiraclient.Project {
	return jiraclient.Project{ID: id, Key: key, Name: "Project " + key}
}

func makeIssue(id, key, projectID string) jiraclient.Issue {
	return jiraclient.Issue{
		ID:  id,
		Key: key,
		Fields: jiraclient.IssueFields{
			Summary:   "Summary " + key,
			Status:    jiraclient.IssueStatus{StatusCategory: jiraclient.IssueStatusCategory{Key: "new"}},
			IssueType: jiraclient.IssueType{Name: "Bug"},
			Project:   jiraclient.IssueProject{ID: projectID, Key: "P"},
		},
	}
}

// ----------------------------------------------------------------
// Tests
// ----------------------------------------------------------------

func TestRunFullSync_Success(t *testing.T) {
	jira := &mockJiraClient{
		projects: []jiraclient.Project{makeProject("10", "PROJ")},
		issues:   []jiraclient.Issue{makeIssue("1", "PROJ-1", "10")},
	}
	repo := &mockRepository{
		syncLogID:    42,
		projectIDMap: map[string]int64{"10": 1},
	}

	syncer := newTestSyncer(jira, repo)
	err := syncer.RunFullSync(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if repo.upsertProjectsCount != 1 {
		t.Errorf("expected 1 project upserted, got %d", repo.upsertProjectsCount)
	}
	if repo.upsertIssuesCount != 1 {
		t.Errorf("expected 1 issue upserted, got %d", repo.upsertIssuesCount)
	}
	if repo.finishedStatus != "SUCCESS" {
		t.Errorf("expected SUCCESS status, got %s", repo.finishedStatus)
	}
}

func TestRunFullSync_NoProjects(t *testing.T) {
	jira := &mockJiraClient{projects: []jiraclient.Project{}}
	repo := &mockRepository{syncLogID: 1, projectIDMap: map[string]int64{}}

	syncer := newTestSyncer(jira, repo)
	err := syncer.RunFullSync(context.Background())
	if err != nil {
		t.Fatalf("expected no error for empty projects, got: %v", err)
	}
	if repo.finishedStatus != "SUCCESS" {
		t.Errorf("expected SUCCESS, got %s", repo.finishedStatus)
	}
}

func TestRunFullSync_JiraProjectsFetchError(t *testing.T) {
	jira := &mockJiraClient{projectsErr: errors.New("jira down")}
	repo := &mockRepository{syncLogID: 1}

	syncer := newTestSyncer(jira, repo)
	err := syncer.RunFullSync(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if repo.finishedStatus != "FAILURE" {
		t.Errorf("expected FAILURE status, got %s", repo.finishedStatus)
	}
	if repo.finishedErrMsg == "" {
		t.Error("expected non-empty error message in sync log")
	}
}

func TestRunFullSync_UpsertProjectsError(t *testing.T) {
	jira := &mockJiraClient{projects: []jiraclient.Project{makeProject("1", "P")}}
	repo := &mockRepository{
		syncLogID:         1,
		upsertProjectsErr: errors.New("db write error"),
	}

	syncer := newTestSyncer(jira, repo)
	err := syncer.RunFullSync(context.Background())
	if err == nil {
		t.Fatal("expected error from upsert projects, got nil")
	}
	if repo.finishedStatus != "FAILURE" {
		t.Errorf("expected FAILURE status, got %s", repo.finishedStatus)
	}
}

func TestRunFullSync_IssuesFetchErrorContinues(t *testing.T) {
	// チケット取得エラーは警告扱いで処理継続し SUCCESS になる
	jira := &mockJiraClient{
		projects:  []jiraclient.Project{makeProject("1", "P"), makeProject("2", "Q")},
		issuesErr: errors.New("project not found"),
	}
	repo := &mockRepository{
		syncLogID:    1,
		projectIDMap: map[string]int64{"1": 1, "2": 2},
	}

	syncer := newTestSyncer(jira, repo)
	err := syncer.RunFullSync(context.Background())
	if err != nil {
		t.Fatalf("expected no error when issue fetch fails per project, got: %v", err)
	}
	if repo.finishedStatus != "SUCCESS" {
		t.Errorf("expected SUCCESS despite issue fetch errors, got %s", repo.finishedStatus)
	}
}

func TestRunFullSync_ParallelFetch(t *testing.T) {
	// 複数プロジェクトのチケット取得が並列実行されることを検証
	projects := []jiraclient.Project{
		makeProject("1", "P1"),
		makeProject("2", "P2"),
		makeProject("3", "P3"),
	}
	jira := &mockJiraClient{
		projects: projects,
		issues:   []jiraclient.Issue{makeIssue("i1", "P1-1", "1")},
	}
	repo := &mockRepository{
		syncLogID:    1,
		projectIDMap: map[string]int64{"1": 1, "2": 2, "3": 3},
	}

	syncer := newTestSyncer(jira, repo)
	if err := syncer.RunFullSync(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 3プロジェクト分 SearchIssues が呼ばれたはず
	if got := atomic.LoadInt64(&jira.searchCallCount); got != 3 {
		t.Errorf("expected SearchIssues called 3 times, got %d", got)
	}
}

func TestRunFullSync_SyncLogAlwaysFinished(t *testing.T) {
	// Jira エラーでも sync_log は必ず更新される
	jira := &mockJiraClient{projectsErr: errors.New("timeout")}
	repo := &mockRepository{syncLogID: 99}

	syncer := newTestSyncer(jira, repo)
	syncer.RunFullSync(context.Background())

	if repo.finishedStatus == "" {
		t.Error("FinishSyncLog must be called even when sync fails")
	}
}

func TestNewSyncer_DefaultWorkerCount(t *testing.T) {
	s := NewSyncer(nil, nil, zap.NewNop(), 0)
	if s.workerCount != defaultWorkerCount {
		t.Errorf("expected default worker count %d, got %d", defaultWorkerCount, s.workerCount)
	}
}

// ----------------------------------------------------------------
// Delta Sync Tests
// ----------------------------------------------------------------

func TestRunDeltaSync_Success(t *testing.T) {
	lastSync := time.Now().Add(-30 * time.Minute)
	jira := &mockJiraClient{
		issues: []jiraclient.Issue{makeIssue("1", "PROJ-1", "10")},
	}
	repo := &mockRepository{
		syncLogID:    42,
		lastSyncTime: &lastSync,
		projectIDMap: map[string]int64{"10": 1},
	}

	syncer := newTestSyncer(jira, repo)
	err := syncer.RunDeltaSync(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if repo.upsertIssuesCount != 1 {
		t.Errorf("expected 1 issue upserted, got %d", repo.upsertIssuesCount)
	}
	if repo.finishedStatus != "SUCCESS" {
		t.Errorf("expected SUCCESS status, got %s", repo.finishedStatus)
	}
}

func TestRunDeltaSync_NoLastSync_UsesFallback(t *testing.T) {
	// 前回の記録がない場合は1時間前をフォールバックとして使用し、正常に完了する
	jira := &mockJiraClient{
		issues: []jiraclient.Issue{makeIssue("1", "PROJ-1", "10")},
	}
	repo := &mockRepository{
		syncLogID:    1,
		lastSyncTime: nil, // 前回記録なし
		projectIDMap: map[string]int64{"10": 1},
	}

	syncer := newTestSyncer(jira, repo)
	err := syncer.RunDeltaSync(context.Background())
	if err != nil {
		t.Fatalf("expected no error with fallback, got: %v", err)
	}
	if repo.finishedStatus != "SUCCESS" {
		t.Errorf("expected SUCCESS, got %s", repo.finishedStatus)
	}
}

func TestRunDeltaSync_NoIssues(t *testing.T) {
	lastSync := time.Now().Add(-10 * time.Minute)
	jira := &mockJiraClient{issues: []jiraclient.Issue{}}
	repo := &mockRepository{
		syncLogID:    1,
		lastSyncTime: &lastSync,
		projectIDMap: map[string]int64{},
	}

	syncer := newTestSyncer(jira, repo)
	err := syncer.RunDeltaSync(context.Background())
	if err != nil {
		t.Fatalf("expected no error for empty issues, got: %v", err)
	}
	if repo.upsertIssuesCount != 0 {
		t.Errorf("expected 0 issues upserted, got %d", repo.upsertIssuesCount)
	}
	if repo.finishedStatus != "SUCCESS" {
		t.Errorf("expected SUCCESS, got %s", repo.finishedStatus)
	}
}

func TestRunDeltaSync_SearchError(t *testing.T) {
	lastSync := time.Now().Add(-1 * time.Hour)
	jira := &mockJiraClient{issuesErr: errors.New("jira search failed")}
	repo := &mockRepository{
		syncLogID:    1,
		lastSyncTime: &lastSync,
	}

	syncer := newTestSyncer(jira, repo)
	err := syncer.RunDeltaSync(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if repo.finishedStatus != "FAILURE" {
		t.Errorf("expected FAILURE status, got %s", repo.finishedStatus)
	}
	if repo.finishedErrMsg == "" {
		t.Error("expected non-empty error message in sync log")
	}
}

func TestRunDeltaSync_SyncLogAlwaysFinished(t *testing.T) {
	// Jira エラーでも sync_log は必ず更新される
	jira := &mockJiraClient{issuesErr: errors.New("timeout")}
	repo := &mockRepository{syncLogID: 99}

	syncer := newTestSyncer(jira, repo)
	syncer.RunDeltaSync(context.Background())

	if repo.finishedStatus == "" {
		t.Error("FinishSyncLog must be called even when delta sync fails")
	}
}
