# Backend API (Go)

全社プロジェクト進捗可視化プラットフォームのバックエンドAPI。

## 技術スタック

- **言語**: Go 1.21+
- **Webフレームワーク**: Gin
- **データベース**: PostgreSQL 15+ (sqlx)
- **ログ**: Uber Zap
- **環境変数**: godotenv
- **マイグレーション**: golang-migrate

## プロジェクト構造

```
backend/
├── cmd/
│   └── api/                    # アプリケーションエントリーポイント
│       └── main.go
├── internal/
│   ├── domain/                 # ドメイン層（エンティティ・ビジネスルール）
│   │   ├── organization.go
│   │   ├── project.go
│   │   ├── issue.go
│   │   └── sync_log.go
│   ├── usecase/                # ユースケース層（アプリケーションロジック）
│   │   ├── organization_usecase.go
│   │   ├── project_usecase.go
│   │   ├── issue_usecase.go
│   │   └── dashboard_usecase.go
│   ├── interface/              # インターフェース層
│   │   ├── handler/            # HTTPハンドラー
│   │   │   ├── organization_handler.go
│   │   │   ├── project_handler.go
│   │   │   ├── issue_handler.go
│   │   │   └── dashboard_handler.go
│   │   └── repository/         # リポジトリインターフェース
│   │       ├── organization_repository.go
│   │       ├── project_repository.go
│   │       └── issue_repository.go
│   └── infrastructure/         # インフラストラクチャ層
│       ├── postgres/           # PostgreSQL実装
│       │   ├── organization_repository_impl.go
│       │   ├── project_repository_impl.go
│       │   └── issue_repository_impl.go
│       └── router/             # ルーティング設定
│           └── router.go
├── pkg/                        # 共通パッケージ
│   ├── logger/                 # ロガー
│   │   └── logger.go
│   └── config/                 # 設定管理
│       └── config.go
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

## アーキテクチャ

Clean Architectureに基づいたレイヤー分離:

1. **Domain層** (`internal/domain`)
   - エンティティ定義
   - ビジネスルール
   - フレームワーク非依存

2. **Usecase層** (`internal/usecase`)
   - アプリケーションロジック
   - ドメインオブジェクトの操作
   - リポジトリインターフェースへの依存

3. **Interface層** (`internal/interface`)
   - HTTPハンドラー（Ginを使用）
   - リポジトリインターフェース定義
   - リクエスト/レスポンスのDTO

4. **Infrastructure層** (`internal/infrastructure`)
   - データベースアクセス実装
   - ルーティング設定
   - 外部サービス連携

## セットアップ

### 前提条件

- Go 1.21以上
- PostgreSQL 15以上
- Docker（推奨）

### ローカル開発環境

1. 依存関係のインストール:
```bash
cd backend
go mod download
```

2. 環境変数の設定:
```bash
cp .env.example .env
# .envファイルを編集
```

3. データベースの起動:
```bash
make db-up
make db-migrate
```

4. アプリケーションの起動:
```bash
make run
```

サーバーが `http://localhost:8080` で起動します。

## Makeコマンド

```bash
# アプリケーション実行
make run

# ビルド
make build

# テスト実行
make test

# テストカバレッジ
make coverage

# リンター実行
make lint

# フォーマット
make fmt

# 依存関係の更新
make deps

# すべてのチェック（lint + test）
make check
```

## 環境変数

`.env` ファイルで以下の環境変数を設定:

```
# サーバー設定
PORT=8080
GIN_MODE=debug

# データベース設定
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=project_visualization
DB_SSLMODE=disable

# ログ設定
LOG_LEVEL=debug
LOG_FORMAT=json
```

## API エンドポイント

### 組織管理

- `GET /api/v1/organizations` - 組織一覧取得
- `GET /api/v1/organizations/:id` - 組織詳細取得
- `POST /api/v1/organizations` - 組織作成
- `PUT /api/v1/organizations/:id` - 組織更新
- `DELETE /api/v1/organizations/:id` - 組織削除
- `GET /api/v1/organizations/:id/children` - 子組織一覧取得

### プロジェクト管理

- `GET /api/v1/projects` - プロジェクト一覧取得
- `GET /api/v1/projects/:id` - プロジェクト詳細取得
- `PUT /api/v1/projects/:id` - プロジェクト更新
- `PUT /api/v1/projects/:id/organization` - プロジェクトの組織紐付け

### チケット管理

- `GET /api/v1/issues` - チケット一覧取得（フィルタ対応）
- `GET /api/v1/issues/:id` - チケット詳細取得
- `GET /api/v1/projects/:id/issues` - プロジェクトのチケット一覧

### ダッシュボード

- `GET /api/v1/dashboard/summary` - 全社サマリ取得
- `GET /api/v1/dashboard/organizations/:id` - 組織別サマリ取得
- `GET /api/v1/dashboard/projects/:id` - プロジェクト別サマリ取得

### ヘルスチェック

- `GET /health` - ヘルスチェック
- `GET /ready` - Readinessチェック

## テスト

```bash
# すべてのテストを実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# カバレッジレポート生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Docker

### ビルド

```bash
docker build -t project-viz-backend .
```

### 実行

```bash
docker run -p 8080:8080 --env-file .env project-viz-backend
```

### Docker Compose

```bash
# すべてのサービスを起動（PostgreSQL + Backend）
docker-compose up -d

# ログ確認
docker-compose logs -f backend

# 停止
docker-compose down
```

## コーディング規約

### 命名規則

- パッケージ名: 小文字、単数形
- インターフェース名: `~er` サフィックス（例: `Organizer`）
- 構造体: PascalCase
- メソッド: PascalCase（公開）、camelCase（非公開）
- 変数: camelCase

### エラーハンドリング

- エラーは常に返却する
- エラーメッセージは小文字で開始
- カスタムエラー型を活用

### ログ

- 構造化ログを使用（Zap）
- ログレベルを適切に設定
  - Debug: デバッグ情報
  - Info: 一般的な情報
  - Warn: 警告
  - Error: エラー

## デプロイ

### AWS Fargate / ECS

Terraformで構築したインフラにデプロイ:

```bash
# ECRにプッシュ
make docker-push

# ECSサービス更新
make deploy
```

## トラブルシューティング

### データベース接続エラー

1. PostgreSQLが起動しているか確認
2. 環境変数が正しく設定されているか確認
3. ファイアウォール設定を確認

### マイグレーションエラー

```bash
# マイグレーション状態を確認
make db-version

# dirty状態からの回復
make db-force V=1
```

## 関連ドキュメント

- [SPEC.md](../SPEC.md) - プロジェクト要件定義書
- [database/README.md](../database/README.md) - データベース管理
- [API Documentation](./docs/api.md) - API仕様書（準備中）
