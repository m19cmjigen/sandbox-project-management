.PHONY: help db-up db-down db-migrate db-rollback db-version db-force db-migration db-drop db-create db-reset db-seed \
        backend-run backend-build backend-test backend-lint backend-fmt backend-deps backend-clean \
        sync-once sync-scheduler sync-build \
        perf-smoke perf-load perf-stress perf-large-data

# デフォルトのデータベースURL（環境変数で上書き可能）
DATABASE_URL ?= postgres://admin:admin123@localhost:5432/project_visualization?sslmode=disable
MIGRATION_PATH = database/migrations
BACKEND_DIR = backend

# ヘルプ
help:
	@echo "利用可能なコマンド:"
	@echo ""
	@echo "データベース:"
	@echo "  make db-up          - PostgreSQLコンテナを起動"
	@echo "  make db-down        - PostgreSQLコンテナを停止・削除"
	@echo "  make db-migrate     - マイグレーションを適用"
	@echo "  make db-rollback    - マイグレーションを1ステップロールバック"
	@echo "  make db-version     - 現在のマイグレーションバージョンを確認"
	@echo "  make db-force V=1   - マイグレーションバージョンを強制設定"
	@echo "  make db-migration name=feature_name - 新しいマイグレーションファイルを作成"
	@echo "  make db-reset       - データベースをリセット（全削除後に再作成）"
	@echo "  make db-seed        - シードデータを投入（開発・テスト用）"
	@echo ""
	@echo "バックエンド:"
	@echo "  make backend-run    - バックエンドAPIを起動"
	@echo "  make backend-build  - バックエンドをビルド"
	@echo "  make backend-test   - バックエンドのテストを実行"
	@echo "  make backend-lint   - バックエンドのリントを実行"
	@echo "  make backend-fmt    - バックエンドのフォーマットを実行"
	@echo "  make backend-deps   - バックエンドの依存関係を更新"
	@echo "  make backend-clean  - バックエンドのビルド成果物を削除"
	@echo ""
	@echo "Jira同期:"
	@echo "  make sync-once      - Jira同期を1回実行"
	@echo "  make sync-scheduler - Jira同期をスケジューラーモードで実行"
	@echo "  make sync-build     - 同期CLIツールをビルド"
	@echo ""
	@echo "パフォーマンステスト:"
	@echo "  make perf-smoke     - スモークテスト（疎通確認）"
	@echo "  make perf-load      - 負荷テスト（50-100 VU）"
	@echo "  make perf-stress    - ストレステスト（最大200 VU）"
	@echo "  make perf-large-data - 大量データ生成（10,000+ issues）"
	@echo ""
	@echo "Docker Compose:"
	@echo "  make up             - すべてのサービスを起動"
	@echo "  make down           - すべてのサービスを停止"
	@echo "  make logs           - ログを表示"

# PostgreSQLコンテナの起動
db-up:
	@echo "PostgreSQLコンテナを起動中..."
	docker run --name project-viz-db \
		-e POSTGRES_USER=admin \
		-e POSTGRES_PASSWORD=admin123 \
		-e POSTGRES_DB=project_visualization \
		-p 5432:5432 \
		-d postgres:15-alpine
	@echo "PostgreSQLが起動しました: localhost:5432"

# PostgreSQLコンテナの停止・削除
db-down:
	@echo "PostgreSQLコンテナを停止中..."
	docker stop project-viz-db || true
	docker rm project-viz-db || true
	@echo "PostgreSQLコンテナを削除しました"

# マイグレーションの適用
db-migrate:
	@echo "マイグレーションを適用中..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" up
	@echo "マイグレーションが完了しました"

# マイグレーションのロールバック
db-rollback:
	@echo "マイグレーションをロールバック中..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down 1
	@echo "ロールバックが完了しました"

# マイグレーション全体のロールバック
db-rollback-all:
	@echo "すべてのマイグレーションをロールバック中..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" down
	@echo "すべてのロールバックが完了しました"

# マイグレーションバージョンの確認
db-version:
	@migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" version

# マイグレーションバージョンの強制設定（エラーリカバリ用）
db-force:
	@if [ -z "$(V)" ]; then \
		echo "エラー: バージョンを指定してください。例: make db-force V=1"; \
		exit 1; \
	fi
	@echo "マイグレーションバージョンを $(V) に強制設定中..."
	migrate -path $(MIGRATION_PATH) -database "$(DATABASE_URL)" force $(V)
	@echo "バージョンが設定されました"

# 新しいマイグレーションファイルの作成
db-migration:
	@if [ -z "$(name)" ]; then \
		echo "エラー: マイグレーション名を指定してください。例: make db-migration name=add_users"; \
		exit 1; \
	fi
	@echo "マイグレーションファイルを作成中: $(name)"
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name)
	@echo "マイグレーションファイルが作成されました"

# データベースのリセット（開発環境用）
db-reset: db-rollback-all db-migrate
	@echo "データベースがリセットされました"

# データベースへの接続（psql）
db-connect:
	@echo "PostgreSQLに接続中..."
	psql "$(DATABASE_URL)"

# シードデータの投入（開発・テスト用）
db-seed:
	@echo "シードデータを投入中..."
	cd $(BACKEND_DIR) && DATABASE_URL="$(DATABASE_URL)" go run ./scripts/seed_data.go
	@echo "シードデータの投入が完了しました"

# =====================================
# バックエンドコマンド
# =====================================

# バックエンドAPIの起動
backend-run:
	@echo "バックエンドAPIを起動中..."
	cd $(BACKEND_DIR) && go run cmd/api/main.go

# バックエンドのビルド
backend-build:
	@echo "バックエンドをビルド中..."
	cd $(BACKEND_DIR) && go build -o bin/api cmd/api/main.go
	@echo "ビルド完了: backend/bin/api"

# バックエンドのテスト実行
backend-test:
	@echo "バックエンドのテストを実行中..."
	cd $(BACKEND_DIR) && go test -v ./...

# バックエンドのテストカバレッジ
backend-coverage:
	@echo "テストカバレッジを生成中..."
	cd $(BACKEND_DIR) && go test -coverprofile=coverage.out ./...
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html
	@echo "カバレッジレポート: backend/coverage.html"

# バックエンドのリント
backend-lint:
	@echo "リントを実行中..."
	cd $(BACKEND_DIR) && golangci-lint run ./...

# バックエンドのフォーマット
backend-fmt:
	@echo "コードフォーマット中..."
	cd $(BACKEND_DIR) && go fmt ./...
	cd $(BACKEND_DIR) && goimports -w .

# バックエンドの依存関係更新
backend-deps:
	@echo "依存関係を更新中..."
	cd $(BACKEND_DIR) && go mod download
	cd $(BACKEND_DIR) && go mod tidy

# バックエンドのビルド成果物削除
backend-clean:
	@echo "ビルド成果物を削除中..."
	cd $(BACKEND_DIR) && rm -rf bin/
	cd $(BACKEND_DIR) && rm -f coverage.out coverage.html

# =====================================
# Docker Composeコマンド
# =====================================

# すべてのサービスを起動
up:
	@echo "Docker Composeでサービスを起動中..."
	docker-compose up -d
	@echo "サービスが起動しました"

# すべてのサービスを停止
down:
	@echo "Docker Composeでサービスを停止中..."
	docker-compose down
	@echo "サービスが停止しました"

# ログを表示
logs:
	docker-compose logs -f

# バックエンドのログを表示
logs-backend:
	docker-compose logs -f backend

# PostgreSQLのログを表示
logs-db:
	docker-compose logs -f postgres

# =====================================
# Jira同期コマンド
# =====================================

# 同期CLIツールのビルド
sync-build:
	@echo "同期CLIツールをビルド中..."
	cd $(BACKEND_DIR) && go build -o bin/sync cmd/sync/main.go
	@echo "ビルド完了: backend/bin/sync"

# Jira同期を1回実行
sync-once:
	@echo "Jira同期を実行中（1回のみ）..."
	cd $(BACKEND_DIR) && go run cmd/sync/main.go -mode=once -org-id=1

# Jira同期をスケジューラーモードで実行
sync-scheduler:
	@echo "Jira同期をスケジューラーモードで実行中..."
	@echo "同期間隔: 1時間"
	cd $(BACKEND_DIR) && go run cmd/sync/main.go -mode=scheduler -org-id=1 -interval=1h

# =====================================
# パフォーマンステストコマンド
# =====================================

# スモークテスト（疎通確認）
perf-smoke:
	@echo "スモークテストを実行中..."
	@command -v k6 >/dev/null 2>&1 || { echo "エラー: k6がインストールされていません。README.mdを参照してください。"; exit 1; }
	k6 run performance/smoke-test.js

# 負荷テスト（通常運用を想定）
perf-load:
	@echo "負荷テストを実行中..."
	@command -v k6 >/dev/null 2>&1 || { echo "エラー: k6がインストールされていません。"; exit 1; }
	k6 run performance/api-load-test.js

# ストレステスト（高負荷・限界確認）
perf-stress:
	@echo "ストレステストを実行中..."
	@command -v k6 >/dev/null 2>&1 || { echo "エラー: k6がインストールされていません。"; exit 1; }
	k6 run performance/stress-test.js

# 大量データ生成（パフォーマンステスト用）
perf-large-data:
	@echo "大量テストデータを生成中（10,000+ issues）..."
	@echo "警告: この処理には数分かかる場合があります"
	cd $(BACKEND_DIR) && DATABASE_URL="$(DATABASE_URL)" go run ./scripts/generate_large_dataset.go
	@echo "大量データ生成が完了しました"
