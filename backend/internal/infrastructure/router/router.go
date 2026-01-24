package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/postgres"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/handler"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
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

	// リポジトリの初期化
	orgRepo := postgres.NewOrganizationRepository(db)
	projectRepo := postgres.NewProjectRepository(db)
	issueRepo := postgres.NewIssueRepository(db)

	// ユースケースの初期化
	orgUsecase := usecase.NewOrganizationUsecase(orgRepo)

	// ハンドラーの初期化
	orgHandler := handler.NewOrganizationHandler(orgUsecase, log)

	// ヘルスチェックエンドポイント
	r.GET("/health", healthCheckHandler(db))
	r.GET("/ready", readinessCheckHandler(db))

	// API v1 ルートグループ
	v1 := r.Group("/api/v1")
	{
		// 組織管理
		organizations := v1.Group("/organizations")
		{
			organizations.GET("", orgHandler.ListOrganizations)
			organizations.GET("/tree", orgHandler.GetOrganizationTree)
			organizations.GET("/:id", orgHandler.GetOrganization)
			organizations.POST("", orgHandler.CreateOrganization)
			organizations.PUT("/:id", orgHandler.UpdateOrganization)
			organizations.DELETE("/:id", orgHandler.DeleteOrganization)
			organizations.GET("/:id/children", orgHandler.GetOrganizationChildren)
		}

		// プロジェクト管理
		projects := v1.Group("/projects")
		{
			projects.GET("", listProjectsHandler(projectRepo))
			projects.GET("/:id", getProjectHandler(projectRepo))
			projects.PUT("/:id", updateProjectHandler(projectRepo))
			projects.PUT("/:id/organization", assignProjectToOrganizationHandler(projectRepo))
			projects.GET("/:id/issues", listProjectIssuesHandler(issueRepo))
		}

		// チケット管理
		issues := v1.Group("/issues")
		{
			issues.GET("", listIssuesHandler(issueRepo))
			issues.GET("/:id", getIssueHandler(issueRepo))
		}

		// ダッシュボード
		dashboard := v1.Group("/dashboard")
		{
			dashboard.GET("/summary", getDashboardSummaryHandler(projectRepo))
			dashboard.GET("/organizations/:id", getOrganizationSummaryHandler(orgRepo, projectRepo))
			dashboard.GET("/projects/:id", getProjectSummaryHandler(projectRepo))
		}
	}

	return r
}

// LoggerMiddleware はリクエストログを出力するミドルウェア
func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Info("Request processed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
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
			"status":  "ok",
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
			"status":   "ready",
			"database": "connected",
		})
	}
}

// 以下はプレースホルダーハンドラー（BACK-004以降で実装）

func listProjectsHandler(repo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list projects - not implemented yet"})
	}
}

func getProjectHandler(repo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get project - not implemented yet"})
	}
}

func updateProjectHandler(repo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "update project - not implemented yet"})
	}
}

func assignProjectToOrganizationHandler(repo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "assign project to organization - not implemented yet"})
	}
}

func listProjectIssuesHandler(repo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list project issues - not implemented yet"})
	}
}

func listIssuesHandler(repo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "list issues - not implemented yet"})
	}
}

func getIssueHandler(repo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get issue - not implemented yet"})
	}
}

func getDashboardSummaryHandler(repo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get dashboard summary - not implemented yet"})
	}
}

func getOrganizationSummaryHandler(orgRepo, projRepo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get organization summary - not implemented yet"})
	}
}

func getProjectSummaryHandler(repo interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "get project summary - not implemented yet"})
	}
}
