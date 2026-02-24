# Jira連携セットアップガイド

このガイドでは、JiraからプロジェクトとチケットをDBに同期する手順を説明します。

---

## 全体の流れ

```
① Jira APIトークンを発行
    ↓
② 環境変数を設定 (.env)
    ↓
③ バッチを実行してDBに同期
    ↓
④ ブラウザで確認
```

---

## ① Jira APIトークンを発行する

バッチ処理はJiraにAPIトークン（Basic認証）でアクセスします。

1. [https://id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens) を開く
2. **Create API token** をクリック
3. ラベルに任意の名前を入力（例: `project-viz-batch`）して **Create** をクリック
4. 表示されたトークン文字列をコピーして保管する（再表示できない）

> **注意**: トークンはJiraにログインしているアカウントに紐づきます。専用のサービスアカウントを用意することを推奨します。

### 必要なJira権限

バッチで使用するアカウントには以下の権限が必要です:

| 権限 | 説明 |
|------|------|
| Browse Projects | プロジェクト一覧を取得するため |
| Browse Boards | チケット一覧を取得するため（JQL検索）|

Jira管理者に「すべてのプロジェクトの閲覧権限」を付与してもらってください。

---

## ② 環境変数を設定する

### ローカル開発環境

```bash
cd backend
cp .env.example .env
```

`.env` を開いて Jira の値を入力します:

```env
JIRA_BASE_URL=https://your-org.atlassian.net   # ← Atlassianのサブドメイン部分を変える
JIRA_EMAIL=your-email@example.com              # ← トークンを発行したアカウントのメール
JIRA_API_TOKEN=ATATT3xFfGF0...                 # ← ①でコピーしたトークン
```

> **確認方法**: Jiraの画面URLが `https://acme.atlassian.net` なら `JIRA_BASE_URL=https://acme.atlassian.net` です。末尾のスラッシュは不要です。

### 本番環境 (ECS)

本番ではAWS Secrets Managerを使います。詳細は [docs/secrets-management.md](./secrets-management.md) を参照してください。

---

## ③ バッチを実行してDBに同期する

### 前提条件

バッチを実行する前に、DBとバックエンドが起動していることを確認します:

```bash
make up          # PostgreSQL + バックエンドを起動
make db-migrate  # マイグレーション適用（初回のみ）
```

### バッチの実行

```bash
# バッチバイナリをビルド
cd backend
go build -o bin/batch ./cmd/batch/

# フル同期（全プロジェクト・全チケットを取得）
./bin/batch

# または差分同期（前回以降に更新されたチケットのみ）
BATCH_SYNC_MODE=delta ./bin/batch
```

### Dockerで実行する場合

```bash
docker compose run --rm backend \
  sh -c "BATCH_SYNC_MODE=full /app/batch"
```

---

## ④ 同期結果を確認する

### ログで確認

バッチ実行中に以下のようなログが表示されます:

```json
{"level":"info","msg":"full sync started"}
{"level":"info","msg":"fetched projects","count":12}
{"level":"info","msg":"fetched issues","count":284}
{"level":"info","msg":"full sync finished","status":"SUCCESS","projects_synced":12,"issues_synced":284,"duration":"8.3s"}
```

### DBで確認

```bash
make db-connect
```

```sql
-- 同期されたプロジェクト数を確認
SELECT COUNT(*) FROM projects;

-- 同期されたチケット数を確認
SELECT COUNT(*) FROM issues;

-- 最新の同期ログを確認
SELECT sync_type, status, projects_synced, issues_synced,
       duration_seconds, executed_at
FROM sync_logs
ORDER BY executed_at DESC
LIMIT 5;
```

### ブラウザで確認

フロントエンドを起動してダッシュボードを開きます:

```bash
cd frontend && npm run dev
```

ブラウザで http://localhost:3000 を開き、プロジェクトとチケットが表示されることを確認します。

---

## 定期実行の設定（本番）

本番環境では CloudWatch EventBridge でバッチを定期実行します。

| 種別 | スケジュール | 用途 |
|------|------------|------|
| Full Sync | 毎日 01:00 JST | 全件を最新状態に同期 |
| Delta Sync | 1時間ごと | 直近の変更をリアルタイムに反映 |

詳細は [docs/batch-schedule.md](./batch-schedule.md) を参照してください。

---

## トラブルシューティング

### `missing required environment variables: JIRA_BASE_URL` と表示される

`backend/.env` ファイルが存在しないか、変数が未設定です。

```bash
cat backend/.env | grep JIRA_  # 設定値を確認
```

### `jira API error: HTTP 401` が出る

認証情報が間違っています。以下を確認してください:

- `JIRA_EMAIL` がAPIトークンを発行したアカウントのメールアドレスであること
- `JIRA_API_TOKEN` が正しくコピーされていること（スペースや改行が混入していないか）
- Atlassianアカウントのパスワードではなく **APIトークン** を使っていること

### `jira API error: HTTP 403` が出る

アカウントの権限不足です。Jira管理者に「全プロジェクトの閲覧権限」を付与してもらってください。

### `jira API error: HTTP 429` が出る

JiraのAPIレート制限に達しました。バッチは自動でリトライします。プロジェクト数が多い場合は `BATCH_WORKER_COUNT` を減らしてください:

```env
BATCH_WORKER_COUNT=2
```

### プロジェクトは同期されたがチケットが0件

JQLの検索結果が空の可能性があります。Jira画面で以下のJQLを試してください:

```
project = YOUR_PROJECT_KEY ORDER BY updated ASC
```

チケットが表示されない場合、プロジェクトのチケットがアカウントから閲覧できていません。
