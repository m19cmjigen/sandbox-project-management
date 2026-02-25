package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// --- maskToken tests ---

func TestMaskToken_ShortToken(t *testing.T) {
	assert.Equal(t, "•••••", maskToken("abc"))
}

func TestMaskToken_LongToken(t *testing.T) {
	assert.Equal(t, "•••••1234", maskToken("supersecrettoken1234"))
}

// --- getJiraSettingsHandler tests ---

func TestGetJiraSettingsHandler_NotConfigured(t *testing.T) {
	db, mock := newTestDB(t)
	// DBにレコードがない場合、StructScan は sql.ErrNoRows を返す
	mock.ExpectQuery(`SELECT id, jira_url, email, api_token, created_at, updated_at`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "jira_url", "email", "api_token", "created_at", "updated_at"}))

	handler := getJiraSettingsHandler(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/settings/jira", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp jiraSettingsResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.False(t, resp.Configured)
}

func TestGetJiraSettingsHandler_Configured(t *testing.T) {
	db, mock := newTestDB(t)
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "jira_url", "email", "api_token", "created_at", "updated_at"}).
		AddRow(1, "https://example.atlassian.net", "user@example.com", "supersecrettoken1234", now, now)
	mock.ExpectQuery(`SELECT id, jira_url, email, api_token, created_at, updated_at`).
		WillReturnRows(rows)

	handler := getJiraSettingsHandler(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/settings/jira", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp jiraSettingsResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.True(t, resp.Configured)
	assert.Equal(t, "https://example.atlassian.net", resp.JiraURL)
	assert.Equal(t, "•••••1234", resp.APITokenMask)
}

// --- updateJiraSettingsHandler tests ---

func TestUpdateJiraSettingsHandler_MissingFields(t *testing.T) {
	db, _ := newTestDB(t)
	handler := updateJiraSettingsHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/settings/jira", bytes.NewBufferString(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateJiraSettingsHandler_InvalidURL(t *testing.T) {
	db, _ := newTestDB(t)
	handler := updateJiraSettingsHandler(db)

	body := bytes.NewBufferString(`{"jira_url":"not-a-url","email":"user@example.com","api_token":"token"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/settings/jira", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- triggerSyncHandler tests ---

func TestTriggerSyncHandler_AlreadyRunning(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM sync_logs WHERE status = 'RUNNING'`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	handler := triggerSyncHandler(db, zap.NewNop())
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/settings/jira/sync", nil)

	handler(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestTriggerSyncHandler_NotConfigured(t *testing.T) {
	db, mock := newTestDB(t)
	// running count = 0
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM sync_logs WHERE status = 'RUNNING'`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// jira_settings が未登録 → 空行を返す
	mock.ExpectQuery(`SELECT jira_url, email, api_token FROM jira_settings`).
		WillReturnRows(sqlmock.NewRows([]string{"jira_url", "email", "api_token"}))

	handler := triggerSyncHandler(db, zap.NewNop())
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/settings/jira/sync", nil)

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp["error"], "not configured")
}

// --- listSyncLogsHandler tests ---

func TestListSyncLogsHandler_ReturnsList(t *testing.T) {
	db, mock := newTestDB(t)
	now := time.Now()
	durSec := 10
	rows := sqlmock.NewRows([]string{
		"id", "sync_type", "executed_at", "completed_at", "status",
		"projects_synced", "issues_synced", "error_message", "duration_seconds",
	}).AddRow(1, "FULL", now, &now, "SUCCESS", 5, 100, nil, &durSec)
	mock.ExpectQuery(`SELECT id, sync_type, executed_at`).WillReturnRows(rows)

	handler := listSyncLogsHandler(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/sync-logs", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string][]syncLogRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp["data"], 1)
	assert.Equal(t, "FULL", resp["data"][0].SyncType)
}
