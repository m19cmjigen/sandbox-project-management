# 運用手順書

## 日常運用タスク

### サービス起動・停止

```bash
# 起動
make up && make db-migrate

# 停止
make down

# 再起動 (バックエンドのみ)
docker compose restart backend

# 状態確認
docker compose ps
```

### ログ確認

```bash
# リアルタイムログ (全サービス)
make logs

# バックエンドのみ
make logs-backend

# PostgreSQLのみ
make logs-db

# 直近100行を表示して終了
docker compose logs --tail=100 backend
```

バックエンドはJSON形式でログを出力します（Uber Zap使用）。

```json
{"level":"info","ts":"2026-02-23T12:00:00.000Z","msg":"Request processed","method":"GET","path":"/api/v1/projects","status":200,"ip":"172.18.0.1","user_agent":"Mozilla/5.0..."}
```

### ヘルスチェック

```bash
# 稼働確認
curl http://localhost:8080/health
# {"status":"ok","service":"project-visualization-api"}

# DB接続確認
curl http://localhost:8080/ready
# {"status":"ready","database":"connected"}
```

## データベース運用

### バックアップ

```bash
# pg_dumpでバックアップ
docker exec project-viz-db \
  pg_dump -U admin project_visualization > backup_$(date +%Y%m%d).sql

# 圧縮バックアップ
docker exec project-viz-db \
  pg_dump -U admin -Fc project_visualization > backup_$(date +%Y%m%d).dump
```

### リストア

```bash
# SQLファイルからリストア
docker exec -i project-viz-db \
  psql -U admin project_visualization < backup_20260223.sql

# カスタム形式からリストア
docker exec -i project-viz-db \
  pg_restore -U admin -d project_visualization backup_20260223.dump
```

### データベース接続

```bash
# psqlで直接接続
make db-connect

# または
docker exec -it project-viz-db psql -U admin project_visualization
```

### マイグレーション状態確認

```bash
make db-version
```

## パフォーマンス監視

### k6パフォーマンステストの定期実行

本番デプロイ後やリリース前に実施します。

```bash
# バックエンド起動確認後にテスト実行
BASE_URL=http://localhost:8080 ./performance/run.sh all
```

各テストの目標値:

| テスト | VU数 | p(95)目標 | エラー率目標 |
|---|---|---|---|
| スモーク | 1 | < 500ms | < 1% |
| 負荷 | 最大10 | < 500ms | < 1% |
| ストレス | 最大50 | < 2000ms | < 10% |

実績値 (2026-02-23):
- スモーク: p(95)=16ms, エラー率=0%
- 負荷: p(95)=13ms, エラー率=0%
- ストレス: p(95)=18ms, エラー率=0%

### データベースのスロークエリ確認

```sql
-- 実行時間の長いクエリを確認
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### テーブルサイズ確認

```sql
SELECT
  relname AS table_name,
  pg_size_pretty(pg_total_relation_size(relid)) AS total_size,
  pg_size_pretty(pg_relation_size(relid)) AS data_size
FROM pg_catalog.pg_statio_user_tables
ORDER BY pg_total_relation_size(relid) DESC;
```

## Jiraデータ同期

バッチ処理でJiraのプロジェクト・チケットをDBに同期します。
初回セットアップ手順は [docs/jira-setup.md](./jira-setup.md) を参照してください。

### 手動でバッチを実行する

```bash
cd backend

# フル同期（全件）
go build -o bin/batch ./cmd/batch/ && ./bin/batch

# 差分同期（前回以降の変更のみ）
BATCH_SYNC_MODE=delta ./bin/batch
```

### 同期ステータス確認

```sql
-- 最新の同期ログを確認
SELECT sync_type, status, projects_synced, issues_synced,
       duration_seconds, executed_at, error_message
FROM sync_logs
ORDER BY executed_at DESC
LIMIT 10;

-- 失敗した同期を確認
SELECT *
FROM sync_logs
WHERE status = 'FAILURE'
ORDER BY executed_at DESC;
```

### バッチスケジュール（本番）

本番環境では CloudWatch EventBridge で自動実行します。
詳細は [docs/batch-schedule.md](./batch-schedule.md) を参照してください。

## セキュリティ運用

### 認証情報のローテーション

DBパスワードのローテーション手順:

1. 新しいパスワードを生成
2. DBユーザーのパスワードを変更

```sql
ALTER USER admin WITH PASSWORD 'new_password';
```

3. `.env`または環境変数を更新
4. バックエンドを再起動

```bash
docker compose restart backend
```

### アクセスログの確認

```bash
# エラーが多いパスの確認
docker compose logs backend | grep '"status":5' | head -20

# 直近のアクセス統計
docker compose logs backend | grep '"msg":"Request processed"' | \
  jq -r '"\(.status) \(.path)"' | sort | uniq -c | sort -rn | head -20
```

## コンテナリソース監視

```bash
# コンテナのリソース使用量をリアルタイムで確認
docker stats project-viz-backend project-viz-db

# 一回だけ確認
docker stats --no-stream
```
