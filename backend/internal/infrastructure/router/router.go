package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/auth"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/config"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
)

// NewRouter は新しいGinルーターを作成する
func NewRouter(cfg *config.Config, db *sqlx.DB, log *logger.Logger) *gin.Engine {
	// Ginモードの設定
	gin.SetMode(cfg.Server.GinMode)

	r := gin.New()

	tm := auth.NewTokenManager(cfg.Auth.JWTSecret)

	// 共通ミドルウェア
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware(log))
	r.Use(SecurityHeadersMiddleware())
	r.Use(CORSMiddleware(cfg.Auth.AllowedOrigins))

	// 公開エンドポイント（認証不要）
	r.GET("/health", healthCheckHandler(db))
	r.GET("/ready", readinessCheckHandler(db))

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 認証エンドポイント（JWT不要）
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", loginHandler(db, tm))
		}

		// 認証が必要なエンドポイント
		protected := v1.Group("")
		protected.Use(auth.Middleware(tm))
		{
			// 認証ユーザー情報
			protected.GET("/auth/me", meHandler())

			// 組織管理
			organizations := protected.Group("/organizations")
			{
				organizations.GET("", listOrganizationsHandlerWithDB(db))
				organizations.GET("/:id", getOrganizationHandlerWithDB(db))
				organizations.GET("/:id/children", getChildOrganizationsHandlerWithDB(db))
				// admin のみ書き込み可
				organizations.POST("", auth.RequireRole("admin"), createOrganizationHandlerWithDB(db))
				organizations.PUT("/:id", auth.RequireRole("admin"), updateOrganizationHandlerWithDB(db))
				organizations.DELETE("/:id", auth.RequireRole("admin"), deleteOrganizationHandlerWithDB(db))
			}

			// プロジェクト管理
			projects := protected.Group("/projects")
			{
				projects.GET("", listProjectsHandlerWithDB(db))
				projects.GET("/:id", getProjectHandlerWithDB(db))
				projects.GET("/:id/issues", listProjectIssuesHandlerWithDB(db))
				// admin のみ書き込み可
				projects.PUT("/:id", auth.RequireRole("admin"), updateProjectHandlerWithDB(db))
				// admin + project_manager が組織割り当て可能
				projects.PUT("/:id/organization", auth.RequireRole("admin", "project_manager"), assignProjectToOrganizationHandlerWithDB(db))
			}

			// ユーザー管理 (admin のみ)
			users := protected.Group("/users")
			users.Use(auth.RequireRole("admin"))
			{
				users.GET("", listUsersHandlerWithDB(db))
				users.POST("", createUserHandlerWithDB(db))
				users.PUT("/:id", updateUserHandlerWithDB(db))
				users.DELETE("/:id", deleteUserHandlerWithDB(db))
			}

			// 設定管理 (admin のみ)
			settings := protected.Group("/settings")
			settings.Use(auth.RequireRole("admin"))
			{
				settings.GET("/jira", getJiraSettingsHandler(db))
				settings.PUT("/jira", updateJiraSettingsHandler(db))
				settings.POST("/jira/test", testJiraConnectionHandler(db))
				settings.POST("/jira/sync", triggerSyncHandler(db))
			}

			// 同期ログ (admin のみ)
			protected.GET("/sync-logs", auth.RequireRole("admin"), listSyncLogsHandler(db))

			// チケット管理（読み取り専用）
			issues := protected.Group("/issues")
			{
				issues.GET("", listIssuesHandlerWithDB(db))
				issues.GET("/:id", getIssueHandlerWithDB(db))
			}

			// ダッシュボード（読み取り専用）
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/summary", getDashboardSummaryHandlerWithDB(db))
				dashboard.GET("/organizations/:id", getOrganizationSummaryHandlerWithDB(db))
				dashboard.GET("/projects/:id", getProjectSummaryHandlerWithDB(db))
			}
		}
	}

	return r
}

// LoggerMiddleware はリクエストログを出力するミドルウェア
func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Request.Context().Value("start")
		if start == nil {
			start = c.GetTime("start")
		}

		c.Next()

		log.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// SecurityHeadersMiddleware はセキュリティ関連HTTPヘッダーを設定するミドルウェア
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'none'")
		c.Next()
	}
}

// CORSMiddleware はCORS設定を行うミドルウェア。
// allowedOrigins にカンマ区切りでオリジンを指定する。空の場合はワイルドカード（開発用）。
func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
	// 許可オリジンセットを事前構築
	originSet := make(map[string]struct{})
	if allowedOrigins != "" {
		for _, o := range strings.Split(allowedOrigins, ",") {
			if o = strings.TrimSpace(o); o != "" {
				originSet[o] = struct{}{}
			}
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowOrigin := "*"
		if len(originSet) > 0 {
			if _, ok := originSet[origin]; ok {
				allowOrigin = origin
			} else {
				// 許可されていないオリジンには CORS ヘッダーを付与しない
				c.Next()
				return
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, Cache-Control")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// healthCheckHandler はヘルスチェックハンドラー
func healthCheckHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "project-visualization-api",
		})
	}
}

// readinessCheckHandler はReadinessチェックハンドラー
func readinessCheckHandler(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// データベース接続確認
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"error":  "database connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"database": "connected",
		})
	}
}

// 以下はプレースホルダーハンドラー（実装は後続チケットで行う）

func listOrganizationsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list organizations - not implemented yet"})
}

func getOrganizationHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get organization - not implemented yet"})
}

func createOrganizationHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "create organization - not implemented yet"})
}

func updateOrganizationHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update organization - not implemented yet"})
}

func deleteOrganizationHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "delete organization - not implemented yet"})
}

func getChildOrganizationsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get child organizations - not implemented yet"})
}

func listProjectsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list projects - not implemented yet"})
}

func getProjectHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get project - not implemented yet"})
}

func updateProjectHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update project - not implemented yet"})
}

func assignProjectToOrganizationHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "assign project to organization - not implemented yet"})
}

func listProjectIssuesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list project issues - not implemented yet"})
}

func listIssuesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list issues - not implemented yet"})
}

func getIssueHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get issue - not implemented yet"})
}

func getDashboardSummaryHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get dashboard summary - not implemented yet"})
}

func getOrganizationSummaryHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get organization summary - not implemented yet"})
}

func getProjectSummaryHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get project summary - not implemented yet"})
}
