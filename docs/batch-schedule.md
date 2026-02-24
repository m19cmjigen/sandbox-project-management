# Batch Schedule — CloudWatch EventBridge

このドキュメントでは、バッチ処理の実行スケジュールを CloudWatch EventBridge で設定する方法を説明します。

## バッチの種類

| 種別 | 環境変数 | 説明 |
|------|----------|------|
| Full Sync | `BATCH_SYNC_MODE=full`（デフォルト）| 全プロジェクト・全チケットを取得して DB を更新 |
| Delta Sync | `BATCH_SYNC_MODE=delta` | 前回成功した Delta Sync 以降に更新されたチケットのみを取得・upsert |

## EventBridge スケジュールルール設定

### Full Sync — 毎日 01:00 JST

```
cron(0 16 * * ? *)
```

> JST 01:00 = UTC 16:00 (前日)

**コンソール設定例**

| 項目 | 値 |
|------|----|
| Schedule expression | `cron(0 16 * * ? *)` |
| Target | ECS Task または Lambda |
| Environment variable | `BATCH_SYNC_MODE=full` |

### Delta Sync — 1時間ごと

```
rate(1 hour)
```

**コンソール設定例**

| 項目 | 値 |
|------|----|
| Schedule expression | `rate(1 hour)` |
| Target | ECS Task または Lambda |
| Environment variable | `BATCH_SYNC_MODE=delta` |

## 環境変数一覧

| 変数名 | 必須 | デフォルト | 説明 |
|--------|------|-----------|------|
| `JIRA_BASE_URL` | Yes | — | Jira Cloud ベース URL（例: `https://your-org.atlassian.net`）|
| `JIRA_EMAIL` | Yes | — | Jira 認証用メールアドレス |
| `JIRA_API_TOKEN` | Yes | — | Jira API トークン |
| `BATCH_SYNC_MODE` | No | `full` | 実行モード: `full` または `delta` |
| `BATCH_WORKER_COUNT` | No | `5` | Full Sync 時のプロジェクト並列フェッチ数 |

## Delta Sync のフォールバック動作

前回の Delta Sync 成功記録が `sync_logs` に存在しない場合、**現在時刻から1時間前**をフォールバックとして使用します。
これにより、初回実行時または sync_logs がリセットされた場合でも安全に動作します。

## 実行履歴の確認

```sql
SELECT sync_type, status, executed_at, completed_at,
       issues_synced, duration_seconds, error_message
FROM sync_logs
ORDER BY executed_at DESC
LIMIT 20;
```
