# バッチ監視・アラート設定

このドキュメントでは、バッチ処理の CloudWatch メトリクス監視・アラート・ダッシュボードの設定方法を説明します。

## メトリクス仕様

バッチは **CloudWatch Embedded Metrics Format (EMF)** を使用してメトリクスを送信します。
ECS タスクに CloudWatch Logs エージェントが設定されていれば、標準出力に書き出すだけで
自動的にカスタムメトリクスとして収集されます。

### 送信メトリクス一覧

| メトリクス名 | 単位 | 説明 |
|------------|------|------|
| `SyncSuccess` | Count | 成功=1、失敗=0 |
| `DurationSeconds` | Seconds | バッチ実行時間（秒） |
| `IssuesSynced` | Count | upsertしたチケット数 |
| `ProjectsSynced` | Count | upsertしたプロジェクト数（Full Syncのみ有効） |

**Namespace**: 環境変数 `METRICS_NAMESPACE` で指定（推奨: `SandboxProjectManagement/Batch`）
**Dimension**: `SyncType` = `FULL` または `DELTA`

### 有効化方法

ECS タスク定義の環境変数に以下を追加:

```
METRICS_NAMESPACE=SandboxProjectManagement/Batch
```

## CloudWatch Alarms 設定

### アラーム一覧

#### 1. バッチ失敗アラーム（Full Sync）

```json
{
  "AlarmName": "BatchFullSync-Failure",
  "Namespace": "SandboxProjectManagement/Batch",
  "MetricName": "SyncSuccess",
  "Dimensions": [{"Name": "SyncType", "Value": "FULL"}],
  "Statistic": "Minimum",
  "Period": 86400,
  "EvaluationPeriods": 1,
  "Threshold": 1,
  "ComparisonOperator": "LessThanThreshold",
  "TreatMissingData": "breaching",
  "AlarmActions": ["<SNS_TOPIC_ARN>"]
}
```

#### 2. バッチ失敗アラーム（Delta Sync）

```json
{
  "AlarmName": "BatchDeltaSync-Failure",
  "Namespace": "SandboxProjectManagement/Batch",
  "MetricName": "SyncSuccess",
  "Dimensions": [{"Name": "SyncType", "Value": "DELTA"}],
  "Statistic": "Minimum",
  "Period": 3600,
  "EvaluationPeriods": 2,
  "Threshold": 1,
  "ComparisonOperator": "LessThanThreshold",
  "TreatMissingData": "breaching",
  "AlarmActions": ["<SNS_TOPIC_ARN>"]
}
```

#### 3. Full Sync 実行時間超過アラーム

```json
{
  "AlarmName": "BatchFullSync-DurationHigh",
  "Namespace": "SandboxProjectManagement/Batch",
  "MetricName": "DurationSeconds",
  "Dimensions": [{"Name": "SyncType", "Value": "FULL"}],
  "Statistic": "Maximum",
  "Period": 86400,
  "EvaluationPeriods": 1,
  "Threshold": 300,
  "ComparisonOperator": "GreaterThanThreshold",
  "AlarmActions": ["<SNS_TOPIC_ARN>"]
}
```

### AWS CLI でのアラーム作成例

```bash
aws cloudwatch put-metric-alarm \
  --alarm-name "BatchFullSync-Failure" \
  --namespace "SandboxProjectManagement/Batch" \
  --metric-name "SyncSuccess" \
  --dimensions Name=SyncType,Value=FULL \
  --statistic Minimum \
  --period 86400 \
  --evaluation-periods 1 \
  --threshold 1 \
  --comparison-operator LessThanThreshold \
  --treat-missing-data breaching \
  --alarm-actions "<SNS_TOPIC_ARN>"
```

## SNS トピック設定

### トピック作成

```bash
aws sns create-topic --name batch-alerts
```

### メール通知の登録

```bash
aws sns subscribe \
  --topic-arn "<SNS_TOPIC_ARN>" \
  --protocol email \
  --notification-endpoint "ops-team@example.com"
```

### Slack 通知（Lambda 経由の場合）

Lambda (Node.js/Python) を SNS サブスクライバーとして登録し、
CloudWatch Alarm → SNS → Lambda → Slack Incoming Webhook の経路で通知を送る。

## CloudWatch Dashboard 設定

### ダッシュボード作成（AWS CLI）

```bash
aws cloudwatch put-dashboard \
  --dashboard-name "BatchSyncMonitoring" \
  --dashboard-body file://dashboard.json
```

`dashboard.json` の構成例:

```json
{
  "widgets": [
    {
      "type": "metric",
      "properties": {
        "title": "Sync Success Rate",
        "metrics": [
          ["SandboxProjectManagement/Batch", "SyncSuccess", "SyncType", "FULL"],
          [".", ".", ".", "DELTA"]
        ],
        "stat": "Minimum",
        "period": 3600,
        "view": "timeSeries"
      }
    },
    {
      "type": "metric",
      "properties": {
        "title": "Issues Synced",
        "metrics": [
          ["SandboxProjectManagement/Batch", "IssuesSynced", "SyncType", "FULL"],
          [".", ".", ".", "DELTA"]
        ],
        "stat": "Sum",
        "period": 3600,
        "view": "timeSeries"
      }
    },
    {
      "type": "metric",
      "properties": {
        "title": "Sync Duration (seconds)",
        "metrics": [
          ["SandboxProjectManagement/Batch", "DurationSeconds", "SyncType", "FULL"],
          [".", ".", ".", "DELTA"]
        ],
        "stat": "Maximum",
        "period": 3600,
        "view": "timeSeries"
      }
    }
  ]
}
```

## ローカルでの動作確認

EMF 出力をローカルで確認する場合は、`METRICS_NAMESPACE` を設定して実行します:

```bash
METRICS_NAMESPACE=TestBatch/Local \
BATCH_SYNC_MODE=full \
JIRA_BASE_URL=... \
go run ./cmd/batch/
```

正常に動作した場合、標準出力に以下のような JSON が出力されます:

```json
{
  "_aws": {
    "Timestamp": 1234567890000,
    "CloudWatchMetrics": [{
      "Namespace": "TestBatch/Local",
      "Dimensions": [["SyncType"]],
      "Metrics": [
        {"Name": "SyncSuccess", "Unit": "Count"},
        {"Name": "DurationSeconds", "Unit": "Seconds"},
        {"Name": "IssuesSynced", "Unit": "Count"},
        {"Name": "ProjectsSynced", "Unit": "Count"}
      ]
    }]
  },
  "SyncType": "FULL",
  "SyncSuccess": 1,
  "DurationSeconds": 12.5,
  "IssuesSynced": 150,
  "ProjectsSynced": 5
}
```
