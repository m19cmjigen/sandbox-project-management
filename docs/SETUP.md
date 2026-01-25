# セットアップガイド

## 目次

1. [前提条件](#前提条件)
2. [環境変数の設定](#環境変数の設定)
3. [データベースのセットアップ](#データベースのセットアップ)
4. [バックエンドのセットアップ](#バックエンドのセットアップ)
5. [フロントエンドのセットアップ](#フロントエンドのセットアップ)
6. [Jira統合の設定](#jira統合の設定)
7. [トラブルシューティング](#トラブルシューティング)

---

## 前提条件

以下のソフトウェアがインストールされている必要があります：

### 必須
- **Docker** 20.10+
- **Docker Compose** 2.0+
- **Go** 1.21+
- **Node.js** 18+
- **npm** 9+
- **golang-migrate** CLI

### 推奨
- **Make** （Makefileの実行用）
- **Git**
- **PostgreSQL Client** （psql、デバッグ用）

### インストール方法

#### golang-migrate

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Windows
scoop install migrate
```

---

## 環境変数の設定

### 1. 環境変数ファイルの作成

```bash
cp .env.example .env
```

### 2. 必須の環境変数を設定

`.env`ファイルを開いて以下を設定：

```bash
# Database Configuration
DATABASE_URL=postgres://admin:admin123@localhost:5432/project_visualization?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=project_visualization

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
GIN_MODE=release

# Jira Cloud API Configuration（必須）
JIRA_BASE_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-jira-api-token

# Batch Job Configuration
SYNC_INTERVAL=1h
DEFAULT_ORGANIZATION_ID=1

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### 3. Jira APIトークンの取得

1. [Atlassian Account Settings](https://id.atlassian.com/manage-profile/security/api-tokens) にアクセス
2. 「Create API token」をクリック
3. トークン名を入力（例: "Project Visualization Platform"）
4. 生成されたトークンをコピーして`JIRA_API_TOKEN`に設定

---

## データベースのセットアップ

### 1. PostgreSQLコンテナの起動

```bash
make db-up
```

または手動で：

```bash
docker run --name project-viz-db \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=admin123 \
  -e POSTGRES_DB=project_visualization \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### 2. データベース接続の確認

```bash
# psqlで接続
make db-connect

# または
psql "postgres://admin:admin123@localhost:5432/project_visualization?sslmode=disable"
```

### 3. マイグレーションの適用

```bash
make db-migrate
```

マイグレーションが成功すると、以下のテーブルが作成されます：
- `organizations`
- `projects`
- `issues`
- `sync_logs`
- `project_statistics`（ビュー）

### 4. マイグレーションの確認

```bash
make db-version
```

現在のバージョン番号が表示されれば成功です。

---

## バックエンドのセットアップ

### 1. 依存関係のインストール

```bash
cd backend
go mod download
go mod tidy
```

### 2. ビルド

```bash
make backend-build
```

バイナリが `backend/bin/api` に生成されます。

### 3. APIサーバーの起動

```bash
make backend-run
```

または直接実行：

```bash
cd backend
go run cmd/api/main.go
```

### 4. 動作確認

別のターミナルで：

```bash
# ヘルスチェック
curl http://localhost:8080/health

# Readinessチェック
curl http://localhost:8080/ready

# 組織一覧（空の配列が返る）
curl http://localhost:8080/api/v1/organizations
```

---

## フロントエンドのセットアップ

### 1. 依存関係のインストール

```bash
cd frontend
npm install
```

### 2. 開発サーバーの起動

```bash
npm run dev
```

ブラウザで `http://localhost:5173` にアクセスします。

### 3. プロダクションビルド

```bash
npm run build
```

ビルド成果物は `frontend/dist/` に生成されます。

---

## Jira統合の設定

### 1. 初期組織の作成

まず、プロジェクトを割り当てるための組織を作成します：

```bash
curl -X POST http://localhost:8080/api/v1/organizations \
  -H "Content-Type: application/json" \
  -d '{
    "name": "全社",
    "parent_id": null
  }'
```

レスポンスから組織IDを確認します（例: `"id": 1`）。

### 2. 同期の実行

#### 方法1: CLIツールを使用

```bash
# 一度だけ同期
make sync-once

# スケジューラーモード（1時間ごと）
make sync-scheduler

# カスタム間隔
cd backend
go run cmd/sync/main.go -mode=scheduler -org-id=1 -interval=30m
```

#### 方法2: APIを使用

```bash
curl -X POST http://localhost:8080/api/v1/sync/trigger \
  -H "Content-Type: application/json" \
  -d '{
    "organization_id": 1
  }'
```

#### 方法3: UIから実行

1. ブラウザで `http://localhost:5173/admin` にアクセス
2. 「Jira同期」タブをクリック
3. 「今すぐ同期」ボタンをクリック

### 3. 同期結果の確認

```bash
# 同期ログを確認
curl http://localhost:8080/api/v1/sync/logs

# 最新の同期ログ
curl http://localhost:8080/api/v1/sync/logs/latest

# プロジェクト一覧を確認
curl http://localhost:8080/api/v1/projects?with_stats=true
```

---

## Docker Composeでの起動

すべてのサービスを一度に起動する場合：

```bash
# サービスの起動
make up

# ログの確認
make logs

# サービスの停止
make down
```

---

## トラブルシューティング

### データベース接続エラー

**症状:**
```
failed to connect to database: connection refused
```

**解決方法:**
```bash
# PostgreSQLコンテナが起動しているか確認
docker ps | grep postgres

# 起動していない場合は再起動
make db-down
make db-up

# 数秒待ってから再試行
sleep 5
make db-migrate
```

### マイグレーションエラー

**症状:**
```
Dirty database version 1. Fix and force version.
```

**解決方法:**
```bash
# 現在のバージョンを確認
make db-version

# 強制的にバージョンを設定
make db-force V=1

# マイグレーションを再実行
make db-migrate
```

### Jira API接続エラー

**症状:**
```
failed to fetch Jira projects: jira API error (status 401)
```

**解決方法:**
1. `JIRA_BASE_URL`が正しいか確認（末尾にスラッシュ不要）
2. `JIRA_EMAIL`が正しいか確認
3. `JIRA_API_TOKEN`が有効か確認（再生成してみる）
4. Jiraアカウントに適切な権限があるか確認

### ポート競合エラー

**症状:**
```
bind: address already in use
```

**解決方法:**
```bash
# ポート8080を使用しているプロセスを確認
lsof -i :8080

# プロセスを終了
kill <PID>

# または環境変数で別のポートを指定
export SERVER_PORT=8081
```

### フロントエンドのビルドエラー

**症状:**
```
error: Cannot find module '@mui/material'
```

**解決方法:**
```bash
# node_modulesを削除して再インストール
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### 同期が遅い

**原因:**
大量のプロジェクト/Issueがある場合、初回同期に時間がかかります。

**対処:**
1. バックエンドログで進捗を確認：
   ```bash
   make logs-backend
   ```

2. 特定のプロジェクトのみ同期：
   ```bash
   curl -X POST http://localhost:8080/api/v1/sync/projects/1
   ```

3. スケジューラーの間隔を調整：
   ```bash
   # 2時間ごとに変更
   cd backend
   go run cmd/sync/main.go -mode=scheduler -interval=2h
   ```

---

## 次のステップ

1. [API ドキュメント](./API.md)を確認
2. 組織階層を作成（UI: `/admin`）
3. プロジェクトを組織に割り当て
4. ダッシュボードで可視化を確認（UI: `/`）

---

## よくある質問

### Q: 本番環境へのデプロイ方法は？

A: 現在、本番デプロイ用の設定は未実装です。以下を実装予定：
- Docker マルチステージビルド
- CI/CD パイプライン
- 環境別設定管理

### Q: 認証機能はありますか？

A: 現在は未実装です。将来的にJWT認証を追加予定です。

### Q: データベースのバックアップ方法は？

A: PostgreSQLの標準ツールを使用：
```bash
# バックアップ
docker exec project-viz-db pg_dump -U admin project_visualization > backup.sql

# リストア
docker exec -i project-viz-db psql -U admin project_visualization < backup.sql
```

### Q: 複数の組織を管理できますか？

A: はい。組織は階層構造で無制限に作成できます。UIまたはAPIから管理可能です。
