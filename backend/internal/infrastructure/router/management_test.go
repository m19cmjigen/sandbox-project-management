package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// newTestDB returns a sqlx.DB backed by go-sqlmock.
func newTestDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	t.Cleanup(func() { sqlxDB.Close() })
	return sqlxDB, mock
}

// --- createOrganizationHandlerWithDB tests ---

func TestCreateOrganizationHandler_MissingName(t *testing.T) {
	db, _ := newTestDB(t)
	handler := createOrganizationHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBufferString(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp["error"], "name")
}

func TestCreateOrganizationHandler_EmptyName(t *testing.T) {
	db, _ := newTestDB(t)
	handler := createOrganizationHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(`{"name":"   "}`)
	c.Request = httptest.NewRequest(http.MethodPost, "/organizations", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- updateOrganizationHandlerWithDB tests ---

func TestUpdateOrganizationHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := updateOrganizationHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/organizations/abc", bytes.NewBufferString(`{"name":"X"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOrganizationHandler_MissingName(t *testing.T) {
	db, _ := newTestDB(t)
	handler := updateOrganizationHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/organizations/1", bytes.NewBufferString(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- deleteOrganizationHandlerWithDB tests ---

func TestDeleteOrganizationHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := deleteOrganizationHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/organizations/xyz", nil)
	c.Params = gin.Params{{Key: "id", Value: "xyz"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteOrganizationHandler_HasChildren(t *testing.T) {
	db, mock := newTestDB(t)
	handler := deleteOrganizationHandlerWithDB(db)

	// Mock: COUNT child organizations returns 2
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM organizations WHERE parent_id`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/organizations/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp["error"], "child")
}

func TestDeleteOrganizationHandler_HasProjects(t *testing.T) {
	db, mock := newTestDB(t)
	handler := deleteOrganizationHandlerWithDB(db)

	// Mock: no children
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM organizations WHERE parent_id`).
		WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	// Mock: has projects
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM projects WHERE organization_id`).
		WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/organizations/5", nil)
	c.Params = gin.Params{{Key: "id", Value: "5"}}

	handler(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp["error"], "project")
}

// --- assignProjectToOrganizationHandlerWithDB tests ---

func TestAssignProjectHandler_InvalidProjectID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := assignProjectToOrganizationHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/projects/bad/organization", bytes.NewBufferString(`{"organization_id":1}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "bad"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
