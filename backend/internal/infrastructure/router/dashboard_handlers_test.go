package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- getDashboardSummaryHandlerWithDB tests ---

func TestGetDashboardSummaryHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := getDashboardSummaryHandlerWithDB(db)

	// 1st query: project counts (CTE + COUNT)
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(sqlmock.NewRows([]string{"total", "red", "yellow", "green"}).
			AddRow(10, 3, 2, 5))

	// 2nd query: issue counts
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(sqlmock.NewRows([]string{"total", "red", "yellow", "green"}).
			AddRow(80, 15, 10, 55))

	// 3rd query: per-org stats
	orgStatsCols := []string{"id", "name", "parent_id", "level", "total_projects", "red_projects", "yellow_projects", "green_projects", "delay_status"}
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(sqlmock.NewRows(orgStatsCols).
			AddRow(1, "開発本部", nil, 0, 6, 2, 1, 3, "RED"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/summary", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp DashboardSummaryResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, 10, resp.TotalProjects)
	assert.Equal(t, 3, resp.RedProjects)
	assert.Len(t, resp.Organizations, 1)
	assert.Equal(t, "開発本部", resp.Organizations[0].Name)
}

// --- getOrganizationSummaryHandlerWithDB tests ---

func TestGetOrganizationSummaryHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := getOrganizationSummaryHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/organizations/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetOrganizationSummaryHandler_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := getOrganizationSummaryHandlerWithDB(db)

	orgStatsCols := []string{"id", "name", "parent_id", "level", "total_projects", "red_projects", "yellow_projects", "green_projects", "delay_status"}
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows(orgStatsCols))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/organizations/999", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetOrganizationSummaryHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := getOrganizationSummaryHandlerWithDB(db)

	now := time.Now()
	orgStatsCols := []string{"id", "name", "parent_id", "level", "total_projects", "red_projects", "yellow_projects", "green_projects", "delay_status"}
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(orgStatsCols).
			AddRow(1, "開発本部", nil, 0, 4, 1, 1, 2, "RED"))

	projectCols := []string{
		"id", "jira_project_id", "key", "name",
		"lead_account_id", "lead_email", "organization_id",
		"red_count", "yellow_count", "green_count", "open_count", "total_count",
		"created_at", "updated_at",
	}
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(projectCols).
			AddRow(1, "JIRA-1", "PROJ", "テストPJ", nil, nil, 1, 1, 0, 3, 4, 4, now, now))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/organizations/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	org := resp["organization"].(map[string]interface{})
	assert.Equal(t, "開発本部", org["name"])
}

// --- getProjectSummaryHandlerWithDB tests ---

func TestGetProjectSummaryHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := getProjectSummaryHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/projects/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "invalid project id", resp["error"])
}

func TestGetProjectSummaryHandler_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := getProjectSummaryHandlerWithDB(db)

	// Project query returns no rows
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/projects/999", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "project not found", resp["error"])
}

func TestGetProjectSummaryHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := getProjectSummaryHandlerWithDB(db)

	now := time.Now()

	// Project query
	projectCols := []string{
		"id", "jira_project_id", "key", "name",
		"lead_account_id", "lead_email", "organization_id", "is_active",
		"created_at", "updated_at",
		"red_count", "yellow_count", "green_count", "open_count", "total_count",
	}
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(projectCols).
			AddRow(1, "JIRA-1", "PROJ", "Test Project", nil, nil, nil, true, now, now, 2, 1, 5, 3, 8))

	// Delayed issues query
	issueCols := []string{
		"id", "jira_issue_id", "jira_issue_key", "project_id",
		"project_key", "project_name", "summary", "status",
		"status_category", "due_date", "assignee_name", "assignee_account_id",
		"delay_status", "priority", "issue_type",
		"last_updated_at", "created_at", "updated_at",
	}
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(issueCols).
			AddRow(101, "JIRA-I-1", "PROJ-1", 1, "PROJ", "Test Project",
				"Fix bug", "In Progress", "In Progress",
				"2026-02-20", nil, nil, "RED", nil, nil, now, now, now))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/projects/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp ProjectSummaryResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(1), resp.Project.ID)
	assert.Equal(t, "RED", resp.Project.DelayStatus)
	assert.Equal(t, 2, resp.Summary.RedCount)
	assert.Equal(t, 1, resp.Summary.YellowCount)
	assert.Len(t, resp.DelayedIssues, 1)
	assert.Equal(t, "RED", resp.DelayedIssues[0].DelayStatus)
}

func TestGetDashboardSummaryHandler_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`WITH project_stats`).WillReturnError(sqlmock.ErrCancelled)

	handler := getDashboardSummaryHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/dashboard/summary", nil)

	handler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
