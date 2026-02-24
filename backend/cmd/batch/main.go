package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/m19cmjigen/sandbox-project-management/backend/internal/batch"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/config"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/jiraclient"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/logger"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/metrics"
	"github.com/m19cmjigen/sandbox-project-management/backend/pkg/secrets"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "batch error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log, err := logger.New(cfg.Log.Level, cfg.Log.Format)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer log.Sync()

	// DB 接続
	db, err := sqlx.Connect("postgres", cfg.Database.GetDSN())
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()

	// Jira 認証情報を環境変数から読み込む
	// 本番環境では ECS タスク定義の secrets フィールドで AWS Secrets Manager から注入される
	// 詳細は docs/secrets-management.md を参照
	jiraCreds, err := secrets.LoadJira()
	if err != nil {
		return fmt.Errorf("load jira credentials: %w", err)
	}

	jiraClient := jiraclient.New(jiraclient.Config{
		BaseURL:  jiraCreds.BaseURL,
		Email:    jiraCreds.Email,
		APIToken: jiraCreds.APIToken,
	})

	workerCount, _ := strconv.Atoi(getEnv("BATCH_WORKER_COUNT", "5"))
	// BATCH_SYNC_MODE: "full"（デフォルト）または "delta"
	syncMode := getEnv("BATCH_SYNC_MODE", "full")
	// METRICS_NAMESPACE: CloudWatch メトリクスのネームスペース。空の場合はメトリクス送信を無効化
	metricsNamespace := getEnv("METRICS_NAMESPACE", "")

	repo := batch.NewRepository(db)
	syncer := batch.NewSyncer(jiraClient, repo, log.Logger, workerCount)

	// METRICS_NAMESPACE が設定されている場合は CloudWatch EMF でメトリクスを送信する
	if metricsNamespace != "" {
		syncer.SetRecorder(metrics.NewEMFRecorder(metricsNamespace, os.Stdout))
		log.Info("metrics enabled", zap.String("namespace", metricsNamespace))
	}

	log.Info("starting batch",
		zap.String("jira_base_url", jiraCreds.BaseURL),
		zap.Int("worker_count", workerCount),
		zap.String("sync_mode", syncMode),
	)

	switch syncMode {
	case "delta":
		return syncer.RunDeltaSync(context.Background())
	default:
		return syncer.RunFullSync(context.Background())
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
