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

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/auth"
)

// setTestClaims injects auth.Claims into the Gin context to simulate an authenticated request.
func setTestClaims(c *gin.Context, userID int64) {
	claims := &auth.Claims{UserID: userID}
	c.Set("claims", claims)
}

// notificationCols is the ordered list of columns returned by SELECT on notifications.
var notificationCols = []string{
	"id", "user_id", "type", "title", "body", "is_read", "related_log_id", "created_at",
}

// --- listNotificationsHandlerWithDB tests ---

func TestListNotificationsHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	now := time.Now()
	logID := int64(10)

	rows := sqlmock.NewRows(notificationCols).
		AddRow(1, int64(1), "SYNC_COMPLETED", "Jira sync completed", "All synced.", false, &logID, now).
		AddRow(2, int64(1), "SYNC_FAILED", "Jira sync failed", "An error occurred.", true, nil, now)
	mock.ExpectQuery(`SELECT id, user_id, type, title, body, is_read, related_log_id, created_at`).
		WillReturnRows(rows)

	handler := listNotificationsHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/notifications", nil)
	setTestClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp struct {
		Data        []notificationRow `json:"data"`
		UnreadCount int               `json:"unread_count"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp.Data, 2)
	assert.Equal(t, 1, resp.UnreadCount)
}

func TestListNotificationsHandler_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectQuery(`SELECT id, user_id, type, title, body, is_read, related_log_id, created_at`).
		WillReturnError(sqlmock.ErrCancelled)

	handler := listNotificationsHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/notifications", nil)
	setTestClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestListNotificationsHandler_Unauthorized(t *testing.T) {
	db, _ := newTestDB(t)
	handler := listNotificationsHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/notifications", nil)
	// No claims injected

	handler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// --- readNotificationHandlerWithDB tests ---

func TestReadNotificationHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectExec(`UPDATE notifications SET is_read = TRUE`).
		WithArgs(int64(5), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	handler := readNotificationHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/notifications/5/read", nil)
	c.Params = gin.Params{{Key: "id", Value: "5"}}
	setTestClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "notification marked as read", resp["message"])
}

func TestReadNotificationHandler_NotFound(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectExec(`UPDATE notifications SET is_read = TRUE`).
		WithArgs(int64(99), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	handler := readNotificationHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/notifications/99/read", nil)
	c.Params = gin.Params{{Key: "id", Value: "99"}}
	setTestClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestReadNotificationHandler_InvalidID(t *testing.T) {
	db, _ := newTestDB(t)
	handler := readNotificationHandlerWithDB(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/notifications/abc/read", nil)
	c.Params = gin.Params{{Key: "id", Value: "abc"}}
	setTestClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- readAllNotificationsHandlerWithDB tests ---

func TestReadAllNotificationsHandler_Success(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectExec(`UPDATE notifications SET is_read = TRUE`).
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(3, 3))

	handler := readAllNotificationsHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/notifications/read-all", nil)
	setTestClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "all notifications marked as read", resp["message"])
}

func TestReadAllNotificationsHandler_DBError(t *testing.T) {
	db, mock := newTestDB(t)
	mock.ExpectExec(`UPDATE notifications SET is_read = TRUE`).
		WithArgs(int64(1)).
		WillReturnError(sqlmock.ErrCancelled)

	handler := readAllNotificationsHandlerWithDB(db)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/notifications/read-all", nil)
	setTestClaims(c, 1)

	handler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
