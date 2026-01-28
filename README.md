# 全社プロジェクト進捗可視化プラットフォーム

組織ごとに分散しているJiraプロジェクトの進捗状況（特に納期遅延）を一元管理し、経営層・PMO・管理職が早期に対策を打てる状態にするプラットフォーム。

## プロジェクト概要

### ターゲットユーザー
- 経営層
- PMO（プロジェクトマネジメントオフィス）
- 部門長
- プロジェクトマネージャー

### 主な機能

#### 📊 可視化・ダッシュボード
1. **ダッシュボード**: 全社・組織別・プロジェクト別の遅延状況を可視化
   - 🔴 RED (遅延): 期限超過
   - 🟡 YELLOW (注意): 期限間近または期限未設定
   - 🟢 GREEN (正常): 余裕あり
2. **ヒートマップ表示**: プロジェクト遅延状況の視覚的把握

#### 🏢 組織・プロジェクト管理
3. **組織階層管理**: ツリー表示、階層検索、プロジェクト紐付け
4. **プロジェクト管理**: Jira連携、組織割り当て、統計表示
5. **チケット管理**: 遅延フィルタリング、詳細表示

#### 🔒 セキュリティ
6. **JWT認証**: ロールベースアクセス制御（Admin/Manager/Viewer）
7. **監査ログ**: 全APIリクエストの記録とトレーサビリティ
8. **パスワード暗号化**: bcryptによる安全なパスワード管理

#### 🔄 自動化・同期
9. **Jira Cloud同期**: Full/Delta同期、スケジューラー、手動実行
10. **バッチ処理**: リトライ機構、エラーハンドリング、進捗ログ

#### 👥 管理機能
11. **ユーザー管理**: ロール設定、アクティブ/非アクティブ管理
12. **同期ログ管理**: 履歴確認、エラー追跡

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

# シードデータ投入（開発・テスト用）
make db-seed

# 接続確認
make db-connect
```

**シードデータ**: 開発環境用のサンプルデータ（4ユーザー、9組織、5プロジェクト、30-50チケット）を簡単に投入できます。詳細は [database/README_SEED.md](database/README_SEED.md) を参照してください。

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
make backend-run      # バックエンド起動
make backend-build    # ビルド
make backend-test     # テスト実行
make backend-lint     # リント実行
make backend-fmt      # フォーマット
make backend-coverage # テストカバレッジ生成
```

### パフォーマンステスト
```bash
make perf-smoke       # スモークテスト（疎通確認）
make perf-load        # 負荷テスト（50-100 VU）
make perf-stress      # ストレステスト（最大200 VU）
make perf-large-data  # 大量データ生成（10,000+ issues）
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
- **Infrastructure**: 5チケット (13人日) - *AWS環境構築（環境依存のため保留）*
- **Database**: 4チケット (8人日) - ✅ **完了**
- **Backend**: 6チケット (24人日) - ✅ **完了**
- **Batch Worker**: 5チケット (19人日) - ✅ **完了**
- **Frontend**: 7チケット (31人日) - ✅ **完了**
- **Security**: 3チケット (11人日) - ✅ **完了**
- **Testing**: 3チケット (22人日) - ⚙️ **進行中**（ドメイン・リポジトリ層完了、E2E完了、パフォーマンステスト完了）
- **Deployment**: 2チケット (8人日) - ⚙️ **部分完了**（CI/CD完了、本番デプロイは環境依存）
- **Documentation**: 2チケット (8人日) - ✅ **完了**

**合計見積もり**: 144人日
**完了率**: 約85%（実装可能な範囲では95%完了）

### 開発フェーズ

1. **Phase 1**: 基盤構築 (DB, Backend, Frontend セットアップ) ✅ **完了**
2. **Phase 2**: コア機能開発 (API, Batch, UI実装) ✅ **完了**
3. **Phase 3**: 管理機能・監視 ✅ **完了**
4. **Phase 4**: テスト・品質保証 ⚙️ **進行中**（85%完了）
5. **Phase 5**: デプロイ・ドキュメント ✅ **完了**（CI/CD完了、ドキュメント完備）

### 実装済み機能

#### ✅ コア機能
- ダッシュボード（全社・組織別・プロジェクト別）
- 組織階層管理（トリガー、ビューあり）
- プロジェクト管理（Jira連携）
- チケット管理（遅延ステータス可視化）
- Jira Cloud API同期（Full/Delta）

#### ✅ セキュリティ
- JWT認証・認可（ロールベース: Admin/Manager/Viewer）
- パスワードハッシュ化（bcrypt）
- 監査ログ（全APIリクエスト記録）
- セキュアなAPI設計

#### ✅ テスト・品質
- ドメイン層ユニットテスト（100%カバレッジ）
- リポジトリ層統合テスト（20+テストケース）
- E2Eテスト（Playwright、6スイート）
- パフォーマンステスト（k6、3シナリオ）
- CI/CD自動テスト

#### ✅ 開発体験
- シードデータスクリプト（開発環境すぐ構築）
- 大量データ生成（10,000+チケット）
- Makefile統合（40+コマンド）
- 包括的ドキュメント（日本語・英語）

## アーキテクチャ

### Clean Architecture（レイヤードアーキテクチャ）

```
┌─────────────────────────────────────────────────────────────┐
│                        Presentation Layer                    │
│  ┌────────────────┐              ┌──────────────────────┐   │
│  │  React UI      │              │   REST API (Gin)     │   │
│  │  - Components  │◄────────────►│   - Handlers         │   │
│  │  - Pages       │   HTTP/JSON  │   - Middleware       │   │
│  └────────────────┘              └──────────────────────┘   │
└────────────────────────────────────────┬────────────────────┘
                                         │
┌────────────────────────────────────────┴────────────────────┐
│                        Use Case Layer                        │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Business Logic                                      │   │
│  │  - OrganizationUsecase                              │   │
│  │  - ProjectUsecase                                   │   │
│  │  - IssueUsecase                                     │   │
│  │  - DashboardUsecase                                 │   │
│  │  - AuthUsecase                                      │   │
│  │  - SyncUsecase (Jira同期)                          │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────────────────────┬────────────────────┘
                                         │
┌────────────────────────────────────────┴────────────────────┐
│                        Domain Layer                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Entities & Value Objects                           │   │
│  │  - Organization, Project, Issue                     │   │
│  │  - User, SyncLog, AuditLog                         │   │
│  │  - DelayStatus, StatusCategory                     │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────────────────────┬────────────────────┘
                                         │
┌────────────────────────────────────────┴────────────────────┐
│                   Infrastructure Layer                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │  PostgreSQL  │  │  Jira API    │  │  JWT / bcrypt    │  │
│  │  Repository  │  │  Client      │  │  Auth Service    │  │
│  └──────────────┘  └──────────────┘  └──────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### システム構成（本番環境）

```
┌──────────────┐
│   User       │
│   Browser    │
└──────┬───────┘
       │ HTTPS
       ↓
┌──────────────────────────────────────────────────────┐
│              CloudFront + S3                         │
│              (React SPA)                             │
└──────┬───────────────────────────────────────────────┘
       │ API Request
       ↓
┌──────────────────────────────────────────────────────┐
│              ALB / API Gateway                       │
└──────┬───────────────────────────────────────────────┘
       │
       ↓
┌──────────────────────────────────────────────────────┐
│         Backend API (Go on Fargate/ECS)              │
│  ┌────────────┐  ┌────────────┐  ┌──────────────┐   │
│  │   API      │  │   Batch    │  │    Auth      │   │
│  │  Handler   │  │   Worker   │  │   Service    │   │
│  └────────────┘  └────────────┘  └──────────────┘   │
└──────┬─────────────────┬───────────────────────┬─────┘
       │                 │                       │
       ↓                 ↓                       ↓
┌──────────────┐  ┌────────────────┐  ┌────────────────┐
│  Aurora DB   │  │  Jira Cloud    │  │  CloudWatch    │
│ (PostgreSQL) │  │  REST API      │  │  Logs          │
└──────────────┘  └────────────────┘  └────────────────┘
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

**テストカバレッジ**:
- ドメイン層: 100% ✅
- リポジトリ層: 統合テスト完備（20+テストケース）✅
- E2Eテスト: Playwright（6テストスイート）✅

**テスト内容**:
- `internal/domain/*_test.go` - ドメインロジックのユニットテスト
- `internal/infrastructure/postgres/*_test.go` - リポジトリ層の統合テスト
- データベース接続が利用できない場合は自動スキップ

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

**テストスイート**（6スイート）:
- 認証フロー（Authentication Flow）
- ダッシュボード（Dashboard）
- プロジェクト管理（Projects）
- チケット管理（Issues）
- 組織管理（Organizations）
- 管理者機能（Admin）

**詳細**: [frontend/e2e/README.md](frontend/e2e/README.md)

### パフォーマンステスト（k6）

```bash
# k6のインストール（macOS）
brew install k6

# スモークテスト（疎通確認）
make perf-smoke

# 負荷テスト（本番想定）
make perf-load

# ストレステスト（限界確認）
make perf-stress
```

**パフォーマンス目標**:
- Dashboard API: p95 < 500ms
- Project List: p95 < 300ms
- Issue Search: p95 < 500ms
- Error Rate: < 1%

**詳細**: [performance/README.md](performance/README.md)

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
- [performance/README.md](performance/README.md) - パフォーマンステストガイド（k6）
- [database/README_SEED.md](database/README_SEED.md) - シードデータガイド
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

## 機能一覧

### ダッシュボード機能

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| 全社サマリ | 全プロジェクトの遅延状況を集計表示 | ✅ 完了 |
| 組織別サマリ | 組織ごとの遅延状況を表示 | ✅ 完了 |
| プロジェクト別サマリ | プロジェクトごとの詳細統計 | ✅ 完了 |
| ヒートマップ | 遅延状況を色分けで視覚化 | ✅ 完了 |
| リアルタイム更新 | 最新データを自動取得 | ✅ 完了 |

### 組織管理機能

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| 組織階層表示 | ツリービューで階層構造を表示 | ✅ 完了 |
| 組織CRUD | 作成・更新・削除 | ✅ 完了 |
| プロジェクト紐付け | 組織へのプロジェクト割り当て | ✅ 完了 |
| 階層検索 | パスベースの検索機能 | ✅ 完了 |

### プロジェクト管理機能

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| プロジェクト一覧 | 全プロジェクトの表示 | ✅ 完了 |
| プロジェクト詳細 | 統計・チケット一覧 | ✅ 完了 |
| 組織割り当て | プロジェクトと組織の紐付け | ✅ 完了 |
| フィルタリング | 遅延ステータス別表示 | ✅ 完了 |
| ソート機能 | 各種条件でソート | ✅ 完了 |

### チケット管理機能

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| チケット一覧 | 全チケットの表示 | ✅ 完了 |
| チケット詳細 | 詳細情報の表示 | ✅ 完了 |
| 遅延フィルタ | RED/YELLOW/GREEN別表示 | ✅ 完了 |
| プロジェクトフィルタ | プロジェクト別表示 | ✅ 完了 |
| ページネーション | 大量データの効率的表示 | ✅ 完了 |
| Jiraリンク | Jira画面への直接リンク | ✅ 完了 |

### 認証・認可機能

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| ログイン | JWT認証 | ✅ 完了 |
| ログアウト | セッション終了 | ✅ 完了 |
| トークン更新 | 自動リフレッシュ | ✅ 完了 |
| ロールベース制御 | Admin/Manager/Viewer | ✅ 完了 |
| パスワード変更 | セキュアなパスワード更新 | ✅ 完了 |

### ユーザー管理機能（管理者のみ）

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| ユーザー一覧 | 全ユーザーの表示 | ✅ 完了 |
| ユーザー作成 | 新規ユーザー登録 | ✅ 完了 |
| ユーザー更新 | 情報・ロール変更 | ✅ 完了 |
| ユーザー削除 | 論理削除（非アクティブ化） | ✅ 完了 |
| パスワードリセット | 管理者によるリセット | ✅ 完了 |

### Jira同期機能

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| Full Sync | 全データの同期 | ✅ 完了 |
| Delta Sync | 差分同期 | ✅ 完了 |
| 手動実行 | API経由での同期トリガー | ✅ 完了 |
| スケジューラー | 定期自動実行 | ✅ 完了 |
| リトライ機構 | エラー時の自動リトライ | ✅ 完了 |
| 同期ログ | 実行履歴の記録 | ✅ 完了 |
| エラー通知 | 失敗時のロギング | ✅ 完了 |

### 監査・セキュリティ機能

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| 監査ログ | 全APIリクエストの記録 | ✅ 完了 |
| ログ検索 | フィルタリング・検索 | ✅ 完了 |
| ログ削除 | 古いログの削除 | ✅ 完了 |
| パスワード暗号化 | bcryptハッシュ | ✅ 完了 |
| トークン暗号化 | JWTシークレット管理 | ✅ 完了 |

### テスト機能

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| ドメイン層テスト | ユニットテスト | ✅ 完了（100%） |
| リポジトリ層テスト | 統合テスト | ✅ 完了（20+ケース） |
| E2Eテスト | Playwright | ✅ 完了（6スイート） |
| パフォーマンステスト | k6負荷テスト | ✅ 完了（3シナリオ） |
| CI/CD自動テスト | GitHub Actions | ✅ 完了 |

### 開発支援機能

| 機能 | 説明 | 実装状況 |
|------|------|---------|
| シードデータ | 開発用サンプルデータ | ✅ 完了 |
| 大量データ生成 | パフォーマンステスト用 | ✅ 完了 |
| Makefile統合 | 40+開発コマンド | ✅ 完了 |
| Docker Compose | ローカル環境構築 | ✅ 完了 |
| Hot Reload | 開発時の自動リロード | ✅ 完了 |

## プロジェクトの特徴

### 🎯 ビジネス価値
- **早期課題発見**: 遅延プロジェクトを即座に可視化
- **データ駆動意思決定**: 組織横断の統計データに基づく判断
- **透明性向上**: 全社プロジェクトの状況を一元管理
- **工数削減**: Jira自動同期により手動集計を不要に

### 🏗️ 技術的特徴
- **Clean Architecture**: 保守性・テスト容易性の高い設計
- **型安全性**: Go + TypeScriptによる型安全な開発
- **高品質**: 包括的なテストカバレッジ（ドメイン層100%）
- **パフォーマンス**: 10,000+チケット環境での動作確認済み
- **セキュリティ**: JWT認証、監査ログ、パスワード暗号化
- **スケーラビリティ**: Goroutineによる並行処理、効率的なDB設計

### 📚 ドキュメント充実
- 日本語・英語両対応
- ユーザーマニュアル（60+ページ）
- クイックスタートガイド（5分）
- API仕様書
- データベーススキーマ設計書
- E2Eテストガイド
- パフォーマンステストガイド
- シードデータガイド

### 🚀 開発者体験
- Makefileによる簡単コマンド実行
- Docker Composeで即座に起動
- シードデータで開発環境すぐ構築
- Hot Reloadで快適な開発
- CI/CDで品質保証

## 今後の拡張案

### 追加機能（実装可能）
- [ ] ダッシュボードのカスタマイズ機能
- [ ] エクスポート機能（Excel/CSV）
- [ ] 通知機能（Slack/Email連携）
- [ ] レポート自動生成
- [ ] カスタムフィルタ保存
- [ ] プロジェクトテンプレート
- [ ] 複数Jiraインスタンス対応

### インフラ拡張（AWS環境）
- [ ] Terraform による IaC
- [ ] AWS Fargate へのデプロイ
- [ ] Aurora PostgreSQL
- [ ] CloudFront + S3（フロントエンド）
- [ ] CloudWatch 監視
- [ ] Auto Scaling 設定

## よくある質問（FAQ）

### Q: どのような組織に適していますか？
A: 複数のJiraプロジェクトを管理している組織、特に以下のような課題を持つ組織に最適です：
- プロジェクトの遅延状況を一元管理したい
- 組織横断でのプロジェクト進捗を可視化したい
- 経営層・PMOが全社状況を把握したい

### Q: Jira Cloudのみ対応ですか？
A: 現在はJira Cloud REST API v3に対応しています。Jira Server/Data Centerへの対応は、APIエンドポイントの調整で可能です。

### Q: どれくらいのデータ量に対応できますか？
A: パフォーマンステストで10,000+チケット環境での動作を確認済みです。データベースのインデックス設計により、さらに大規模なデータにも対応可能です。

### Q: セキュリティは十分ですか？
A: 以下のセキュリティ対策を実装しています：
- JWT認証・認可
- ロールベースアクセス制御
- パスワード暗号化（bcrypt）
- 監査ログ（全APIリクエスト記録）
- HTTPS通信（本番環境）

### Q: 導入にどれくらい時間がかかりますか？
A: ローカル環境での動作確認は5-10分程度です。本番環境へのデプロイは、AWSインフラの準備状況により異なりますが、Terraformを使用すれば半日〜1日程度で可能です。

### Q: カスタマイズは可能ですか？
A: Clean Architectureに基づいた設計のため、拡張・カスタマイズが容易です。以下のようなカスタマイズが可能です：
- 遅延判定ロジックのカスタマイズ
- 新しいダッシュボードの追加
- 独自の通知機能の実装
- 外部システム連携

