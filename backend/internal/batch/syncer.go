package batch

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/jiraclient"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/normalizer"
)

const defaultWorkerCount = 5

// JiraClient defines the Jira API operations used by the syncer.
// This interface allows the syncer to be tested without a real Jira instance.
type JiraClient interface {
	GetAllProjects() ([]jiraclient.Project, error)
	SearchIssues(opts jiraclient.IssueSearchOptions) ([]jiraclient.Issue, error)
}

// Syncer orchestrates the Jira → DB synchronization process.
type Syncer struct {
	jira        JiraClient
	repo        Repository
	log         *zap.Logger
	workerCount int
}

// NewSyncer creates a new Syncer. workerCount controls the number of concurrent
// per-project issue fetch goroutines; pass 0 to use the default (5).
func NewSyncer(jira JiraClient, repo Repository, log *zap.Logger, workerCount int) *Syncer {
	if workerCount <= 0 {
		workerCount = defaultWorkerCount
	}
	return &Syncer{jira: jira, repo: repo, log: log, workerCount: workerCount}
}

// RunFullSync fetches all Jira projects and their issues, then upserts them into the DB.
// It records execution details in the sync_logs table.
func (s *Syncer) RunFullSync(ctx context.Context) error {
	start := time.Now()
	s.log.Info("full sync started")

	logID, err := s.repo.StartSyncLog(ctx, "FULL")
	if err != nil {
		return fmt.Errorf("start sync log: %w", err)
	}

	projectsSynced, issuesSynced, syncErr := s.runSync(ctx)

	// sync_logs の更新は sync 失敗でも必ず行う
	status := "SUCCESS"
	errMsg := ""
	if syncErr != nil {
		status = "FAILURE"
		errMsg = syncErr.Error()
	}

	if finishErr := s.repo.FinishSyncLog(ctx, logID, status, projectsSynced, issuesSynced, errMsg); finishErr != nil {
		s.log.Error("failed to finish sync log", zap.Error(finishErr))
	}

	s.log.Info("full sync finished",
		zap.String("status", status),
		zap.Int("projects_synced", projectsSynced),
		zap.Int("issues_synced", issuesSynced),
		zap.Duration("duration", time.Since(start)),
	)

	return syncErr
}

// RunDeltaSync fetches only the issues updated since the last successful DELTA sync
// and upserts them into the DB. It records execution details in the sync_logs table.
func (s *Syncer) RunDeltaSync(ctx context.Context) error {
	start := time.Now()
	s.log.Info("delta sync started")

	logID, err := s.repo.StartSyncLog(ctx, "DELTA")
	if err != nil {
		return fmt.Errorf("start sync log: %w", err)
	}

	issuesSynced, syncErr := s.runDeltaSync(ctx)

	// sync_logs の更新は sync 失敗でも必ず行う
	status := "SUCCESS"
	errMsg := ""
	if syncErr != nil {
		status = "FAILURE"
		errMsg = syncErr.Error()
	}

	if finishErr := s.repo.FinishSyncLog(ctx, logID, status, 0, issuesSynced, errMsg); finishErr != nil {
		s.log.Error("failed to finish sync log", zap.Error(finishErr))
	}

	s.log.Info("delta sync finished",
		zap.String("status", status),
		zap.Int("issues_synced", issuesSynced),
		zap.Duration("duration", time.Since(start)),
	)

	return syncErr
}

// runDeltaSync fetches issues updated since the last successful DELTA sync and upserts them.
func (s *Syncer) runDeltaSync(ctx context.Context) (issuesSynced int, err error) {
	// 1. 前回成功した DELTA sync の実行時刻を取得
	lastSync, err := s.repo.GetLastSuccessfulSyncTime(ctx, "DELTA")
	if err != nil {
		return 0, fmt.Errorf("get last successful sync time: %w", err)
	}

	// フォールバック: 前回記録がなければ1時間前を使用
	var since time.Time
	if lastSync != nil {
		since = *lastSync
	} else {
		since = time.Now().Add(-1 * time.Hour)
		s.log.Info("no previous delta sync found, using 1-hour fallback")
	}

	// 2. JQL でプロジェクト横断の差分チケットを一括取得
	jql := fmt.Sprintf(`updated >= "%s" ORDER BY updated ASC`, since.Format("2006/01/02 15:04"))
	s.log.Info("fetching delta issues",
		zap.String("since", since.Format(time.RFC3339)),
		zap.String("jql", jql),
	)

	issues, err := s.jira.SearchIssues(jiraclient.IssueSearchOptions{JQL: jql})
	if err != nil {
		return 0, fmt.Errorf("search delta issues: %w", err)
	}
	s.log.Info("fetched delta issues", zap.Int("count", len(issues)))

	if len(issues) == 0 {
		return 0, nil
	}

	// 3. チケットに対応する project_id を解決
	projectIDMap, err := s.repo.GetProjectIDMap(ctx)
	if err != nil {
		return 0, fmt.Errorf("get project id map: %w", err)
	}

	// 4. チケットを正規化して DB に upsert
	now := normalizer.Now()
	dbIssues := make([]normalizer.DBIssue, len(issues))
	for i, issue := range issues {
		dbIssues[i] = normalizer.ConvertIssue(issue, now)
	}
	issuesSynced, err = s.repo.UpsertIssues(ctx, dbIssues, projectIDMap)
	if err != nil {
		return 0, fmt.Errorf("upsert issues: %w", err)
	}

	return issuesSynced, nil
}

// runSync is the core sync logic, separated for testability.
// Returns the number of projects and issues synced, plus any error.
func (s *Syncer) runSync(ctx context.Context) (projectsSynced, issuesSynced int, err error) {
	// 1. プロジェクト一覧を取得
	s.log.Info("fetching projects from Jira")
	jiraProjects, err := s.jira.GetAllProjects()
	if err != nil {
		return 0, 0, fmt.Errorf("get projects: %w", err)
	}
	s.log.Info("fetched projects", zap.Int("count", len(jiraProjects)))

	// 2. プロジェクトを正規化して DB に upsert
	dbProjects := make([]normalizer.DBProject, len(jiraProjects))
	for i, p := range jiraProjects {
		dbProjects[i] = normalizer.ConvertProject(p)
	}
	projectsSynced, err = s.repo.UpsertProjects(ctx, dbProjects)
	if err != nil {
		return 0, 0, fmt.Errorf("upsert projects: %w", err)
	}

	// 3. jira_project_id → DB id のマップを取得
	projectIDMap, err := s.repo.GetProjectIDMap(ctx)
	if err != nil {
		return projectsSynced, 0, fmt.Errorf("get project id map: %w", err)
	}

	// 4. プロジェクトごとのチケット取得を並列化（worker pool）
	allIssues := s.fetchIssuesParallel(ctx, jiraProjects)
	s.log.Info("fetched issues", zap.Int("count", len(allIssues)))

	// 5. チケットを正規化して DB に upsert
	now := normalizer.Now()
	dbIssues := make([]normalizer.DBIssue, len(allIssues))
	for i, issue := range allIssues {
		dbIssues[i] = normalizer.ConvertIssue(issue, now)
	}
	issuesSynced, err = s.repo.UpsertIssues(ctx, dbIssues, projectIDMap)
	if err != nil {
		return projectsSynced, 0, fmt.Errorf("upsert issues: %w", err)
	}

	return projectsSynced, issuesSynced, nil
}

// fetchIssuesParallel fetches issues for all projects concurrently using a worker pool.
// Errors from individual projects are logged as warnings; processing continues.
func (s *Syncer) fetchIssuesParallel(ctx context.Context, projects []jiraclient.Project) []jiraclient.Issue {
	type result struct {
		issues []jiraclient.Issue
	}

	// semaphore で同時実行数を制限する
	sem := make(chan struct{}, s.workerCount)
	resultCh := make(chan result, len(projects))

	var wg sync.WaitGroup
	for _, p := range projects {
		p := p // ループ変数キャプチャ
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release

			// context がキャンセルされていたらスキップ
			select {
			case <-ctx.Done():
				return
			default:
			}

			jql := fmt.Sprintf("project = %s ORDER BY updated ASC", p.Key)
			issues, err := s.jira.SearchIssues(jiraclient.IssueSearchOptions{JQL: jql})
			if err != nil {
				// 1プロジェクトの失敗は警告扱いとし、他プロジェクトの処理を継続する
				s.log.Warn("failed to fetch issues for project",
					zap.String("project_key", p.Key),
					zap.Error(err),
				)
				return
			}
			resultCh <- result{issues: issues}
		}()
	}

	// 全 goroutine の完了後にチャネルを閉じる
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var all []jiraclient.Issue
	for r := range resultCh {
		all = append(all, r.issues...)
	}
	return all
}
