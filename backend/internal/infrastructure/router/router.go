package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/auth"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/jira"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/postgres"
	httpHandler "github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/http"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/interface/http/middleware"
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
	syncLogRepo := postgres.NewSyncLogRepository(db)
	userRepo := postgres.NewUserRepository(db)
	auditLogRepo := postgres.NewAuditLogRepository(db)

	// 認証サービスの初期化
	jwtService := auth.NewJWTService(auth.JWTConfig{
		SecretKey:       cfg.JWT.SecretKey,
		ExpirationHours: cfg.JWT.ExpirationHours,
	})
	passwordService := auth.NewPasswordService()

	// ユースケースの初期化
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService, passwordService)
	auditUsecase := usecase.NewAuditUsecase(auditLogRepo)
	orgUsecase := usecase.NewOrganizationUsecase(orgRepo)
	projectUsecase := usecase.NewProjectUsecase(projectRepo)
	issueUsecase := usecase.NewIssueUsecase(issueRepo)
	dashboardUsecase := usecase.NewDashboardUsecase(orgRepo, projectRepo, issueRepo)

	// Jira統合（環境変数から設定を読み込み）
	var syncUsecase usecase.SyncUsecase
	jiraConfig, err := jira.LoadConfig()
	if err == nil {
		jiraClient := jira.NewClient(jiraConfig)
		syncUsecase = usecase.NewSyncUsecaseWithRetry(jiraClient, projectRepo, issueRepo, syncLogRepo, nil)
		log.Info("Jira integration enabled")
	} else {
		log.Warn("Jira integration disabled", zap.Error(err))
	}

	// ハンドラーの初期化
	authHandler := httpHandler.NewAuthHandler(authUsecase)
	auditHandler := httpHandler.NewAuditHandler(auditUsecase)
	orgHandler := handler.NewOrganizationHandler(orgUsecase, log)
	projectHandler := handler.NewProjectHandler(projectUsecase, log)
	issueHandler := handler.NewIssueHandler(issueUsecase, log)
	dashboardHandler := handler.NewDashboardHandler(dashboardUsecase, log)

	var syncHandler *handler.SyncHandler
	if syncUsecase != nil {
		syncHandler = handler.NewSyncHandler(syncUsecase, syncLogRepo)
	}

	// ヘルスチェックエンドポイント
	r.GET("/health", healthCheckHandler(db))
	r.GET("/ready", readinessCheckHandler(db))

	// API v1 ルートグループ
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuditMiddleware(auditUsecase)) // 監査ログミドルウェア
	{
		// 認証エンドポイント（認証不要）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.GET("/me", middleware.AuthMiddleware(authUsecase), authHandler.GetCurrentUser)
		}

		// ユーザー管理（管理者のみ）
		users := v1.Group("/users", middleware.AuthMiddleware(authUsecase))
		{
			users.GET("", middleware.RequireAdmin(), authHandler.ListUsers)
			users.POST("", middleware.RequireAdmin(), authHandler.CreateUser)
			users.GET("/:id", authHandler.GetUser)
			users.PUT("/:id", middleware.RequireAdmin(), authHandler.UpdateUser)
			users.DELETE("/:id", middleware.RequireAdmin(), authHandler.DeleteUser)
			users.POST("/:id/password", authHandler.ChangePassword)
		}

		// 認証が必要な保護されたルート
		protected := v1.Group("", middleware.AuthMiddleware(authUsecase))
		{
			// 組織管理（読み取りは全員、変更はマネージャー以上）
			organizations := protected.Group("/organizations")
			{
				organizations.GET("", orgHandler.ListOrganizations)
				organizations.GET("/tree", orgHandler.GetOrganizationTree)
				organizations.GET("/:id", orgHandler.GetOrganization)
				organizations.GET("/:id/children", orgHandler.GetOrganizationChildren)

				// 変更操作はマネージャー以上
				organizations.POST("", middleware.RequireManagerOrAdmin(), orgHandler.CreateOrganization)
				organizations.PUT("/:id", middleware.RequireManagerOrAdmin(), orgHandler.UpdateOrganization)
				organizations.DELETE("/:id", middleware.RequireManagerOrAdmin(), orgHandler.DeleteOrganization)
			}

			// プロジェクト管理（読み取りは全員、変更はマネージャー以上）
			projects := protected.Group("/projects")
			{
				projects.GET("", projectHandler.ListProjects)
				projects.GET("/:id", projectHandler.GetProject)
				projects.GET("/:id/issues", issueHandler.ListProjectIssues)

				// 変更操作はマネージャー以上
				projects.PUT("/:id/organization", middleware.RequireManagerOrAdmin(), projectHandler.AssignProjectToOrganization)
			}

			// チケット管理（読み取り専用）
			issues := protected.Group("/issues")
			{
				issues.GET("", issueHandler.ListIssues)
				issues.GET("/:id", issueHandler.GetIssue)
			}

			// ダッシュボード（読み取り専用）
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/summary", dashboardHandler.GetDashboardSummary)
				dashboard.GET("/organizations/:id", dashboardHandler.GetOrganizationSummary)
				dashboard.GET("/projects/:id", dashboardHandler.GetProjectSummary)
			}

			// Jira同期（マネージャー以上、Jira統合が有効な場合のみ）
			if syncHandler != nil {
				sync := protected.Group("/sync", middleware.RequireManagerOrAdmin())
				{
					sync.POST("/trigger", syncHandler.TriggerSync)
					sync.POST("/projects/:id", syncHandler.SyncProject)
					sync.GET("/logs", syncHandler.GetSyncLogs)
					sync.GET("/logs/latest", syncHandler.GetLatestSyncLog)
					sync.GET("/logs/:id", syncHandler.GetSyncLog)
				}
			}

			// 監査ログ（管理者のみ）
			audit := protected.Group("/audit", middleware.RequireAdmin())
			{
				audit.GET("/logs", auditHandler.ListAuditLogs)
				audit.GET("/logs/:id", auditHandler.GetAuditLog)
				audit.DELETE("/logs/cleanup", auditHandler.CleanupOldLogs)
			}
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

