# データベース管理

## ディレクトリ構成

```
database/
├── schema/               # スキーマ設計ドキュメント
│   ├── schema_design.md  # スキーマ設計書
│   ├── er_diagram.mmd    # ER図（Mermaid形式）
│   └── schema.sql        # 完全なDDL
├── migrations/           # マイグレーションファイル
│   ├── 000001_initial_schema.up.sql    # 初期スキーマ（UP）
│   └── 000001_initial_schema.down.sql  # 初期スキーマ（DOWN）
└── README.md            # このファイル
```

## マイグレーションツール

このプロジェクトでは [golang-migrate](https://github.com/golang-migrate/migrate) を使用します。

### インストール

#### macOS
```bash
brew install golang-migrate
```

#### Linux
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

#### Go経由
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## マイグレーション実行

### 環境変数設定

```bash
export DATABASE_URL="postgres://username:password@localhost:5432/dbname?sslmode=disable"
```

### マイグレーションの適用

```bash
# すべてのマイグレーションを適用
migrate -path database/migrations -database "${DATABASE_URL}" up

# 特定のバージョンまで適用
migrate -path database/migrations -database "${DATABASE_URL}" goto 1

# 1ステップだけ適用
migrate -path database/migrations -database "${DATABASE_URL}" up 1
```

### ロールバック

```bash
# すべてのマイグレーションをロールバック
migrate -path database/migrations -database "${DATABASE_URL}" down

# 1ステップだけロールバック
migrate -path database/migrations -database "${DATABASE_URL}" down 1
```

### マイグレーション状態の確認

```bash
migrate -path database/migrations -database "${DATABASE_URL}" version
```

### マイグレーションの強制実行（エラーリカバリ）

```bash
# dirty状態をクリアして特定バージョンに強制設定
migrate -path database/migrations -database "${DATABASE_URL}" force 1
```

## Makefileコマンド

プロジェクトルートの`Makefile`に以下のコマンドを追加予定：

```bash
# マイグレーション適用
make db-migrate

# ロールバック
make db-rollback

# マイグレーション状態確認
make db-version

# 新しいマイグレーション作成
make db-migration name=add_new_feature
```

## 開発環境セットアップ

### Dockerを使用したPostgreSQL起動

```bash
docker run --name project-viz-db \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=admin123 \
  -e POSTGRES_DB=project_visualization \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### マイグレーション適用

```bash
export DATABASE_URL="postgres://admin:admin123@localhost:5432/project_visualization?sslmode=disable"
migrate -path database/migrations -database "${DATABASE_URL}" up
```

### 接続確認

```bash
psql -h localhost -U admin -d project_visualization
```

## マイグレーションファイルの命名規則

```
{version}_{description}.{up|down}.sql
```

例:
- `000001_initial_schema.up.sql`
- `000001_initial_schema.down.sql`
- `000002_add_user_table.up.sql`
- `000002_add_user_table.down.sql`

## スキーマ設計の参照

詳細なスキーマ設計については `schema/schema_design.md` を参照してください。

## テーブル一覧

1. **organizations** - 組織階層マスタ
2. **projects** - Jiraプロジェクト情報
3. **issues** - Jiraチケット情報
4. **sync_logs** - バッチ実行ログ

## ビュー一覧

1. **project_delay_summary** - プロジェクトごとの遅延チケットサマリ
2. **organization_delay_summary** - 組織ごとの遅延プロジェクトサマリ

## CI/CD統合

GitHub Actions等のCI/CDパイプラインでマイグレーションを自動実行する場合：

```yaml
- name: Run database migrations
  env:
    DATABASE_URL: ${{ secrets.DATABASE_URL }}
  run: |
    migrate -path database/migrations -database "${DATABASE_URL}" up
```

## トラブルシューティング

### dirty状態からの回復

マイグレーション途中でエラーが発生した場合：

1. エラーの原因を確認
2. 手動でSQLを修正
3. dirty状態をクリア:
   ```bash
   migrate -path database/migrations -database "${DATABASE_URL}" force {version}
   ```

### マイグレーションファイルの修正

すでに適用済みのマイグレーションを修正する場合：

1. 新しいマイグレーションファイルを作成して修正を適用
2. 既存のマイグレーションファイルは変更しない

## 関連ドキュメント

- [SPEC.md](../SPEC.md) - プロジェクト要件定義書
- [schema/schema_design.md](schema/schema_design.md) - スキーマ設計書
