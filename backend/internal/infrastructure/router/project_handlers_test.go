package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// projectCols lists the columns returned by the project aggregate query.
var projectCols = []string{
	"id", "jira_project_id", "key", "name",
	"lead_account_id", "lead_email", "organization_id", "is_active",
	"created_at", "updated_at",
	"red_count", "yellow_count", "green_count", "open_count", "total_count",
}

// --- listProjectsHandlerWithDB tests ---

func TestListProjectsHandler_EmptyResult(t *testing.T) {
	db, mock := newTestDB(t)
	handler := listProjectsHandlerWithDB(db)

	// COUNT query for total
	mock.ExpectQuery(`SELECT COUNT`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// Data query returns empty
	mock.ExpectQuery(`SELECT`).
		WillReturnRows(sqlmock.NewRows(projectCols))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/projects", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp ProjectListResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Empty(t, resp.Data)
	assert.Equal(t, 0, resp.Pagination.Total)
}

// --- getProjectHandlerWithDB tests ---

func TestGetProjectHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := getProjectHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/projects/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "invalid project id", resp["error"])
}

func TestGetProjectHandler_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows(projectCols))

	handler := getProjectHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/projects/999", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetProjectHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	now := time.Now()

	rows := sqlmock.NewRows(projectCols).
		AddRow(1, "JIRA-1", "PROJ", "テストプロジェクト", nil, nil, nil, true, now, now, 2, 1, 5, 3, 8)
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	handler := getProjectHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/projects/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp ProjectRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, "PROJ", resp.Key)
	assert.Equal(t, "RED", resp.DelayStatus)
	assert.Equal(t, 2, resp.RedCount)
}

// --- updateProjectHandlerWithDB tests ---

func TestUpdateProjectHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := updateProjectHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/projects/abc", bytes.NewBufferString(`{"is_active":true}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "invalid project id", resp["error"])
}

func TestUpdateProjectHandler_MissingIsActive(t *testing.T) {
	db, _ := newTestDB(t)
	handler := updateProjectHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/projects/1", bytes.NewBufferString(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "is_active is required", resp["error"])
}

func TestUpdateProjectHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := updateProjectHandlerWithDB(db)

	mock.ExpectExec(`UPDATE projects SET is_active`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/projects/1", bytes.NewBufferString(`{"is_active":false}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "project updated", resp["message"])
}

func TestListProjectsHandler_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT`).WillReturnError(sqlmock.ErrCancelled)

	handler := listProjectsHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/projects", nil)

	handler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
