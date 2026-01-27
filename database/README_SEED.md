# シードデータ (Seed Data)

開発・テスト用のサンプルデータを簡単に投入するためのツールです。

## 概要

このシードデータスクリプトは以下のデータを作成します：

### 作成されるデータ

1. **ユーザー (4名)**
   - `admin` / `admin123` - システム管理者 (Admin)
   - `pmo_manager` / `admin123` - PMOマネージャー (Manager)
   - `dept_manager` / `admin123` - 部門マネージャー (Manager)
   - `viewer` / `admin123` - 一般閲覧者 (Viewer)

2. **組織階層 (9組織)**
   ```
   本社 (Headquarters)
   ├── 営業本部 (Sales Division)
   │   ├── 東京営業所 (Tokyo Sales Office)
   │   └── 大阪営業所 (Osaka Sales Office)
   ├── 開発本部 (Development Division)
   │   ├── フロントエンドチーム (Frontend Team)
   │   ├── バックエンドチーム (Backend Team)
   │   └── インフラチーム (Infrastructure Team)
   └── PMO部 (PMO Department)
   ```

3. **プロジェクト (5プロジェクト)**
   - 新規CRMシステム導入 (SALES)
   - マーケティングオートメーション (MARKET)
   - 社内Webアプリ刷新 (WEBAPP)
   - モバイルアプリ開発 (MOBILE)
   - プロジェクト管理基盤構築 (PMOBASE)

4. **チケット (30-50チケット)**
   - 各プロジェクトに5-10個のチケット
   - 遅延ステータスの分布:
     - 🔴 RED (遅延): 期限超過
     - 🟡 YELLOW (注意): 期限間近 or 期限未設定
     - 🟢 GREEN (正常): 余裕あり or 完了

## 実行方法

### 1. データベースの準備

まず、データベースが起動していることを確認します：

```bash
# PostgreSQL起動
make db-up

# マイグレーション適用
make db-migrate
```

### 2. シードデータの投入

#### Makefileから実行（推奨）

```bash
# ルートディレクトリから
make db-seed
```

#### 直接実行

```bash
# backendディレクトリから
cd backend
go run ./scripts/seed_data.go
```

#### 環境変数の指定

デフォルトでは `localhost:5432` のデータベースに接続します。
異なる接続先の場合は `DATABASE_URL` を指定してください：

```bash
DATABASE_URL="postgres://user:password@host:port/dbname?sslmode=disable" go run ./scripts/seed_data.go
```

## 実行結果の例

```
Connected to database successfully
Starting seed data creation...
Creating users...
Created 4 users
Creating organization hierarchy...
Created 9 organizations
Creating projects...
Created 5 projects
Creating issues...
Created 47 issues
Seed data creation completed successfully!
```

## データのクリーンアップ

シードデータをクリアしたい場合：

```bash
# すべてのデータをクリア（マイグレーションのロールバック&再適用）
make db-reset

# その後、必要に応じて再度シード実行
make db-seed
```

## 注意事項

- シードデータは**開発・テスト環境専用**です
- 本番環境では実行しないでください
- 重複データを避けるため、既存データがある場合は先にクリアしてください
- ユーザー名が重複する場合、エラーメッセージが表示されますが処理は継続します

## カスタマイズ

`backend/scripts/seed_data.go` を編集することで、以下をカスタマイズできます：

- ユーザーの数と権限
- 組織構造
- プロジェクト数とプロジェクト名
- チケット数とステータス分布

## トラブルシューティング

### データベースに接続できない

```
Failed to connect to database: dial tcp 127.0.0.1:5432: connect: connection refused
```

→ PostgreSQLが起動しているか確認してください（`make db-up`）

### ユーザー作成でエラーが出る

```
Warning: failed to create user admin: pq: duplicate key value violates unique constraint
```

→ 既にシードデータが投入されています。`make db-reset` でクリアしてください

### マイグレーションエラー

```
error: Dirty database version...
```

→ `make db-force-version` でマイグレーションバージョンを修正してください

## 関連コマンド

```bash
# データベース関連
make db-up          # PostgreSQL起動
make db-down        # PostgreSQL停止
make db-migrate     # マイグレーション適用
make db-reset       # データベースリセット
make db-seed        # シードデータ投入

# 開発用
make backend-run    # バックエンド起動（シードデータ使用可能）
```
