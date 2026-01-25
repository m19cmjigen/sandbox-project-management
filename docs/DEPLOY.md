# デプロイガイド

## 目次

1. [本番環境の前提条件](#本番環境の前提条件)
2. [環境変数の設定](#環境変数の設定)
3. [Dockerを使用したデプロイ](#dockerを使用したデプロイ)
4. [CI/CDパイプライン](#cicdパイプライン)
5. [モニタリングとロギング](#モニタリングとロギング)
6. [バックアップとリストア](#バックアップとリストア)
7. [トラブルシューティング](#トラブルシューティング)

---

## 本番環境の前提条件

### 必須要件

- **Docker** 20.10+
- **Docker Compose** 2.0+
- **最小リソース**:
  - CPU: 2コア
  - メモリ: 4GB
  - ディスク: 20GB

### 推奨要件

- **CPU**: 4コア以上
- **メモリ**: 8GB以上
- **ディスク**: 50GB以上（SSD推奨）
- **ネットワーク**: 安定したインターネット接続

---

## 環境変数の設定

### 1. 本番用環境ファイルの作成

```bash
cp .env.example .env.production
```

### 2. 本番環境変数の設定

`.env.production`を編集：

```bash
# Database Configuration
DB_USER=prod_user
DB_PASSWORD=your-strong-password-here
DB_NAME=project_visualization

# Server Configuration
BACKEND_PORT=8080
FRONTEND_PORT=80
GIN_MODE=release
LOG_LEVEL=info

# Jira Cloud API（必須）
JIRA_BASE_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-production-api-token

# Batch Job Configuration
SYNC_INTERVAL=1h
DEFAULT_ORGANIZATION_ID=1
```

### セキュリティ上の注意

⚠️ **重要**: 本番環境では以下を必ず実施してください：

1. **強力なパスワードを使用**:
   ```bash
   # ランダムなパスワード生成
   openssl rand -base64 32
   ```

2. **.env.productionをGit管理に含めない**:
   ```bash
   echo ".env.production" >> .gitignore
   ```

3. **環境変数ファイルの権限を制限**:
   ```bash
   chmod 600 .env.production
   ```

---

## Dockerを使用したデプロイ

### 1. イメージのビルド

```bash
# すべてのサービスをビルド
docker-compose -f docker-compose.prod.yml build

# 特定のサービスのみビルド
docker-compose -f docker-compose.prod.yml build backend
docker-compose -f docker-compose.prod.yml build frontend
```

### 2. サービスの起動

```bash
# すべてのサービスを起動（デタッチモード）
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

# ログを確認
docker-compose -f docker-compose.prod.yml logs -f
```

### 3. データベースのマイグレーション

初回デプロイ時のみ実行：

```bash
# データベースコンテナに接続
docker exec -it project-viz-db-prod psql -U prod_user -d project_visualization

# または、マイグレーションツールを使用
docker run --rm \
  --network sandbox-project-management_app-network \
  -v $(pwd)/database/migrations:/migrations \
  migrate/migrate \
  -path=/migrations \
  -database postgres://prod_user:your-password@postgres:5432/project_visualization?sslmode=disable \
  up
```

### 4. 初期同期の実行

```bash
# 手動で初期同期を実行
docker-compose -f docker-compose.prod.yml exec backend ./sync -mode=once -org-id=1
```

### 5. サービスの確認

```bash
# ヘルスチェック
curl http://localhost:8080/health
curl http://localhost:8080/ready

# フロントエンドの確認
curl http://localhost/

# コンテナの状態確認
docker-compose -f docker-compose.prod.yml ps
```

---

## CI/CDパイプライン

### GitHub Actionsの設定

CI/CDワークフローは`.github/workflows/ci.yml`で定義されています。

#### 自動実行されるジョブ

1. **Backend Tests**: Goのユニットテスト
2. **Backend Lint**: golangci-lintによる静的解析
3. **Frontend Tests**: TypeScriptの型チェックとビルド
4. **Docker Build**: Dockerイメージのビルドテスト

#### トリガー条件

- `main`、`develop`、`claude/**`ブランチへのプッシュ
- プルリクエストの作成/更新

#### Secretsの設定

GitHub Actionsで本番デプロイを行う場合、以下のSecretsを設定：

1. GitHubリポジトリ → Settings → Secrets and variables → Actions
2. 以下を追加:
   - `JIRA_BASE_URL`
   - `JIRA_EMAIL`
   - `JIRA_API_TOKEN`
   - `DB_PASSWORD`

---

## モニタリングとロギング

### ログの確認

```bash
# すべてのサービスのログ
docker-compose -f docker-compose.prod.yml logs -f

# 特定のサービスのログ
docker-compose -f docker-compose.prod.yml logs -f backend
docker-compose -f docker-compose.prod.yml logs -f frontend
docker-compose -f docker-compose.prod.yml logs -f postgres

# 最新100行のみ表示
docker-compose -f docker-compose.prod.yml logs --tail=100 backend
```

### ヘルスチェック

すべてのサービスにヘルスチェックが設定されています：

```bash
# バックエンド
curl http://localhost:8080/health

# フロントエンド
curl http://localhost/health.html

# データベース（コンテナ内）
docker exec project-viz-db-prod pg_isready -U prod_user
```

### リソース使用状況の確認

```bash
# コンテナのリソース使用状況
docker stats

# ディスク使用状況
docker system df

# 詳細なディスク使用状況
docker system df -v
```

---

## バックアップとリストア

### データベースのバックアップ

#### 自動バックアップスクリプト

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/path/to/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/backup_$DATE.sql"

docker exec project-viz-db-prod pg_dump -U prod_user project_visualization > "$BACKUP_FILE"

# 圧縮
gzip "$BACKUP_FILE"

# 30日以上前のバックアップを削除
find "$BACKUP_DIR" -name "backup_*.sql.gz" -mtime +30 -delete

echo "Backup completed: $BACKUP_FILE.gz"
```

#### バックアップの実行

```bash
chmod +x backup.sh
./backup.sh
```

#### Cronで自動バックアップ

```bash
# 毎日午前2時にバックアップ
0 2 * * * /path/to/backup.sh >> /var/log/backup.log 2>&1
```

### データベースのリストア

```bash
# バックアップファイルを解凍
gunzip backup_20240101_020000.sql.gz

# リストア
docker exec -i project-viz-db-prod psql -U prod_user project_visualization < backup_20240101_020000.sql
```

---

## ローリングアップデート

### ゼロダウンタイムデプロイ

```bash
# 1. 新しいイメージをビルド
docker-compose -f docker-compose.prod.yml build

# 2. サービスを順次更新
docker-compose -f docker-compose.prod.yml up -d --no-deps --build backend

# 3. ヘルスチェック
curl http://localhost:8080/health

# 4. フロントエンドを更新
docker-compose -f docker-compose.prod.yml up -d --no-deps --build frontend
```

### ロールバック

```bash
# 前のイメージに戻す
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d

# または特定のイメージタグを指定
docker tag project-viz-backend:previous project-viz-backend:latest
docker-compose -f docker-compose.prod.yml up -d backend
```

---

## スケーリング

### 水平スケーリング

```bash
# バックエンドを3インスタンスに
docker-compose -f docker-compose.prod.yml up -d --scale backend=3

# ロードバランサー（nginx）の設定例
upstream backend {
    server backend_1:8080;
    server backend_2:8080;
    server backend_3:8080;
}
```

---

## トラブルシューティング

### コンテナが起動しない

```bash
# ログを確認
docker-compose -f docker-compose.prod.yml logs backend

# コンテナの状態を確認
docker-compose -f docker-compose.prod.yml ps

# コンテナを再作成
docker-compose -f docker-compose.prod.yml up -d --force-recreate backend
```

### データベース接続エラー

```bash
# PostgreSQLのログを確認
docker-compose -f docker-compose.prod.yml logs postgres

# データベースに直接接続
docker exec -it project-viz-db-prod psql -U prod_user -d project_visualization

# 接続テスト
docker exec project-viz-db-prod pg_isready -U prod_user
```

### ディスク容量不足

```bash
# 未使用のDockerリソースをクリーンアップ
docker system prune -a --volumes

# 古いイメージを削除
docker image prune -a

# ログファイルをローテート
docker-compose -f docker-compose.prod.yml restart
```

### パフォーマンス問題

```bash
# クエリのパフォーマンス分析
docker exec -it project-viz-db-prod psql -U prod_user -d project_visualization
# \timing on
# SELECT ... （遅いクエリ）

# PostgreSQLの設定を最適化
# shared_buffers, effective_cache_size, work_mem などを調整
```

---

## セキュリティベストプラクティス

1. **HTTPSを有効にする**: Let's EncryptやCloudflareを使用
2. **ファイアウォールを設定**: 必要なポートのみ開放
3. **定期的な更新**: Dockerイメージとパッケージの更新
4. **ログ監視**: 異常なアクセスを検知
5. **バックアップの暗号化**: バックアップファイルを暗号化

---

## パフォーマンスチューニング

### PostgreSQL

```sql
-- postgresql.conf の推奨設定
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 4MB
min_wal_size = 1GB
max_wal_size = 4GB
```

### バックエンド

```bash
# Go 環境変数
GOMAXPROCS=4
GOGC=100
```

---

## 参考リンク

- [Docker公式ドキュメント](https://docs.docker.com/)
- [PostgreSQL公式ドキュメント](https://www.postgresql.org/docs/)
- [GitHub Actions公式ドキュメント](https://docs.github.com/actions)
