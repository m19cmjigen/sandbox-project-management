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
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080

### Jira連携のセットアップ

Jiraからプロジェクト・チケットを取得するバッチ処理を使う場合は、追加でJira認証情報の設定が必要です。

```bash
# 1. APIトークンを https://id.atlassian.com/manage-profile/security/api-tokens で発行

# 2. .env に認証情報を追加
cd backend && cp .env.example .env
# .env の JIRA_BASE_URL / JIRA_EMAIL / JIRA_API_TOKEN を設定

# 3. バッチを実行
go build -o bin/batch ./cmd/batch/ && ./bin/batch
```

詳細な手順（権限設定・確認方法・トラブルシューティング）: [docs/jira-setup.md](docs/jira-setup.md)

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
- `organizations`: 組織階層マスタ
- `projects`: Jiraプロジェクト情報
- `issues`: Jiraチケット情報
- `sync_logs`: バッチ実行ログ

詳細: [database/schema/schema_design.md](database/schema/schema_design.md)

## API エンドポイント

### 組織管理
- `GET /api/v1/organizations` - 組織一覧
- `POST /api/v1/organizations` - 組織作成
- `PUT /api/v1/organizations/:id` - 組織更新
- `DELETE /api/v1/organizations/:id` - 組織削除

### プロジェクト管理
- `GET /api/v1/projects` - プロジェクト一覧
- `PUT /api/v1/projects/:id/organization` - 組織紐付け

### ダッシュボード
- `GET /api/v1/dashboard/summary` - 全社サマリ
- `GET /api/v1/dashboard/organizations/:id` - 組織別サマリ

## テスト

```bash
# バックエンド
cd backend
make backend-test
make backend-coverage

# フロントエンド
cd frontend
npm run test
npm run test:coverage
```

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

```env
# サーバー
PORT=8080
GIN_MODE=debug

# データベース
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=project_visualization
DB_SSLMODE=disable

# Jira連携（バッチ処理に必要）
JIRA_BASE_URL=https://your-org.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-api-token-here

# バッチ設定
BATCH_SYNC_MODE=full       # full または delta
BATCH_WORKER_COUNT=5       # 並列フェッチ数
```

テンプレート: `backend/.env.example`
Jira APIトークンの取得方法: [docs/jira-setup.md](docs/jira-setup.md)

### Frontend
API通信は `/api` プレフィックスを使用し、Viteプロキシで `http://localhost:8080` に転送。

## ドキュメント

- [SPEC.md](SPEC.md) - プロジェクト要件定義書
- [docs/jira-setup.md](docs/jira-setup.md) - **Jira連携セットアップガイド**
- [docs/secrets-management.md](docs/secrets-management.md) - Jira認証情報のセキュア管理（本番環境）
- [docs/batch-schedule.md](docs/batch-schedule.md) - バッチスケジュール設定
- [docs/deploy.md](docs/deploy.md) - デプロイ手順書
- [docs/operations.md](docs/operations.md) - 運用手順書
- [database/README.md](database/README.md) - データベース管理
- [database/schema/schema_design.md](database/schema/schema_design.md) - スキーマ設計書
- [tickets/README.md](tickets/README.md) - 開発チケット一覧

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
