package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/auth"
)

func newTestTokenManager() *auth.TokenManager {
	return auth.NewTokenManager("test-jwt-secret")
}

// --- loginHandler tests ---

func TestLoginHandler_MissingFields(t *testing.T) {
	db, _ := newTestDB(t)
	tm := newTestTokenManager()
	handler := loginHandler(db, tm)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginHandler_ShortPassword(t *testing.T) {
	db, _ := newTestDB(t)
	tm := newTestTokenManager()
	handler := loginHandler(db, tm)

	body := bytes.NewBufferString(`{"email":"user@example.com","password":"short"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	db, mock := newTestDB(t)
	tm := newTestTokenManager()
	handler := loginHandler(db, tm)

	mock.ExpectQuery(`SELECT id, email, password_hash, role, is_active FROM users`).
		WithArgs("notfound@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash", "role", "is_active"}))

	body := bytes.NewBufferString(`{"email":"notfound@example.com","password":"Password1!"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginHandler_DisabledAccount(t *testing.T) {
	db, mock := newTestDB(t)
	tm := newTestTokenManager()
	handler := loginHandler(db, tm)

	// bcrypt最小コストで高速なハッシュを生成（テスト専用）
	hash, _ := bcrypt.GenerateFromPassword([]byte("Password1!"), bcrypt.MinCost)
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role", "is_active"}).
		AddRow(1, "disabled@example.com", string(hash), "viewer", false)
	mock.ExpectQuery(`SELECT id, email, password_hash, role, is_active FROM users`).
		WithArgs("disabled@example.com").
		WillReturnRows(rows)

	body := bytes.NewBufferString(`{"email":"disabled@example.com","password":"Password1!"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "account is disabled", resp["error"])
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	db, mock := newTestDB(t)
	tm := newTestTokenManager()
	handler := loginHandler(db, tm)

	hash, _ := bcrypt.GenerateFromPassword([]byte("CorrectPass1!"), bcrypt.MinCost)
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role", "is_active"}).
		AddRow(1, "admin@example.com", string(hash), "admin", true)
	mock.ExpectQuery(`SELECT id, email, password_hash, role, is_active FROM users`).
		WithArgs("admin@example.com").
		WillReturnRows(rows)

	body := bytes.NewBufferString(`{"email":"admin@example.com","password":"WrongPass1!"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	tm := newTestTokenManager()
	handler := loginHandler(db, tm)

	hash, _ := bcrypt.GenerateFromPassword([]byte("Password1!"), bcrypt.MinCost)
	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role", "is_active"}).
		AddRow(1, "admin@example.com", string(hash), "admin", true)
	mock.ExpectQuery(`SELECT id, email, password_hash, role, is_active FROM users`).
		WithArgs("admin@example.com").
		WillReturnRows(rows)

	body := bytes.NewBufferString(`{"email":"admin@example.com","password":"Password1!"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", body)
	c.Request.Header.Set("Content-Type", "application/json")

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp loginResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, "Bearer", resp.TokenType)
	assert.Equal(t, "admin@example.com", resp.User.Email)
	assert.Equal(t, "admin", resp.User.Role)
}

// --- meHandler tests ---

func TestMeHandler_NoClaims(t *testing.T) {
	handler := meHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	// Claims を設定しない

	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMeHandler_Success(t *testing.T) {
	handler := meHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	setAdminClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp userInfo
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(1), resp.ID)
	assert.Equal(t, "admin@example.com", resp.Email)
	assert.Equal(t, "admin", resp.Role)
}
