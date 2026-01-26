# 全社プロジェクト進捗可視化プラットフォーム

組織ごとに分散しているJiraプロジェクトの進捗状況（特に納期遅延）を一元管理し、経営層・PMO・管理職が早期に対策を打てる状態にするプラットフォーム。

## プロジェクト概要

### ターゲットユーザー
- 経営層
- PMO（プロジェクトマネジメントオフィス）
- 部門長
- プロジェクトマネージャー

### 主な機能
1. **ダッシュボード**: 全社・組織別のプロジェクト遅延状況の可視化
2. **組織管理**: 組織階層の管理とプロジェクトの紐付け
3. **プロジェクト一覧**: 遅延プロジェクトの確認
4. **チケット詳細**: 遅延チケットのフィルタリングと詳細確認
5. **JWT認証**: ロールベースアクセス制御（Admin/Manager/Viewer）
6. **監査ログ**: 全ユーザーアクションの記録とトレーサビリティ
7. **Jira同期**: 自動・手動でのJira Cloudデータ同期
8. **管理機能**: 組織・ユーザー・同期の統合管理画面

## 技術スタック

### Frontend
- **React** 18.2 + **TypeScript** 5.3
- **Material-UI** (MUI) 5.15
- **React Router** 6.22
- **Zustand** (状態管理)
- **Vite** (ビルドツール)

### Backend
- **Go** 1.21+
- **Gin** (Webフレームワーク)
- **sqlx** (データベースアクセス)
- **PostgreSQL** 15+

### Infrastructure
- **AWS Fargate** (ECS)
- **Amazon Aurora** (PostgreSQL)
- **CloudFront** + **S3** (Frontend配信)
- **Terraform** (IaC)

### Batch Processing
- **Go** (Goroutineによる並行処理)
- **Jira Cloud API** 連携

## プロジェクト構造

```
.
├── backend/              # Go バックエンドAPI
│   ├── cmd/
│   │   └── api/         # エントリーポイント
│   ├── internal/        # 内部パッケージ
│   │   ├── domain/      # ドメイン層
│   │   ├── usecase/     # ユースケース層
│   │   ├── interface/   # インターフェース層
│   │   └── infrastructure/ # インフラ層
│   └── pkg/             # 共通パッケージ
│
├── frontend/            # React フロントエンド
│   └── src/
│       ├── components/  # 共通コンポーネント
│       ├── pages/       # ページコンポーネント
│       ├── hooks/       # カスタムフック
│       ├── services/    # API通信
│       └── store/       # 状態管理
│
├── database/            # データベース管理
│   ├── schema/          # スキーマ設計
│   └── migrations/      # マイグレーションファイル
│
├── tickets/             # 開発チケット (37件)
│
├── SPEC.md              # 要件定義書
├── docker-compose.yml   # ローカル開発環境
└── Makefile            # 開発コマンド
```

## セットアップ

### 前提条件
- Docker & Docker Compose
- Go 1.21+
- Node.js 18.0+
- PostgreSQL 15+ (またはDocker)
- golang-migrate

### クイックスタート

1. **リポジトリのクローン**
```bash
git clone https://github.com/m19cmjigen/sandbox-project-management.git
cd sandbox-project-management
```

2. **Docker Composeで起動**
```bash
# すべてのサービスを起動（PostgreSQL + Backend）
docker-compose up -d

# ログ確認
docker-compose logs -f
```

3. **データベースマイグレーション**
```bash
make db-migrate
```

4. **フロントエンドの起動**
```bash
cd frontend
npm install
npm run dev
```

アプリケーションが以下のURLで起動します：
- フロントエンド: http://localhost:5173
- バックエンドAPI: http://localhost:8080

### 初回ログイン

デフォルトの管理者アカウント：
- **ユーザー名**: `admin`
- **パスワード**: `admin123`

⚠️ **重要**: 本番環境では必ずパスワードを変更してください！

📘 **詳しい使い方**: [docs/QUICK_START.md](docs/QUICK_START.md) を参照

### 個別セットアップ

#### データベース

```bash
# PostgreSQL起動
make db-up

# マイグレーション適用
make db-migrate

# 接続確認
make db-connect
```

詳細: [database/README.md](database/README.md)

#### バックエンド

```bash
cd backend

# 依存関係のインストール
go mod download

# 環境変数設定
cp .env.example .env

# 起動
make backend-run
```

詳細: [backend/README.md](backend/README.md)

#### フロントエンド

```bash
cd frontend

# 依存関係のインストール
npm install

# 起動
npm run dev
```

詳細: [frontend/README.md](frontend/README.md)

## 開発コマンド

### データベース
```bash
make db-up          # PostgreSQL起動
make db-down        # PostgreSQL停止
make db-migrate     # マイグレーション適用
make db-rollback    # マイグレーションロールバック
make db-version     # マイグレーションバージョン確認
```

### バックエンド
```bash
make backend-run    # バックエンド起動
make backend-build  # ビルド
make backend-test   # テスト実行
make backend-lint   # リント実行
make backend-fmt    # フォーマット
```

### Docker Compose
```bash
make up             # すべてのサービス起動
make down           # すべてのサービス停止
make logs           # ログ表示
make logs-backend   # バックエンドログ
make logs-db        # データベースログ
```

## 開発チケット

プロジェクトは37個のチケットに分割されています。

チケット一覧とステータスは [tickets/README.md](tickets/README.md) を参照してください。

### カテゴリ別チケット数
- Infrastructure: 5チケット (13人日)
- Database: 4チケット (8人日)
- Backend: 6チケット (24人日)
- Batch Worker: 5チケット (19人日)
- Frontend: 7チケット (31人日)
- Security: 3チケット (11人日)
- Testing: 3チケット (22人日)
- Deployment: 2チケット (8人日)
- Documentation: 2チケット (8人日)

**合計見積もり**: 144人日

### 開発フェーズ

1. **Phase 1**: 基盤構築 (DB, Backend, Frontend セットアップ) ✅
2. **Phase 2**: コア機能開発 (API, Batch, UI実装)
3. **Phase 3**: 管理機能・監視
4. **Phase 4**: テスト・品質保証
5. **Phase 5**: デプロイ・ドキュメント

## アーキテクチャ

### システム構成

```
[User Browser] --(HTTPS)--> [CloudFront + S3 (React App)]
       |
       +--(API Request)--> [ALB / API Gateway]
                                |
                          [Backend API (Go)] <--> [Aurora DB]
                                |
                          [Batch Worker (Go)] --(REST API)--> [Jira Cloud]
```

### データフロー

1. **バッチ処理**: Jira Cloud APIからデータ取得 → 正規化 → DB保存
2. **API**: クライアントリクエスト → DB照会 → 集計 → レスポンス
3. **UI**: React SPA → API呼び出し → データ表示

## データベーススキーマ

主要テーブル：
- `organizations`: 組織階層マスタ（トリガー、ビューあり）
- `projects`: Jiraプロジェクト情報
- `issues`: Jiraチケット情報
- `sync_logs`: バッチ実行ログ
- `users`: ユーザーアカウント（JWT認証用）
- `audit_logs`: 監査ログ（全ユーザーアクション記録）

マイグレーション：
- 010個のマイグレーションファイル（up/down）
- golang-migrateでバージョン管理

詳細: [database/schema/schema_design.md](database/schema/schema_design.md)

## API エンドポイント

### 認証（公開エンドポイント）
- `POST /api/v1/auth/login` - ログイン
- `POST /api/v1/auth/refresh` - トークン更新

### 組織管理（認証必須）
- `GET /api/v1/organizations` - 組織一覧
- `POST /api/v1/organizations` - 組織作成（マネージャー以上）
- `PUT /api/v1/organizations/:id` - 組織更新（マネージャー以上）
- `DELETE /api/v1/organizations/:id` - 組織削除（マネージャー以上）

### プロジェクト管理（認証必須）
- `GET /api/v1/projects` - プロジェクト一覧
- `GET /api/v1/projects/:id` - プロジェクト詳細
- `PUT /api/v1/projects/:id/organization` - 組織紐付け（マネージャー以上）

### チケット管理（認証必須）
- `GET /api/v1/issues` - チケット一覧
- `GET /api/v1/issues/:id` - チケット詳細

### ダッシュボード（認証必須）
- `GET /api/v1/dashboard/summary` - 全社サマリ
- `GET /api/v1/dashboard/organizations/:id` - 組織別サマリ
- `GET /api/v1/dashboard/projects/:id` - プロジェクト別サマリ

### Jira同期（マネージャー以上）
- `POST /api/v1/sync/trigger` - 手動同期実行
- `GET /api/v1/sync/logs` - 同期履歴
- `GET /api/v1/sync/logs/latest` - 最新同期ログ

### ユーザー管理（管理者のみ）
- `GET /api/v1/users` - ユーザー一覧
- `POST /api/v1/users` - ユーザー作成
- `PUT /api/v1/users/:id` - ユーザー更新
- `DELETE /api/v1/users/:id` - ユーザー削除

### 監査ログ（管理者のみ）
- `GET /api/v1/audit/logs` - 監査ログ一覧
- `GET /api/v1/audit/logs/:id` - 監査ログ詳細

**詳細**: [docs/API.md](docs/API.md) を参照

## テスト

### バックエンド単体テスト

```bash
cd backend

# テスト実行
go test -v ./...

# カバレッジ付き
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Makefileから
make backend-test
make backend-coverage
```

### フロントエンドE2Eテスト（Playwright）

```bash
cd frontend

# E2Eテスト実行
npm run test:e2e

# UIモードで実行（推奨）
npm run test:e2e:ui

# ヘッドモードで実行
npm run test:e2e:headed

# デバッグモード
npm run test:e2e:debug

# テストレポート表示
npm run test:e2e:report
```

**詳細**: [frontend/e2e/README.md](frontend/e2e/README.md)

### CI/CDでのテスト

GitHub Actionsで自動実行：
- バックエンド単体テスト
- golangci-lint静的解析
- フロントエンドE2Eテスト
- Dockerビルドテスト

## デプロイ

### 本番環境デプロイ

```bash
# インフラ構築 (Terraform)
cd terraform
terraform init
terraform plan
terraform apply

# バックエンドデプロイ
cd backend
make docker-push
make deploy

# フロントエンドデプロイ
cd frontend
npm run build
aws s3 sync dist/ s3://your-bucket/
aws cloudfront create-invalidation --distribution-id XXX --paths "/*"
```

## 環境変数

### Backend (.env)
```bash
# Server
PORT=8080
GIN_MODE=release

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=project_visualization
DB_SSLMODE=disable
DATABASE_URL=postgres://admin:admin123@localhost:5432/project_visualization?sslmode=disable

# JWT Authentication
JWT_SECRET_KEY=your-secret-key-change-in-production
JWT_EXPIRATION_HOURS=24

# Jira Integration
JIRA_BASE_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-jira-api-token

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Batch Job
SYNC_INTERVAL=1h
DEFAULT_ORGANIZATION_ID=1
```

### Frontend
API通信は `/api` プレフィックスを使用し、Viteプロキシで `http://localhost:8080` に転送。

**重要**: `.env.example`をコピーして`.env`を作成し、適切な値を設定してください。

## ドキュメント

### 仕様・設計

- [SPEC.md](SPEC.md) - プロジェクト要件定義書
- [database/schema/schema_design.md](database/schema/schema_design.md) - データベーススキーマ設計書
- [docs/API.md](docs/API.md) - API仕様書

### セットアップ・デプロイ

- [docs/SETUP.md](docs/SETUP.md) - 環境セットアップガイド
- [docs/DEPLOY.md](docs/DEPLOY.md) - 本番環境デプロイガイド
- [database/README.md](database/README.md) - データベース管理
- [backend/README.md](backend/README.md) - バックエンド開発ガイド
- [frontend/README.md](frontend/README.md) - フロントエンド開発ガイド

### ユーザー向け

- 📘 [docs/QUICK_START.md](docs/QUICK_START.md) - **5分でわかる使い方**（初めての方はこちら）
- 📖 [docs/USER_MANUAL.md](docs/USER_MANUAL.md) - **ユーザーマニュアル**（完全版）

### 開発・テスト

- [tickets/README.md](tickets/README.md) - 開発チケット一覧
- [frontend/e2e/README.md](frontend/e2e/README.md) - E2Eテストガイド（Playwright）
- [docs/README-ja.md](docs/README-ja.md) - 日本語README

## トラブルシューティング

### データベース接続エラー
```bash
# PostgreSQLの起動確認
docker ps | grep postgres

# ログ確認
docker logs project-viz-db
```

### ポート競合
```bash
# 使用中のポートを確認
lsof -i :8080
lsof -i :3000

# プロセス終了
kill -9 <PID>
```

## ライセンス

このプロジェクトは社内プロジェクトです。

## コントリビューション

開発に参加する場合は、以下の手順に従ってください：

1. 新しいブランチを作成
2. 変更をコミット
3. プルリクエストを作成
4. レビュー後にマージ

## サポート

質問や問題がある場合は、Issue を作成してください。
