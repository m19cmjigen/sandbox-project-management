package router

import (
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

func TestOrgDelayStatus_Red(t *testing.T) {
	org := &OrganizationRow{RedProjects: 3, YellowProjects: 1, GreenProjects: 0}
	assert.Equal(t, "RED", orgDelayStatus(org))
}

func TestOrgDelayStatus_Yellow(t *testing.T) {
	org := &OrganizationRow{RedProjects: 0, YellowProjects: 2, GreenProjects: 1}
	assert.Equal(t, "YELLOW", orgDelayStatus(org))
}

func TestOrgDelayStatus_Green(t *testing.T) {
	org := &OrganizationRow{RedProjects: 0, YellowProjects: 0, GreenProjects: 5}
	assert.Equal(t, "GREEN", orgDelayStatus(org))
}

func TestOrgDelayStatus_AllZero(t *testing.T) {
	org := &OrganizationRow{RedProjects: 0, YellowProjects: 0, GreenProjects: 0}
	assert.Equal(t, "GREEN", orgDelayStatus(org))
}

// orgCols lists the columns returned by orgQuery SELECT.
var orgCols = []string{
	"id", "name", "parent_id", "path", "level",
	"created_at", "updated_at",
	"total_projects", "red_projects", "yellow_projects", "green_projects",
}

// --- listOrganizationsHandlerWithDB tests ---

func TestListOrganizationsHandler_EmptyResult(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows(orgCols))

	handler := listOrganizationsHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/organizations", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []OrganizationRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Empty(t, resp)
}

func TestListOrganizationsHandler_ReturnsList(t *testing.T) {
	db, mock := newTestDB(t)
	now := time.Now()
	parentID := int64(1)
	rows := sqlmock.NewRows(orgCols).
		AddRow(1, "開発本部", nil, "/", 0, now, now, 6, 2, 1, 3).
		AddRow(2, "第一開発部", &parentID, "/", 1, now, now, 3, 0, 1, 2)
	mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	handler := listOrganizationsHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/organizations", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []OrganizationRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp, 2)
	assert.Equal(t, "開発本部", resp[0].Name)
	assert.Equal(t, "RED", resp[0].DelayStatus)
}

// --- getOrganizationHandlerWithDB tests ---

func TestGetOrganizationHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := getOrganizationHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/organizations/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetOrganizationHandler_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(999)).
		WillReturnRows(sqlmock.NewRows(orgCols))

	handler := getOrganizationHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/organizations/999", nil)
	c.Params = gin.Params{{Key: "id", Value: "999"}}

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetOrganizationHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	now := time.Now()
	rows := sqlmock.NewRows(orgCols).AddRow(1, "開発本部", nil, "/", 0, now, now, 4, 1, 0, 3)
	mock.ExpectQuery(`SELECT`).WithArgs(int64(1)).WillReturnRows(rows)

	handler := getOrganizationHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/organizations/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp OrganizationRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, "RED", resp.DelayStatus)
}

// --- getChildOrganizationsHandlerWithDB tests ---

func TestGetChildOrganizationsHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := getChildOrganizationsHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/organizations/abc/children", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetChildOrganizationsHandler_EmptyResult(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT`).WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(orgCols))

	handler := getChildOrganizationsHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/organizations/1/children", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []OrganizationRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Empty(t, resp)
}
