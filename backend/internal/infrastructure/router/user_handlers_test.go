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

func TestCreateUserHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := createUserHandlerWithDB(db)

	mock.ExpectQuery(`INSERT INTO users`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role", "is_active"}).
			AddRow(3, "new@example.com", "viewer", true))

	body := bytes.NewBufferString(`{"email":"new@example.com","password":"Password1!","role":"viewer"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp userListItem
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "new@example.com", resp.Email)
	assert.Equal(t, "viewer", resp.Role)
}

func TestUpdateUserHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := updateUserHandlerWithDB(db)

	mock.ExpectExec(`UPDATE users`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(`SELECT id, email, role, is_active FROM users`).
		WithArgs(int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role", "is_active"}).
			AddRow(2, "viewer@example.com", "project_manager", true))

	newRole := "project_manager"
	body, _ := json.Marshal(updateUserRequest{Role: &newRole})
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/2", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "2"}}

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp userListItem
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "project_manager", resp.Role)
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

func TestListUsersHandler_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT id, email, role, is_active FROM users`).
		WillReturnError(sqlmock.ErrCancelled)

	handler := listUsersHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/users", nil)
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteUserHandler_DeleteSuccess(t *testing.T) {
	db, mock := newTestDB(t)
	handler := deleteUserHandlerWithDB(db)

	mock.ExpectExec(`DELETE FROM users WHERE id`).
		WithArgs(int64(2)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/users/2", nil)
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// --- changePasswordHandlerWithDB tests ---

func TestChangePasswordHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	handler := changePasswordHandlerWithDB(db)

	// パスワードハッシュ更新が1行に影響する
	mock.ExpectExec(`UPDATE users SET password_hash`).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := bytes.NewBufferString(`{"new_password":"newpass123"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/2/password", body)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "password updated", resp["message"])
}

func TestChangePasswordHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := changePasswordHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/abc/password", bytes.NewBufferString(`{"new_password":"newpass123"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "invalid user id", resp["error"])
}

func TestChangePasswordHandler_TooShort(t *testing.T) {
	db, _ := newTestDB(t)
	handler := changePasswordHandlerWithDB(db)

	body := bytes.NewBufferString(`{"new_password":"short"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/2/password", body)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChangePasswordHandler_UserNotFound(t *testing.T) {
	db, mock := newTestDB(t)
	handler := changePasswordHandlerWithDB(db)

	// 該当ユーザーが存在しない場合は rowsAffected = 0
	mock.ExpectExec(`UPDATE users SET password_hash`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	body := bytes.NewBufferString(`{"new_password":"newpass123"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/99/password", body)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "99"}}
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "user not found", resp["error"])
}

func TestChangePasswordHandler_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	handler := changePasswordHandlerWithDB(db)

	mock.ExpectExec(`UPDATE users SET password_hash`).
		WillReturnError(sqlmock.ErrCancelled)

	body := bytes.NewBufferString(`{"new_password":"newpass123"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/2/password", body)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: "2"}}
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "failed to update password", resp["error"])
}
