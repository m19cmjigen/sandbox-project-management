package batch

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/normalizer"
)

// newRepoDB returns a sqlx.DB backed by go-sqlmock for repository tests.
func newRepoDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	t.Cleanup(func() { sqlxDB.Close() })
	return sqlxDB, mock
}

// --- UpsertProjects tests ---

func TestUpsertProjects_Empty(t *testing.T) {
	db, _ := newRepoDB(t)
	repo := NewRepository(db)

	n, err := repo.UpsertProjects(context.Background(), nil)

	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestUpsertProjects_Success(t *testing.T) {
	db, mock := newRepoDB(t)
	repo := NewRepository(db)

	mock.ExpectExec(`INSERT INTO projects`).
		WillReturnResult(sqlmock.NewResult(1, 2))

	projects := []normalizer.DBProject{
		{JiraProjectID: "P1", Key: "KEY1", Name: "Project 1"},
		{JiraProjectID: "P2", Key: "KEY2", Name: "Project 2"},
	}
	n, err := repo.UpsertProjects(context.Background(), projects)

	assert.NoError(t, err)
	assert.Equal(t, 2, n)
}

// --- UpsertIssues tests ---

func TestUpsertIssues_Empty(t *testing.T) {
	db, _ := newRepoDB(t)
	repo := NewRepository(db)

	n, err := repo.UpsertIssues(context.Background(), nil, map[string]int64{})

	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestUpsertIssues_UnknownProject(t *testing.T) {
	db, _ := newRepoDB(t)
	repo := NewRepository(db)

	// All issues have unknown jira_project_id → all skipped → no DB call
	issues := []normalizer.DBIssue{
		{JiraIssueID: "I1", JiraProjectID: "UNKNOWN", Summary: "issue 1", LastUpdatedAt: time.Now()},
	}
	n, err := repo.UpsertIssues(context.Background(), issues, map[string]int64{"KNOWN": 1})

	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestUpsertIssues_Success(t *testing.T) {
	db, mock := newRepoDB(t)
	repo := NewRepository(db)

	mock.ExpectExec(`INSERT INTO issues`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	issues := []normalizer.DBIssue{
		{
			JiraIssueID:   "I1",
			JiraIssueKey:  "KEY-1",
			JiraProjectID: "P1",
			Summary:       "Test issue",
			Status:        "In Progress",
			StatusCategory: "In Progress",
			DelayStatus:   "GREEN",
			LastUpdatedAt: time.Now(),
		},
	}
	n, err := repo.UpsertIssues(context.Background(), issues, map[string]int64{"P1": 10})

	assert.NoError(t, err)
	assert.Equal(t, 1, n)
}

// --- GetProjectIDMap tests ---

func TestGetProjectIDMap_Empty(t *testing.T) {
	db, mock := newRepoDB(t)
	repo := NewRepository(db)

	mock.ExpectQuery(`SELECT id, jira_project_id FROM projects`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "jira_project_id"}))

	m, err := repo.GetProjectIDMap(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, m)
}

func TestGetProjectIDMap_WithRows(t *testing.T) {
	db, mock := newRepoDB(t)
	repo := NewRepository(db)

	mock.ExpectQuery(`SELECT id, jira_project_id FROM projects`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "jira_project_id"}).
				AddRow(1, "P1").
				AddRow(2, "P2"),
		)

	m, err := repo.GetProjectIDMap(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, map[string]int64{"P1": 1, "P2": 2}, m)
}

// --- StartSyncLog tests ---

func TestStartSyncLog_Success(t *testing.T) {
	db, mock := newRepoDB(t)
	repo := NewRepository(db)

	mock.ExpectQuery(`INSERT INTO sync_logs`).
		WithArgs("FULL").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	id, err := repo.StartSyncLog(context.Background(), "FULL")

	assert.NoError(t, err)
	assert.Equal(t, int64(42), id)
}

// --- GetLastSuccessfulSyncTime tests ---

func TestGetLastSuccessfulSyncTime_NoRows(t *testing.T) {
	db, mock := newRepoDB(t)
	repo := NewRepository(db)

	mock.ExpectQuery(`SELECT executed_at FROM sync_logs`).
		WithArgs("FULL").
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetLastSuccessfulSyncTime(context.Background(), "FULL")

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGetLastSuccessfulSyncTime_WithRow(t *testing.T) {
	db, mock := newRepoDB(t)
	repo := NewRepository(db)

	now := time.Now().Truncate(time.Second)
	mock.ExpectQuery(`SELECT executed_at FROM sync_logs`).
		WithArgs("FULL").
		WillReturnRows(sqlmock.NewRows([]string{"executed_at"}).AddRow(now))

	result, err := repo.GetLastSuccessfulSyncTime(context.Background(), "FULL")

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, now, *result)
}

// --- FinishSyncLog tests ---

func TestFinishSyncLog_Success(t *testing.T) {
	db, mock := newRepoDB(t)
	repo := NewRepository(db)

	mock.ExpectExec(`UPDATE sync_logs SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.FinishSyncLog(context.Background(), 1, "SUCCESS", 5, 100, "")

	assert.NoError(t, err)
}

func TestFinishSyncLog_WithErrorMessage(t *testing.T) {
	db, mock := newRepoDB(t)
	repo := NewRepository(db)

	mock.ExpectExec(`UPDATE sync_logs SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.FinishSyncLog(context.Background(), 2, "FAILED", 0, 0, "connection refused")

	assert.NoError(t, err)
}
