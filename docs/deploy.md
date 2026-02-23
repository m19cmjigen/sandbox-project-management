# デプロイ手順書

## 概要

本システムはDockerを使用してデプロイします。開発環境・本番環境ともにDockerコンテナで動作します。

## 前提条件

- Docker 20.10以上
- Docker Compose 2.0以上
- make
- golang-migrate (`brew install golang-migrate`)

## ローカル開発環境構築

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd sandbox-project-management
```

### 2. 環境変数ファイルの作成

```bash
cd backend && cp .env.example .env
```

`.env`の主要変数:

```
PORT=8080
GIN_MODE=debug
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=project_visualization
DB_SSLMODE=disable
LOG_LEVEL=debug
LOG_FORMAT=json
```

### 3. サービスの起動

```bash
# PostgreSQL + バックエンドをDockerで起動
make up

# マイグレーション適用
make db-migrate

# フロントエンド開発サーバー起動
cd frontend && npm install && npm run dev
```

### 4. 動作確認

```bash
# バックエンドのヘルスチェック
curl http://localhost:8080/health

# フロントエンド
open http://localhost:3000
```

## バックエンドのビルド

ローカルでバイナリをビルドする場合:

```bash
make backend-build
# backend/bin/api にバイナリが生成される
```

Dockerイメージのビルド:

```bash
docker compose build backend
```

## データベースマイグレーション

### 新規マイグレーション適用

```bash
make db-migrate
```

### ロールバック

```bash
# 1ステップロールバック
make db-rollback

# 全ロールバック
make db-rollback-all
```

### バージョン確認

```bash
make db-version
```

### エラー時のリカバリ

マイグレーションがエラー状態 (dirty) になった場合:

```bash
# バージョンを直前の正常なバージョンに強制設定
make db-force V=<version>
```

## 本番環境デプロイ (想定構成)

本番環境ではAWS上での以下の構成を想定しています。

### インフラ構成

```
Route 53 → ALB → ECS Fargate (Backend)
                       ↓
              Amazon Aurora PostgreSQL
```

### デプロイフロー

1. Dockerイメージのビルド・ECRへのプッシュ

```bash
# ECRへのログイン
aws ecr get-login-password --region ap-northeast-1 | \
  docker login --username AWS --password-stdin <account-id>.dkr.ecr.ap-northeast-1.amazonaws.com

# ビルド & プッシュ
docker build -t project-viz-api ./backend
docker tag project-viz-api:latest <ecr-uri>/project-viz-api:latest
docker push <ecr-uri>/project-viz-api:latest
```

2. マイグレーションの適用

```bash
# ECS Task (migration) をワンショット実行
# または接続可能な環境からmigrate CLIを実行
DATABASE_URL="postgres://<user>:<pass>@<aurora-endpoint>:5432/<db>?sslmode=require" \
  migrate -path database/migrations -database "$DATABASE_URL" up
```

3. ECSサービスの更新

```bash
aws ecs update-service \
  --cluster project-viz \
  --service project-viz-api \
  --force-new-deployment
```

4. フロントエンドのデプロイ

```bash
# ビルド
cd frontend && npm run build

# S3 + CloudFrontへアップロード
aws s3 sync dist/ s3://<bucket-name>/ --delete
aws cloudfront create-invalidation \
  --distribution-id <distribution-id> \
  --paths "/*"
```

### 環境変数 (本番)

| 変数 | 説明 |
|---|---|
| PORT | APIポート (8080) |
| GIN_MODE | `release` に設定 |
| DB_HOST | Aurora エンドポイント |
| DB_PORT | 5432 |
| DB_USER | DBユーザー名 |
| DB_PASSWORD | DBパスワード (Secrets Manager推奨) |
| DB_NAME | データベース名 |
| DB_SSLMODE | `require` |
| LOG_LEVEL | `info` |
| LOG_FORMAT | `json` |

### セキュリティ考慮事項

- DBパスワードはAWS Secrets Managerで管理する
- VPC内にECS・RDSを配置し、外部からの直接アクセスを禁止する
- ALBにHTTPS (ACM証明書) を設定する
- GIN_MODEは本番で必ず`release`にする (デバッグ情報の漏洩防止)

## Docker Composeコマンド一覧

```bash
make up           # 全サービスを起動
make down         # 全サービスを停止
make logs         # 全ログを表示 (フォロー)
make logs-backend # バックエンドのログを表示
make logs-db      # PostgreSQLのログを表示
```
