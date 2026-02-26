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

// --- createOrganizationHandlerWithDB (success paths) tests ---

func TestCreateOrganizationHandler_ParentNotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := createOrganizationHandlerWithDB(db)

	mock.ExpectQuery(`SELECT path, level FROM organizations WHERE id`).
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"path", "level"}))

	parentID := int64(99)
	body, _ := json.Marshal(createOrgRequest{Name: "子部署", ParentID: &parentID})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "parent organization not found", resp["error"])
}

func TestCreateOrganizationHandler_MaxDepthExceeded(t *testing.T) {
	db, mock := newTestDB(t)
	handler := createOrganizationHandlerWithDB(db)

	mock.ExpectQuery(`SELECT path, level FROM organizations WHERE id`).
		WithArgs(int64(3)).
		WillReturnRows(sqlmock.NewRows([]string{"path", "level"}).AddRow("/1/2/3/", 2))

	parentID := int64(3)
	body, _ := json.Marshal(createOrgRequest{Name: "深すぎる部署", ParentID: &parentID})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp["error"], "depth")
}

func TestCreateOrganizationHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := createOrganizationHandlerWithDB(db)
	now := time.Now()

	mock.ExpectQuery(`INSERT INTO organizations`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))
	mock.ExpectExec(`UPDATE organizations SET path`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows(orgCols).AddRow(10, "新部署", nil, "/10/", 0, now, now, 0, 0, 0, 0))

	body, _ := json.Marshal(createOrgRequest{Name: "新部署"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp OrganizationRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(10), resp.ID)
	assert.Equal(t, "新部署", resp.Name)
}

func TestCreateOrganizationHandler_SuccessWithParent(t *testing.T) {
	db, mock := newTestDB(t)
	handler := createOrganizationHandlerWithDB(db)
	now := time.Now()

	// Fetch parent path/level
	mock.ExpectQuery(`SELECT path, level FROM organizations WHERE id`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"path", "level"}).AddRow("/1/", 0))
	// INSERT RETURNING id
	mock.ExpectQuery(`INSERT INTO organizations`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))
	// UPDATE path
	mock.ExpectExec(`UPDATE organizations SET path`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// SELECT org after creation
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(5)).
		WillReturnRows(sqlmock.NewRows(orgCols).AddRow(5, "子部署", int64(1), "/1/5/", 1, now, now, 0, 0, 0, 0))

	parentID := int64(1)
	body, _ := json.Marshal(createOrgRequest{Name: "子部署", ParentID: &parentID})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/organizations", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp OrganizationRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(5), resp.ID)
	assert.Equal(t, int64(1), *resp.ParentID)
}

// --- updateOrganizationHandlerWithDB (success paths) tests ---

func TestUpdateOrganizationHandler_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := updateOrganizationHandlerWithDB(db)

	mock.ExpectExec(`UPDATE organizations SET name`).
		WithArgs("NewName", int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	body, _ := json.Marshal(updateOrgRequest{Name: "NewName"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/organizations/99", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "99"}}

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateOrganizationHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := updateOrganizationHandlerWithDB(db)
	now := time.Now()

	mock.ExpectExec(`UPDATE organizations SET name`).
		WithArgs("更新後", int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(orgCols).AddRow(1, "更新後", nil, "/1/", 0, now, now, 3, 0, 0, 3))

	body, _ := json.Marshal(updateOrgRequest{Name: "更新後"})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/organizations/1", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp OrganizationRow
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "更新後", resp.Name)
}

// --- deleteOrganizationHandlerWithDB (success paths) tests ---

func TestDeleteOrganizationHandler_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := deleteOrganizationHandlerWithDB(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM organizations WHERE parent_id`).
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM projects WHERE organization_id`).
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec(`DELETE FROM organizations WHERE id`).
		WithArgs(int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/organizations/99", nil)
	c.Params = gin.Params{{Key: "id", Value: "99"}}

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteOrganizationHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := deleteOrganizationHandlerWithDB(db)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM organizations WHERE parent_id`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM projects WHERE organization_id`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec(`DELETE FROM organizations WHERE id`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/organizations/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "organization deleted", resp["message"])
}

// --- assignProjectToOrganizationHandlerWithDB (success paths) tests ---

func TestAssignProjectHandler_OrgNotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := assignProjectToOrganizationHandlerWithDB(db)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	orgID := int64(99)
	body, _ := json.Marshal(assignProjectOrgRequest{OrganizationID: &orgID})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/projects/1/organization", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "organization not found", resp["error"])
}

func TestAssignProjectHandler_ProjectNotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := assignProjectToOrganizationHandlerWithDB(db)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectExec(`UPDATE projects SET organization_id`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	orgID := int64(1)
	body, _ := json.Marshal(assignProjectOrgRequest{OrganizationID: &orgID})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/projects/99/organization", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "99"}}

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestAssignProjectHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := assignProjectToOrganizationHandlerWithDB(db)

	mock.ExpectQuery(`SELECT EXISTS`).
		WithArgs(int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectExec(`UPDATE projects SET organization_id`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	orgID := int64(2)
	body, _ := json.Marshal(assignProjectOrgRequest{OrganizationID: &orgID})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/projects/1/organization", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "project assigned successfully", resp["message"])
}

func TestAssignProjectHandler_SuccessNullOrg(t *testing.T) {
	db, mock := newTestDB(t)
	handler := assignProjectToOrganizationHandlerWithDB(db)

	// organization_id = nil → EXISTS check is skipped
	mock.ExpectExec(`UPDATE projects SET organization_id`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := bytes.NewBufferString(`{"organization_id":null}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/projects/1/organization", body)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
