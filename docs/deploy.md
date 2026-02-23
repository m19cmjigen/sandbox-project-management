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

## CI/CDパイプライン

GitHub Actionsを使用してCI/CDを自動化しています。ワークフローファイルは `.github/workflows/` に配置されています。

### ワークフロー構成

| ファイル | トリガー | 内容 |
|---|---|---|
| `ci.yml` | PR作成・mainへのpush | バックエンドテスト・フロントエンドテスト・Dockerビルド確認 |
| `cd.yml` | mainへのpush (staging) / 手動実行 (production) | ECRプッシュ・ECSデプロイ・S3+CloudFrontデプロイ |

### CIの流れ

```
PR作成 / push
    │
    ├── Backend job: go build → go test → カバレッジ計測
    ├── Frontend job: type-check → lint → vitest → npm build
    ├── E2E job: (mainへのpushのみ) Playwright実行
    └── Docker Build Check: イメージビルド確認
```

### CDの流れ

```
mainへのpush
    │
    └── Deploy Staging:
            │
            ├── ECRへDockerイメージをpush
            ├── ECS Run Taskでマイグレーション実行
            ├── ECSサービスを新イメージで更新 (ローリングアップデート)
            ├── フロントエンドをビルドしてS3へsync
            └── CloudFrontキャッシュを無効化

手動実行 (workflow_dispatch, environment=production)
    │
    └── Deploy Production: (GitHub Environments の承認者が承認後に実行)
            │
            └── 上記と同様の手順 (本番用AWSクレデンシャル使用)
```

### 必要なGitHub Secrets

CDパイプラインを動作させるには以下のSecretsをリポジトリに設定してください。

**Staging用:**

| Secret名 | 説明 |
|---|---|
| `AWS_ACCESS_KEY_ID` | stagingデプロイ用IAMアクセスキー |
| `AWS_SECRET_ACCESS_KEY` | stagingデプロイ用IAMシークレットキー |
| `SUBNET_IDS_STAGING` | マイグレーションタスク実行用サブネットID |
| `SG_ID_STAGING` | マイグレーションタスク実行用セキュリティグループID |
| `CLOUDFRONT_DISTRIBUTION_STAGING` | StagingのCloudFrontディストリビューションID |
| `STAGING_API_BASE_URL` | StagingのAPIベースURL |
| `JIRA_BASE_URL` | JiraのベースURL（例: `https://your-company.atlassian.net`）|

**Production用:**

| Secret名 | 説明 |
|---|---|
| `AWS_ACCESS_KEY_ID_PROD` | 本番デプロイ用IAMアクセスキー |
| `AWS_SECRET_ACCESS_KEY_PROD` | 本番デプロイ用IAMシークレットキー |
| `SUBNET_IDS_PROD` | 本番マイグレーションタスク実行用サブネットID |
| `SG_ID_PROD` | 本番マイグレーションタスク実行用セキュリティグループID |
| `CLOUDFRONT_DISTRIBUTION_PRODUCTION` | 本番CloudFrontディストリビューションID |
| `PRODUCTION_API_BASE_URL` | 本番APIベースURL |

### GitHub Environmentsの設定（本番承認フロー）

1. GitHubリポジトリの **Settings → Environments** を開く
2. `production` Environmentを作成
3. **Required reviewers** に承認者を追加する
4. `workflow_dispatch` でデプロイを実行すると、承認者への承認依頼が送信される

## ロールバック手順

デプロイ後に問題が発生した場合のロールバック手順です。

### バックエンドのロールバック (ECS)

**方法1: 前のタスク定義に戻す（推奨）**

```bash
# 現在のタスク定義リビジョンを確認
aws ecs describe-services \
  --cluster project-viz \
  --services project-viz-api-staging \
  --query 'services[0].taskDefinition'

# 前のリビジョンに戻す（例: リビジョン番号を1つ前に）
aws ecs update-service \
  --cluster project-viz \
  --service project-viz-api-staging \
  --task-definition project-viz-api-staging:<PREV_REVISION> \
  --force-new-deployment
```

**方法2: 前のDockerイメージを再デプロイ**

1. GitHub Actionsの **CD** ワークフローを開く
2. 正常だったコミットハッシュを確認する
3. `workflow_dispatch` でそのコミットをターゲットにして再実行する

### DBマイグレーションのロールバック

```bash
# 本番DBへの接続（AWS SSM Session Manager経由推奨）
DATABASE_URL="postgres://<user>:<pass>@<aurora-endpoint>:5432/<db>?sslmode=require" \
  migrate -path database/migrations -database "$DATABASE_URL" down 1

# バージョン確認
DATABASE_URL="postgres://..." \
  migrate -path database/migrations -database "$DATABASE_URL" version
```

**注意**: DBマイグレーションのロールバックはデータ損失リスクがあります。必ずバックアップを取ってから実施してください。

### フロントエンドのロールバック (S3 + CloudFront)

```bash
# 前のビルド成果物で上書き（GitHub ActionsのArtifactからダウンロードして再デプロイ）
aws s3 sync <前のビルドディレクトリ>/ s3://project-viz-frontend-staging/ --delete

# CloudFrontキャッシュを無効化
aws cloudfront create-invalidation \
  --distribution-id <DISTRIBUTION_ID> \
  --paths "/*"
```

### ローカル開発環境のロールバック

```bash
# 前のコミットのイメージをビルドして起動
git checkout <prev-commit>
docker compose build backend
docker compose up -d backend
```
