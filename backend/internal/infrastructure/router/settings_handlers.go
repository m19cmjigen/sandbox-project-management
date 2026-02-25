package router

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/batch"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/jiraclient"
)

// jiraSettingsRow represents a row in the jira_settings table.
type jiraSettingsRow struct {
	ID        int64     `db:"id"         json:"id"`
	JiraURL   string    `db:"jira_url"   json:"jira_url"`
	Email     string    `db:"email"      json:"email"`
	APIToken  string    `db:"api_token"  json:"-"` // never serialised
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// jiraSettingsResponse is the public shape of Jira settings (token is masked).
type jiraSettingsResponse struct {
	ID            int64     `json:"id"`
	JiraURL       string    `json:"jira_url"`
	Email         string    `json:"email"`
	APITokenMask  string    `json:"api_token_mask"` // "•••••<last4>"
	Configured    bool      `json:"configured"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type updateJiraSettingsRequest struct {
	JiraURL  string `json:"jira_url"  binding:"required,url"`
	Email    string `json:"email"     binding:"required,email"`
	APIToken string `json:"api_token" binding:"required"`
}

// syncLogRow maps to the sync_logs table.
type syncLogRow struct {
	ID             int64      `db:"id"              json:"id"`
	SyncType       string     `db:"sync_type"       json:"sync_type"`
	ExecutedAt     time.Time  `db:"executed_at"     json:"executed_at"`
	CompletedAt    *time.Time `db:"completed_at"    json:"completed_at"`
	Status         string     `db:"status"          json:"status"`
	ProjectsSynced int        `db:"projects_synced" json:"projects_synced"`
	IssuesSynced   int        `db:"issues_synced"   json:"issues_synced"`
	ErrorMessage   *string    `db:"error_message"   json:"error_message"`
	DurationSec    *int       `db:"duration_seconds" json:"duration_seconds"`
}

// maskToken returns "•••••<last4>" when the token is long enough, otherwise "•••••".
func maskToken(token string) string {
	if len(token) <= 4 {
		return "•••••"
	}
	return "•••••" + token[len(token)-4:]
}

// getJiraSettingsHandler handles GET /api/v1/settings/jira.
func getJiraSettingsHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var row jiraSettingsRow
		err := db.QueryRowx(`SELECT id, jira_url, email, api_token, created_at, updated_at
		                     FROM jira_settings ORDER BY id LIMIT 1`).StructScan(&row)
		if err != nil {
			// 設定が未登録の場合は空レスポンスを返す
			c.JSON(http.StatusOK, jiraSettingsResponse{Configured: false})
			return
		}
		c.JSON(http.StatusOK, jiraSettingsResponse{
			ID:           row.ID,
			JiraURL:      row.JiraURL,
			Email:        row.Email,
			APITokenMask: maskToken(row.APIToken),
			Configured:   true,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
		})
	}
}

// updateJiraSettingsHandler handles PUT /api/v1/settings/jira.
// Upserts a single settings record (only one row is maintained).
func updateJiraSettingsHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateJiraSettingsRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// JiraURLの末尾スラッシュを除去して統一
		req.JiraURL = strings.TrimRight(req.JiraURL, "/")

		var id int64
		err := db.QueryRowx(`
			INSERT INTO jira_settings (jira_url, email, api_token)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
			RETURNING id`,
			req.JiraURL, req.Email, req.APIToken,
		).Scan(&id)

		if err != nil || id == 0 {
			// 既存レコードを更新
			_, err = db.Exec(`
				UPDATE jira_settings
				SET jira_url = $1, email = $2, api_token = $3, updated_at = CURRENT_TIMESTAMP
				WHERE id = (SELECT id FROM jira_settings ORDER BY id LIMIT 1)`,
				req.JiraURL, req.Email, req.APIToken,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save settings"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Jira settings saved"})
	}
}

// testJiraConnectionHandler handles POST /api/v1/settings/jira/test.
// Uses stored settings (or request body) to verify the Jira connection.
func testJiraConnectionHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// リクエストボディで上書きも可能（保存前テスト用）
		var req struct {
			JiraURL  string `json:"jira_url"`
			Email    string `json:"email"`
			APIToken string `json:"api_token"`
		}
		_ = c.ShouldBindJSON(&req)

		// リクエストが不完全な場合はDBの設定を使用
		if req.JiraURL == "" || req.Email == "" || req.APIToken == "" {
			var row jiraSettingsRow
			if err := db.QueryRowx(`SELECT jira_url, email, api_token FROM jira_settings ORDER BY id LIMIT 1`).StructScan(&row); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Jira settings not configured"})
				return
			}
			if req.JiraURL == "" {
				req.JiraURL = row.JiraURL
			}
			if req.Email == "" {
				req.Email = row.Email
			}
			if req.APIToken == "" {
				req.APIToken = row.APIToken
			}
		}

		client := jiraclient.New(jiraclient.Config{
			BaseURL:  strings.TrimRight(req.JiraURL, "/"),
			Email:    req.Email,
			APIToken: req.APIToken,
		})

		if err := client.Ping(); err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Jira connection failed: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Jira connection successful"})
	}
}

// triggerSyncHandler handles POST /api/v1/settings/jira/sync.
// Reads Jira settings from DB, constructs a Syncer, and starts a full sync asynchronously.
func triggerSyncHandler(db *sqlx.DB, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 既に RUNNING 状態のジョブがある場合はスキップ
		var running int
		if err := db.QueryRowx(`SELECT COUNT(*) FROM sync_logs WHERE status = 'RUNNING'`).Scan(&running); err == nil && running > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "sync is already running"})
			return
		}

		// Jira 設定を DB から読み込む
		var settings jiraSettingsRow
		if err := db.QueryRowx(`SELECT jira_url, email, api_token FROM jira_settings ORDER BY id LIMIT 1`).StructScan(&settings); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Jira settings not configured"})
			return
		}

		client := jiraclient.New(jiraclient.Config{
			BaseURL:  strings.TrimRight(settings.JiraURL, "/"),
			Email:    settings.Email,
			APIToken: settings.APIToken,
		})
		repo := batch.NewRepository(db)
		syncer := batch.NewSyncer(client, repo, log, 0)

		// フルシンクを非同期で実行（sync_log の管理は Syncer が担当）
		go func() {
			if err := syncer.RunFullSync(context.Background()); err != nil {
				log.Error("full sync failed", zap.Error(err))
			}
		}()

		c.JSON(http.StatusAccepted, gin.H{"message": "sync started"})
	}
}

// listSyncLogsHandler handles GET /api/v1/sync-logs.
// Returns the latest 20 sync log entries.
func listSyncLogsHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logs := make([]syncLogRow, 0)
		err := db.Select(&logs, `
			SELECT id, sync_type, executed_at, completed_at, status,
			       projects_synced, issues_synced, error_message, duration_seconds
			FROM sync_logs
			ORDER BY executed_at DESC
			LIMIT 20`,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sync logs"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": logs})
	}
}
