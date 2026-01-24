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
	projectUsecase := usecase.NewProjectUsecase(projectRepo)
	issueUsecase := usecase.NewIssueUsecase(issueRepo)
	dashboardUsecase := usecase.NewDashboardUsecase(orgRepo, projectRepo, issueRepo)

	// ハンドラーの初期化
	orgHandler := handler.NewOrganizationHandler(orgUsecase, log)
	projectHandler := handler.NewProjectHandler(projectUsecase, log)
	issueHandler := handler.NewIssueHandler(issueUsecase, log)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase, log)

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
			projects.GET("", projectHandler.ListProjects)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id/organization", projectHandler.AssignProjectToOrganization)
			projects.GET("/:id/issues", issueHandler.ListProjectIssues)
		}

		// チケット管理
		issues := v1.Group("/issues")
		{
			issues.GET("", issueHandler.ListIssues)
			issues.GET("/:id", issueHandler.GetIssue)
		}

		// ダッシュボード
		dashboard := v1.Group("/dashboard")
		{
			dashboard.GET("/summary", dashboardHandler.GetDashboardSummary)
			dashboard.GET("/organizations/:id", dashboardHandler.GetOrganizationSummary)
			dashboard.GET("/projects/:id", dashboardHandler.GetProjectSummary)
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

