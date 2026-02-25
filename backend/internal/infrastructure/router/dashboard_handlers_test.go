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
