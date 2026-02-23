# トラブルシューティングガイド

## よくある問題と解決方法

---

### バックエンドが起動しない

**症状**: `make up` 後にバックエンドが`Exited`状態になる。

**確認方法**:
```bash
docker compose ps
docker compose logs backend
```

**原因1**: データベースへの接続失敗

ログに `failed to connect to database` が出る場合:
```bash
# PostgreSQLの起動状態を確認
docker compose ps postgres

# PostgreSQLが起動していない場合
make up  # 再起動

# PostgreSQLのヘルスチェックを確認
docker inspect project-viz-db | grep -A5 Health
```

**原因2**: 環境変数ファイルが不足

```bash
ls backend/.env
# 存在しない場合
cp backend/.env.example backend/.env
```

---

### フロントエンドでAPIエラーが表示される

**症状**: ページが `データの取得に失敗しました` などのエラーを表示する。

**確認方法**:
```bash
# バックエンドが起動しているか確認
curl http://localhost:8080/health

# ブラウザのDevToolsで Network タブを確認
# /api/v1/* へのリクエストのステータスを確認
```

**原因1**: バックエンドが停止している → 「バックエンドが起動しない」を参照

**原因2**: CORSエラー

バックエンドはデフォルトで全オリジンを許可しています (`Access-Control-Allow-Origin: *`)。
CORSエラーが出る場合はバックエンドが正しく起動しているか確認してください。

**原因3**: フロントエンドのプロキシ設定

開発時は`vite.config.ts`のプロキシが `/api` → `http://localhost:8080` に転送します。
バックエンドがポート8080で起動していることを確認してください。

---

### 組織管理ページでエラーになる (`/organizations/manage`)

**症状**: 組織管理ページでエラーが表示される。

**確認方法**:
```bash
# 未割り当てプロジェクト取得APIを直接確認
curl http://localhost:8080/api/v1/projects?unassigned=true
```

**期待されるレスポンス** (空の場合も `[]` を返すべき):
```json
{"data":[],"pagination":{"page":1,"per_page":20,"total":0,"total_pages":1}}
```

**問題のあるレスポンス** (`null` はフロントエンドでエラーになる):
```json
{"data":null,...}
```

**解決方法**: バックエンドのハンドラーで `make([]T, 0)` を使って初期化されているか確認する。
この問題は既に修正済みですが、古いDockerイメージを使用している場合は再ビルドが必要です:
```bash
docker compose build backend
docker compose up -d backend
```

---

### データベースのマイグレーションが失敗する

**症状**: `make db-migrate` でエラーが発生する。

**確認方法**:
```bash
make db-version
# "error: Dirty database version X. Fix and force version." のような出力が出る
```

**解決方法**:

1. エラーの原因となったマイグレーションを確認

```bash
cat database/migrations/000001_initial_schema.up.sql
```

2. データベースを手動で修正するか、前のバージョンに強制設定

```bash
# 前のバージョン (例: バージョン0) に強制設定
make db-force V=0

# 再度マイグレーション実行
make db-migrate
```

3. 開発環境で完全リセットしたい場合

```bash
make db-reset  # 全ロールバック → 全マイグレーション再適用
```

---

### Go ビルドエラー

**症状**: `make backend-build` や Dockerビルドでコンパイルエラーが発生する。

**確認方法**:
```bash
cd backend && go build ./...
```

**原因1**: `go.sum` が最新でない

```bash
cd backend && go mod tidy
```

**原因2**: 依存パッケージのバージョン不整合

```bash
cd backend
go mod download
go mod verify
```

---

### フロントエンドのビルドエラー

**症状**: `npm run build` や `npm run dev` でエラーが発生する。

**確認方法**:
```bash
cd frontend
npm install
npm run type-check
npm run lint
```

**原因1**: `node_modules` が壊れている

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

**原因2**: TypeScript型エラー

```bash
cd frontend && npm run type-check
```

エラーメッセージに従って型定義を修正します。

---

### Playwright E2Eテストが失敗する

**症状**: `npm run test:e2e` で一部テストがタイムアウトまたは失敗する。

**確認方法**:
```bash
cd frontend
npm run test:e2e -- --reporter=list
```

**原因1**: フロントエンド開発サーバーが起動していない

Playwright設定 (`playwright.config.ts`) で `reuseExistingServer: true` の場合は、
事前に `npm run dev` でサーバーを起動してください:
```bash
npm run dev &
npm run test:e2e
```

**原因2**: セレクターが実際のUIと一致しない

デバッグモードで実行して確認:
```bash
npm run test:e2e:ui
```

---

### k6パフォーマンステストが実行できない

**症状**: `./performance/run.sh smoke` でエラーが発生する。

**確認方法**:
```bash
k6 version
```

**原因1**: k6がインストールされていない

```bash
brew install k6   # macOS
```

**原因2**: バックエンドが起動していない

```bash
curl http://localhost:8080/health
# エラーの場合
make up && make db-migrate
```

**原因3**: BASE_URLが間違っている

```bash
# 対象サーバーを明示的に指定
BASE_URL=http://localhost:8080 ./performance/run.sh smoke
```

---

### PostgreSQLへの接続が拒否される

**症状**: `psql: error: connection to server on socket "/tmp/.s.PGSQL.5432" failed`

**確認方法**:
```bash
docker compose ps postgres
```

**解決方法**:
```bash
# PostgreSQLコンテナが停止していたら起動
make up

# ポートが使用されているか確認
lsof -i :5432

# 別プロセスがポート5432を使用している場合
# docker-compose.ymlのポートマッピングを変更するか、他のプロセスを停止
```

---

## ログレベルの変更

デバッグが必要な場合は`.env`でログレベルを変更します:

```
LOG_LEVEL=debug   # debug, info, warn, error
```

変更後はバックエンドを再起動:
```bash
docker compose restart backend
```

## データベースの状態確認クエリ

```sql
-- テーブルのレコード数確認
SELECT 'organizations' AS tbl, COUNT(*) FROM organizations
UNION ALL SELECT 'projects', COUNT(*) FROM projects
UNION ALL SELECT 'issues', COUNT(*) FROM issues
UNION ALL SELECT 'sync_logs', COUNT(*) FROM sync_logs;

-- 遅延ステータスの分布確認
SELECT delay_status, COUNT(*) FROM issues GROUP BY delay_status;

-- 組織に割り当てられていないプロジェクト
SELECT id, key, name FROM projects WHERE organization_id IS NULL;

-- 各組織のプロジェクト数
SELECT o.name, COUNT(p.id) AS project_count
FROM organizations o
LEFT JOIN projects p ON o.id = p.organization_id
GROUP BY o.id, o.name
ORDER BY o.path;
```

## コンテナのリセット (データ消去)

**注意**: データがすべて削除されます。開発環境でのみ実施してください。

```bash
docker compose down -v   # コンテナとボリューム削除
make up                  # 再起動
make db-migrate          # マイグレーション再適用
```
