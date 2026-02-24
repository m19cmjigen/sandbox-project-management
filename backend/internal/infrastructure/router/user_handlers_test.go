package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/auth"
)

// setAdminClaims injects admin Claims into the Gin context to simulate authenticated request.
func setAdminClaims(c *gin.Context, userID int64) {
	c.Set("claims", &auth.Claims{UserID: userID, Email: "admin@example.com", Role: "admin"})
}

// --- listUsersHandlerWithDB tests ---

func TestListUsersHandler_ReturnsList(t *testing.T) {
	db, mock := newTestDB(t)

	rows := sqlmock.NewRows([]string{"id", "email", "role", "is_active"}).
		AddRow(1, "admin@example.com", "admin", true).
		AddRow(2, "viewer@example.com", "viewer", true)
	mock.ExpectQuery(`SELECT id, email, role, is_active FROM users`).WillReturnRows(rows)

	handler := listUsersHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/users", nil)
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string][]userListItem
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp["data"], 2)
}

// --- createUserHandlerWithDB tests ---

func TestCreateUserHandler_MissingFields(t *testing.T) {
	db, _ := newTestDB(t)
	handler := createUserHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUserHandler_InvalidRole(t *testing.T) {
	db, _ := newTestDB(t)
	handler := createUserHandlerWithDB(db)

	body := bytes.NewBufferString(`{"email":"test@example.com","password":"Password1!","role":"superadmin"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUserHandler_ShortPassword(t *testing.T) {
	db, _ := newTestDB(t)
	handler := createUserHandlerWithDB(db)

	body := bytes.NewBufferString(`{"email":"test@example.com","password":"short","role":"viewer"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- updateUserHandlerWithDB tests ---

func TestUpdateUserHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := updateUserHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/abc", bytes.NewBufferString(`{"role":"viewer"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUserHandler_InvalidRole(t *testing.T) {
	db, _ := newTestDB(t)
	handler := updateUserHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	invalidRole := "superadmin"
	body, _ := json.Marshal(updateUserRequest{Role: &invalidRole})
	c.Request = httptest.NewRequest(http.MethodPut, "/users/2", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "2"}}

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- deleteUserHandlerWithDB tests ---

func TestDeleteUserHandler_SelfDeletion(t *testing.T) {
	db, _ := newTestDB(t)
	handler := deleteUserHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	// adminユーザー自身を削除しようとするケース
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "cannot delete yourself", resp["error"])
}

func TestDeleteUserHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := deleteUserHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteUserHandler_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := deleteUserHandlerWithDB(db)

	// 該当ユーザーが存在しない場合は rowsAffected = 0
	mock.ExpectExec(`DELETE FROM users WHERE id`).
		WithArgs(int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/users/99", nil)
	c.Params = gin.Params{{Key: "id", Value: "99"}}
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
