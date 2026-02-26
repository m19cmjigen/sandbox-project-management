package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// testRouter creates a gin router with the given middlewares and a simple 200 handler.
func testRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.GET("/", append(middlewares, func(c *gin.Context) {
		c.Status(http.StatusOK)
	})...)
	return r
}

// --- GetClaims tests ---

func TestGetClaims_NotSet(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	assert.Nil(t, GetClaims(c))
}

func TestGetClaims_Set(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	expected := &Claims{UserID: 42, Email: "test@example.com", Role: "admin"}
	c.Set(claimsKey, expected)

	got := GetClaims(c)
	require.NotNil(t, got)
	assert.Equal(t, expected.UserID, got.UserID)
	assert.Equal(t, expected.Email, got.Email)
	assert.Equal(t, expected.Role, got.Role)
}

// --- Middleware tests ---

func TestMiddleware_MissingHeader(t *testing.T) {
	tm := NewTokenManager("secret")
	router := testRouter(Middleware(tm))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "authorization header is required", resp["error"])
}

func TestMiddleware_BadFormat(t *testing.T) {
	tm := NewTokenManager("secret")
	router := testRouter(Middleware(tm))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "InvalidHeaderWithoutBearer")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp["error"], "Bearer")
}

func TestMiddleware_InvalidToken(t *testing.T) {
	tm := NewTokenManager("secret")
	router := testRouter(Middleware(tm))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer this.is.notvalid")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "invalid token", resp["error"])
}

func TestMiddleware_ExpiredToken(t *testing.T) {
	tm := NewTokenManager("secret")
	router := testRouter(Middleware(tm))

	past := time.Now().Add(-2 * time.Hour)
	claims := Claims{
		UserID: 1,
		Email:  "user@example.com",
		Role:   "viewer",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(past.Add(-1 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(past),
		},
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expired, _ := raw.SignedString(tm.secret)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+expired)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "token has expired", resp["error"])
}

func TestMiddleware_ValidToken(t *testing.T) {
	tm := NewTokenManager("secret")
	router := gin.New()

	// Capture claims set by middleware for assertion
	var gotClaims *Claims
	router.GET("/", Middleware(tm), func(c *gin.Context) {
		gotClaims = GetClaims(c)
		c.Status(http.StatusOK)
	})

	token, err := tm.GenerateAccessToken(7, "admin@example.com", "admin")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	require.NotNil(t, gotClaims)
	assert.Equal(t, int64(7), gotClaims.UserID)
	assert.Equal(t, "admin@example.com", gotClaims.Email)
	assert.Equal(t, "admin", gotClaims.Role)
}

// --- RequireRole tests ---

func TestRequireRole_NoClaims(t *testing.T) {
	// RequireRole without Middleware → no claims in context
	router := testRouter(RequireRole("admin"))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "unauthorized", resp["error"])
}

func TestRequireRole_WrongRole(t *testing.T) {
	tm := NewTokenManager("secret")
	router := testRouter(Middleware(tm), RequireRole("admin"))

	// Generate a viewer token
	token, err := tm.GenerateAccessToken(1, "viewer@example.com", "viewer")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "insufficient permissions", resp["error"])
}

func TestRequireRole_AllowedRole(t *testing.T) {
	tm := NewTokenManager("secret")
	router := testRouter(Middleware(tm), RequireRole("admin", "project_manager"))

	token, err := tm.GenerateAccessToken(1, "pm@example.com", "project_manager")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
