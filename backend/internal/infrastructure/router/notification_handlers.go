package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/auth"
)

// notificationRow maps to the notifications table.
type notificationRow struct {
	ID           int64     `db:"id"             json:"id"`
	UserID       int64     `db:"user_id"         json:"-"`
	Type         string    `db:"type"            json:"type"`
	Title        string    `db:"title"           json:"title"`
	Body         string    `db:"body"            json:"body"`
	IsRead       bool      `db:"is_read"         json:"is_read"`
	RelatedLogID *int64    `db:"related_log_id"  json:"related_log_id"`
	CreatedAt    time.Time `db:"created_at"      json:"created_at"`
}

// listNotificationsHandlerWithDB handles GET /api/v1/notifications.
// Returns the authenticated user's notifications (unread first, up to 50).
func listNotificationsHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c)
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		rows := make([]notificationRow, 0)
		err := db.Select(&rows, `
			SELECT id, user_id, type, title, body, is_read, related_log_id, created_at
			FROM notifications
			WHERE user_id = $1
			ORDER BY is_read ASC, created_at DESC
			LIMIT 50`,
			claims.UserID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notifications"})
			return
		}

		unreadCount := 0
		for _, r := range rows {
			if !r.IsRead {
				unreadCount++
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"data":         rows,
			"unread_count": unreadCount,
		})
	}
}

// readNotificationHandlerWithDB handles PUT /api/v1/notifications/:id/read.
// Marks a specific notification as read for the authenticated user.
func readNotificationHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c)
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
			return
		}

		result, err := db.Exec(`
			UPDATE notifications SET is_read = TRUE
			WHERE id = $1 AND user_id = $2`,
			id, claims.UserID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification"})
			return
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
	}
}

// readAllNotificationsHandlerWithDB handles PUT /api/v1/notifications/read-all.
// Marks all notifications as read for the authenticated user.
func readAllNotificationsHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := auth.GetClaims(c)
		if claims == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		_, err := db.Exec(`
			UPDATE notifications SET is_read = TRUE
			WHERE user_id = $1 AND is_read = FALSE`,
			claims.UserID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notifications"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
	}
}

// broadcastSyncNotification creates a notification for every active user after a sync run.
// Errors are logged but do not panic, as this runs inside a goroutine.
func broadcastSyncNotification(db *sqlx.DB, log *zap.Logger, syncErr error) {
	notifType := "SYNC_COMPLETED"
	title := "Jira sync completed"
	body := "All projects and issues have been synchronized successfully."

	if syncErr != nil {
		notifType = "SYNC_FAILED"
		title = "Jira sync failed"
		body = "An error occurred during synchronization: " + syncErr.Error()
	}

	// 最後の同期ログIDを取得（通知に関連付けるため）
	var logID *int64
	var lastID int64
	if err := db.QueryRowx(`SELECT id FROM sync_logs ORDER BY executed_at DESC LIMIT 1`).Scan(&lastID); err == nil {
		logID = &lastID
	}

	// アクティブなユーザーID一覧を取得
	var userIDs []int64
	if err := db.Select(&userIDs, `SELECT id FROM users WHERE is_active = TRUE`); err != nil {
		log.Error("broadcastSyncNotification: failed to fetch user ids", zap.Error(err))
		return
	}

	for _, uid := range userIDs {
		_, err := db.Exec(`
			INSERT INTO notifications (user_id, type, title, body, related_log_id)
			VALUES ($1, $2, $3, $4, $5)`,
			uid, notifType, title, body, logID,
		)
		if err != nil {
			log.Error("broadcastSyncNotification: failed to insert notification",
				zap.Int64("user_id", uid),
				zap.Error(err),
			)
		}
	}
}
