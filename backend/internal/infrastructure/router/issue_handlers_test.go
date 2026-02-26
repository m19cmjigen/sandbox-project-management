package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIssueHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := getIssueHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/issues/notanid", nil)
	c.Params = gin.Params{{Key: "id", Value: "notanid"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetIssueHandler_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := getIssueHandlerWithDB(db)

	// Mock: no rows returned
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(sqlmock.NewRows([]string{}))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/issues/999", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "issue not found", resp["error"])
}

func TestListIssuesHandler_DefaultParams(t *testing.T) {
	db, mock := newTestDB(t)
	handler := listIssuesHandlerWithDB(db)

	// Mock: COUNT query
	mock.ExpectQuery(`SELECT COUNT\(\*\)`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// Mock: main data query
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "jira_issue_id", "jira_issue_key", "project_id",
			"project_key", "project_name", "summary", "status",
			"status_category", "due_date", "assignee_name", "assignee_account_id",
			"delay_status", "priority", "issue_type",
			"last_updated_at", "created_at", "updated_at",
		}))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/issues", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp IssueListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 0, resp.Pagination.Total)
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 25, resp.Pagination.PerPage)
}

func TestListProjectIssuesHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := listProjectIssuesHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/projects/abc/issues", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "invalid project id", resp["error"])
}

func TestListProjectIssuesHandler_EmptyResult(t *testing.T) {
	db, mock := newTestDB(t)
	handler := listProjectIssuesHandlerWithDB(db)

	// COUNT query
	mock.ExpectQuery(`SELECT COUNT\(\*\)`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// Data query
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(1), 25, 0).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "jira_issue_id", "jira_issue_key", "project_id",
			"project_key", "project_name", "summary", "status",
			"status_category", "due_date", "assignee_name", "assignee_account_id",
			"delay_status", "priority", "issue_type",
			"last_updated_at", "created_at", "updated_at",
		}))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/projects/1/issues", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp IssueListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 0, resp.Pagination.Total)
	assert.Len(t, resp.Data, 0)
}

func TestListIssuesHandler_PaginationBounds(t *testing.T) {
	db, mock := newTestDB(t)
	handler := listIssuesHandlerWithDB(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\)`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "jira_issue_id", "jira_issue_key", "project_id",
			"project_key", "project_name", "summary", "status",
			"status_category", "due_date", "assignee_name", "assignee_account_id",
			"delay_status", "priority", "issue_type",
			"last_updated_at", "created_at", "updated_at",
		}))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// per_page > 100 should be clamped to 25
	c.Request = httptest.NewRequest(http.MethodGet, "/issues?page=0&per_page=9999", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp IssueListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 1, resp.Pagination.Page)
	assert.Equal(t, 25, resp.Pagination.PerPage)
}

func TestListIssuesHandler_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT COUNT`).WillReturnError(sqlmock.ErrCancelled)

	handler := listIssuesHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/issues", nil)

	handler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestListProjectIssuesHandler_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT`).WillReturnError(sqlmock.ErrCancelled)

	handler := listProjectIssuesHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/projects/1/issues", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
