package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/config"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
)

// NewRouter は新しいGinルーターを作成する
func NewRouter(cfg *config.Config, db *sqlx.DB, log *logger.Logger) *gin.Engine {
	// Ginモードの設定
	gin.SetMode(cfg.Server.GinMode)

	r := gin.New()

	// ミドルウェア設定
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware(log))
	r.Use(CORSMiddleware())

	// ヘルスチェックエンドポイント
	r.GET("/health", healthCheckHandler(db))
	r.GET("/ready", readinessCheckHandler(db))

	// API v1 ルートグループ
	v1 := r.Group("/api/v1")
	{
		// 組織管理
		organizations := v1.Group("/organizations")
		{
			organizations.GET("", listOrganizationsHandlerWithDB(db))
			organizations.GET("/:id", getOrganizationHandlerWithDB(db))
			organizations.POST("", createOrganizationHandlerWithDB(db))
			organizations.PUT("/:id", updateOrganizationHandlerWithDB(db))
			organizations.DELETE("/:id", deleteOrganizationHandlerWithDB(db))
			organizations.GET("/:id/children", getChildOrganizationsHandlerWithDB(db))
		}

		// プロジェクト管理
		projects := v1.Group("/projects")
		{
			projects.GET("", listProjectsHandlerWithDB(db))
			projects.GET("/:id", getProjectHandlerWithDB(db))
			projects.PUT("/:id", updateProjectHandler)
			projects.PUT("/:id/organization", assignProjectToOrganizationHandlerWithDB(db))
			projects.GET("/:id/issues", listProjectIssuesHandler)
		}

		// チケット管理
		issues := v1.Group("/issues")
		{
			issues.GET("", listIssuesHandlerWithDB(db))
			issues.GET("/:id", getIssueHandlerWithDB(db))
		}

		// ダッシュボード
		dashboard := v1.Group("/dashboard")
		{
			dashboard.GET("/summary", getDashboardSummaryHandlerWithDB(db))
			dashboard.GET("/organizations/:id", getOrganizationSummaryHandlerWithDB(db))
			dashboard.GET("/projects/:id", getProjectSummaryHandler)
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

// CORSMiddleware はCORS設定を行うミドルウェア
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
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
