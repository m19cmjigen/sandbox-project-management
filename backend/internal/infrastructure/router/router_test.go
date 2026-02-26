package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	logpkg "github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
)

// newTestLogger はミドルウェアテスト用の no-op ロガーを返す
func newTestLogger(t *testing.T) *logpkg.Logger {
	t.Helper()
	log, err := logpkg.New("debug", "json")
	require.NoError(t, err)
	return log
}

// --- SecurityHeadersMiddleware tests ---

func TestSecurityHeadersMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeadersMiddleware())
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
	assert.NotEmpty(t, w.Header().Get("Content-Security-Policy"))
}

// --- CORSMiddleware tests ---

func TestCORSMiddleware_Wildcard(t *testing.T) {
	// allowedOrigins が空の場合は * を返す
	r := gin.New()
	r.Use(CORSMiddleware(""))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	r.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
}

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	r := gin.New()
	r.Use(CORSMiddleware("https://example.com"))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	r.ServeHTTP(w, req)

	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_DisallowedOrigin(t *testing.T) {
	r := gin.New()
	r.Use(CORSMiddleware("https://allowed.com"))
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://evil.com")
	r.ServeHTTP(w, req)

	// 許可されていないオリジンには CORS ヘッダーを付与しない
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_Options(t *testing.T) {
	r := gin.New()
	r.Use(CORSMiddleware(""))
	r.OPTIONS("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set("Origin", "https://example.com")
	r.ServeHTTP(w, req)

	// OPTIONS プリフライトには 204 No Content を返す
	assert.Equal(t, http.StatusNoContent, w.Code)
}

// --- LoggerMiddleware tests ---

func TestLoggerMiddleware_PassThrough(t *testing.T) {
	log := newTestLogger(t)
	r := gin.New()
	r.Use(LoggerMiddleware(log))
	r.GET("/", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	// ミドルウェアがハンドラーへの処理を通過させることを確認
	assert.Equal(t, http.StatusOK, w.Code)
}

// --- healthCheckHandler tests ---

func TestHealthCheckHandler(t *testing.T) {
	db, _ := newTestDB(t)
	handler := healthCheckHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

// --- readinessCheckHandler tests ---

func TestReadinessCheckHandler_DBUp(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectPing()
	handler := readinessCheckHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ready", nil)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ready")
}

func TestReadinessCheckHandler_DBDown(t *testing.T) {
	// DBを閉じることで Ping が失敗するシナリオを再現する
	dbClosed, _ := newTestDB(t)
	dbClosed.Close()

	handler := readinessCheckHandler(dbClosed)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/ready", nil)

	handler(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), "unavailable")
}
