package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/batch"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/jira"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/infrastructure/postgres"
	"github.com/m19cmjigen/sandbox-project-management/backend/internal/usecase"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/retry"
)

func main() {
	// コマンドラインフラグ
	mode := flag.String("mode", "once", "Sync mode: 'once' for one-time sync, 'scheduler' for continuous scheduling")
	orgID := flag.Int64("org-id", 1, "Organization ID to sync projects to")
	interval := flag.Duration("interval", 1*time.Hour, "Sync interval for scheduler mode (e.g., 1h, 30m)")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Jira Sync Tool")

	// データベース接続
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	// Jira設定の読み込み
	jiraConfig, err := jira.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load Jira configuration: %v", err)
	}

	log.Printf("Jira configuration loaded: %s", jiraConfig.BaseURL)

	// リポジトリの初期化
	projectRepo := postgres.NewProjectRepository(db)
	issueRepo := postgres.NewIssueRepository(db)
	syncLogRepo := postgres.NewSyncLogRepository(db)

	// Jira Clientの初期化
	jiraClient := jira.NewClient(jiraConfig)

	// Sync Usecaseの初期化（リトライ機能付き）
	retryConfig := retry.DefaultRetryConfig()
	syncUsecase := usecase.NewSyncUsecaseWithRetry(
		jiraClient,
		projectRepo,
		issueRepo,
		syncLogRepo,
		retryConfig,
	)

	ctx := context.Background()

	switch *mode {
	case "once":
		// 一度だけ同期を実行
		log.Printf("Running one-time sync for organization ID: %d", *orgID)
		runOnce(ctx, syncUsecase, *orgID)

	case "scheduler":
		// スケジューラーモードで継続的に実行
		log.Printf("Starting scheduler mode with interval: %v", *interval)
		runScheduler(ctx, syncUsecase, *orgID, *interval)

	default:
		log.Fatalf("Invalid mode: %s. Use 'once' or 'scheduler'", *mode)
	}
}

// runOnce executes a one-time sync
func runOnce(ctx context.Context, syncUsecase usecase.SyncUsecase, orgID int64) {
	startTime := time.Now()
	log.Println("===== Starting Jira Sync =====")

	syncLog, err := syncUsecase.SyncAllProjects(ctx, orgID)
	if err != nil {
		log.Printf("ERROR: Sync failed: %v", err)
		os.Exit(1)
	}

	duration := time.Since(startTime)

	log.Println("===== Sync Completed =====")
	log.Printf("Duration: %v", duration)
	log.Printf("Status: %s", syncLog.Status)
	log.Printf("Projects Synced: %d", syncLog.ProjectsSynced)
	log.Printf("Issues Synced: %d", syncLog.IssuesSynced)
	log.Printf("Errors: %d", syncLog.ErrorCount)

	if syncLog.ErrorMessage != nil {
		log.Printf("Error Message: %s", *syncLog.ErrorMessage)
	}

	if syncLog.Status == "FAILED" || syncLog.ErrorCount > 0 {
		os.Exit(1)
	}

	log.Println("Sync completed successfully")
}

// runScheduler runs sync jobs on a schedule
func runScheduler(ctx context.Context, syncUsecase usecase.SyncUsecase, orgID int64, interval time.Duration) {
	scheduler := batch.NewScheduler(syncUsecase, orgID, interval)

	// シグナルハンドリングのためのチャネル
	sigChan := make(chan os.Signal, 1)

	// Graceful shutdown
	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping scheduler...")
		scheduler.Stop()
	}()

	log.Printf("Scheduler started with interval: %v", interval)
	log.Println("Press Ctrl+C to stop")

	// スケジューラーを開始（ブロッキング）
	scheduler.Start(ctx)

	log.Println("Scheduler stopped gracefully")
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// formatDuration formats a duration for human-readable output
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}
