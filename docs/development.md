# 開発者ガイド

## 概要

本ドキュメントは、開発環境のセットアップから実装規約、テスト方法までを説明します。

## 開発環境のセットアップ

### 必要なツール

| ツール | バージョン | インストール |
|---|---|---|
| Go | 1.21以上 | `brew install go` |
| Node.js | 18以上 | `brew install node` |
| Docker Desktop | 最新 | https://docs.docker.com/desktop/mac/ |
| golang-migrate | 最新 | `brew install golang-migrate` |
| golangci-lint | 最新 | `brew install golangci-lint` |
| goimports | 最新 | `go install golang.org/x/tools/cmd/goimports@latest` |
| k6 | 最新 | `brew install k6` |

### セットアップ手順

```bash
# 1. リポジトリのクローン
git clone <repository-url>
cd sandbox-project-management

# 2. バックエンド環境変数の設定
cp backend/.env.example backend/.env

# 3. サービス起動
make up
make db-migrate

# 4. フロントエンドの依存関係インストールと起動
cd frontend && npm install && npm run dev

# 5. 動作確認
curl http://localhost:8080/health
open http://localhost:3000
```

## バックエンド開発

### ディレクトリ構成

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # エントリーポイント
├── internal/
│   ├── domain/                  # ドメイン層 (将来実装)
│   ├── usecase/                 # ユースケース層 (将来実装)
│   ├── interface/               # インターフェース層 (将来実装)
│   └── infrastructure/
│       └── router/
│           ├── router.go         # ルーター定義
│           ├── dashboard_handlers.go
│           ├── organization_handlers.go
│           ├── project_handlers.go
│           ├── issue_handlers.go
│           └── management_handlers.go
└── pkg/
    ├── config/                  # 設定
    └── logger/                  # ロガー
```

### よく使うコマンド

```bash
# API起動
make backend-run

# テスト実行
make backend-test

# カバレッジレポート生成 (backend/coverage.html)
make backend-coverage

# リント
make backend-lint

# フォーマット (goimports)
make backend-fmt

# ビルド
make backend-build
```

### コーディング規約

- コメント・ドキュメンテーション (GoDoc) は英語で書く
- エラーハンドリングを省略しない
- SQLクエリには`sqlx`を使用し、プレースホルダーで引数を渡す (SQLインジェクション防止)
- スライスは `make([]T, 0)` で初期化する (JSONシリアライズで`null`にならないように)
- ハンドラーは `handler.go` 内のクロージャとして定義する

### 新しいAPIエンドポイントの追加

1. `backend/internal/infrastructure/router/` に新しいハンドラーファイルを作成
2. `router.go` の `NewRouter()` にルートを追加
3. テストを `_test.go` ファイルに追加

例:

```go
// my_handlers.go
package router

func listMyResourceHandlerWithDB(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        // ...
        c.JSON(http.StatusOK, result)
    }
}

// router.go に追加
myResources := v1.Group("/my-resources")
{
    myResources.GET("", listMyResourceHandlerWithDB(db))
}
```

### データベースマイグレーションの追加

```bash
# 新しいマイグレーションファイルを作成
make db-migration name=add_new_table

# 生成されるファイル:
# database/migrations/000002_add_new_table.up.sql
# database/migrations/000002_add_new_table.down.sql
```

## フロントエンド開発

### ディレクトリ構成

```
frontend/src/
├── pages/              # ルートレベルコンポーネント (1ページ = 1ファイル)
├── components/         # 再利用可能UIコンポーネント
├── stores/             # Zustandストア (状態管理)
├── api/                # Axiosベースのアクセサリ
└── theme/              # Material-UIカスタムテーマ
```

### よく使うコマンド

```bash
cd frontend

npm run dev           # 開発サーバー起動 (port 3000)
npm run build         # 本番ビルド (dist/ に出力)
npm run test          # Vitestテスト実行
npm run test:ui       # Vitest UIダッシュボード
npm run test:coverage # カバレッジレポート生成
npm run test:e2e      # Playwright E2Eテスト
npm run test:e2e:ui   # Playwright UIモード
npm run lint          # ESLint実行
npm run lint:fix      # ESLint自動修正
npm run format        # Prettier実行
npm run type-check    # TypeScript型チェック
```

### コーディング規約

- JSDoc・型定義の説明コメントは英語で書く
- コンポーネント内の実装コメント (背景・理由の説明) は日本語で書く
- Vitestのテスト説明文は英語で書く
- zod-openapi等のスキーマ定義のdescriptionは英語で書く
- 絵文字を使用しない

### APIクライアントの使用

`frontend/src/api/` のアクセサリを使用します:

```typescript
import { fetchProjects } from '../api/projects'

const { data, pagination } = await fetchProjects({ page: 1, per_page: 20 })
```

### 新しいページの追加

1. `frontend/src/pages/NewPage.tsx` を作成
2. `frontend/src/App.tsx` (またはルーターファイル) にルートを追加

## テスト

### バックエンドテスト

```bash
make backend-test          # 全テスト実行
make backend-coverage      # カバレッジ付き実行 → backend/coverage.html を開く
```

テストファイルは対象ファイルと同ディレクトリに `_test.go` として配置します。

### フロントエンドユニットテスト (Vitest)

```bash
cd frontend
npm run test          # テスト実行
npm run test:coverage # カバレッジ付き
```

テストファイルは `*.test.ts(x)` の命名規則に従います。

### E2Eテスト (Playwright)

```bash
cd frontend
npm run test:e2e        # ヘッドレス実行
npm run test:e2e:ui     # UIモード (デバッグ用)
npm run test:e2e:report # HTMLレポート表示
```

テストファイルは `frontend/e2e/specs/*.spec.ts` に配置します。
APIはモック (`frontend/e2e/fixtures/api-mocks.ts`) を使用します。

### パフォーマンステスト (k6)

バックエンド起動後に実行します:

```bash
make up && make db-migrate

# テスト実行
./performance/run.sh smoke   # スモークテスト (1VU/30秒)
./performance/run.sh load    # 負荷テスト (最大10VU/90秒)
./performance/run.sh stress  # ストレステスト (最大50VU/160秒)
./performance/run.sh all     # 全テスト
```

詳細は `performance/README.md` を参照してください。

## ブランチ・コミット規約

### ブランチ名

```
feature/<ticket-id>-<brief-description>
例: feature/backend-001-organization-api
```

### コミットメッセージ

```
<type>: <summary>

例:
feat: add list organizations API
fix: return empty array instead of null for empty results
test: add E2E tests for organizations page
docs: add technical documentation
```

## コードフォーマット

```bash
# バックエンド
make backend-fmt

# フロントエンド
cd frontend && npm run format && npm run lint:fix
```

## PR作成時のチェックリスト

- [ ] `make backend-test` がパスすること
- [ ] `npm run test` がパスすること
- [ ] `npm run lint` がパスすること
- [ ] `npm run type-check` がパスすること
- [ ] 新規APIにはE2Eテストを追加すること
- [ ] 破壊的変更がある場合はAPIドキュメント (`docs/api-spec.yaml`) を更新すること
- [ ] DBスキーマ変更がある場合はマイグレーションファイルを追加すること
