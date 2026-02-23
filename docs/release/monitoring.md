# 監視・アラート設定ガイド

本番環境で設定すべき監視項目とアラートの設定方法をまとめます。

## 監視アーキテクチャ

```
┌─────────────────────────────────────────────────────────────────┐
│                        AWS CloudWatch                           │
│                                                                 │
│  ECS Fargate ─── メトリクス ──┐                                 │
│  ALB         ─── メトリクス ──┼──→ Dashboard + Alarms ──→ SNS  │
│  Aurora      ─── メトリクス ──┘                ↓               │
│  Application ─── ログ ──────→ Log Groups    Slack / Email       │
└─────────────────────────────────────────────────────────────────┘
```

## CloudWatch アラーム設定

### ECS (バックエンド)

| アラーム名 | メトリクス | 閾値 | 評価期間 | 重要度 |
|---|---|---|---|---|
| ecs-cpu-high | ECS/CPUUtilization | > 80% | 5分 × 2 | WARNING |
| ecs-cpu-critical | ECS/CPUUtilization | > 95% | 5分 × 1 | CRITICAL |
| ecs-memory-high | ECS/MemoryUtilization | > 80% | 5分 × 2 | WARNING |
| ecs-task-stopped | ECS RunningTaskCount | < 1 | 1分 × 1 | CRITICAL |

### ALB (ロードバランサー)

| アラーム名 | メトリクス | 閾値 | 評価期間 | 重要度 |
|---|---|---|---|---|
| alb-5xx-rate | ApplicationELB/HTTPCode_Target_5XX_Count | > 10件/分 | 1分 × 2 | CRITICAL |
| alb-4xx-rate | ApplicationELB/HTTPCode_Target_4XX_Count | > 100件/分 | 5分 × 2 | WARNING |
| alb-latency | ApplicationELB/TargetResponseTime | > 1000ms | 5分 × 2 | WARNING |
| alb-unhealthy-host | ApplicationELB/UnHealthyHostCount | >= 1 | 1分 × 1 | CRITICAL |

### Aurora PostgreSQL

| アラーム名 | メトリクス | 閾値 | 評価期間 | 重要度 |
|---|---|---|---|---|
| aurora-cpu-high | RDS/CPUUtilization | > 80% | 5分 × 2 | WARNING |
| aurora-connections-high | RDS/DatabaseConnections | > 80 | 5分 × 2 | WARNING |
| aurora-free-storage | RDS/FreeStorageSpace | < 5GB | 5分 × 1 | CRITICAL |
| aurora-replica-lag | RDS/AuroraReplicaLag | > 1000ms | 5分 × 2 | WARNING |

## CloudWatch アラーム作成（AWS CLI）

### ECSタスク数アラーム

```bash
aws cloudwatch put-metric-alarm \
  --alarm-name "ecs-task-stopped" \
  --alarm-description "ECS running tasks dropped below 1" \
  --metric-name RunningTaskCount \
  --namespace ECS/ContainerInsights \
  --dimensions Name=ClusterName,Value=project-viz-prod \
                Name=ServiceName,Value=project-viz-api-production \
  --statistic Minimum \
  --period 60 \
  --evaluation-periods 1 \
  --threshold 1 \
  --comparison-operator LessThanThreshold \
  --alarm-actions arn:aws:sns:ap-northeast-1:<account-id>:project-viz-alerts \
  --treat-missing-data breaching
```

### ALB 5xxエラーアラーム

```bash
aws cloudwatch put-metric-alarm \
  --alarm-name "alb-5xx-rate" \
  --alarm-description "ALB 5xx errors exceeded threshold" \
  --metric-name HTTPCode_Target_5XX_Count \
  --namespace AWS/ApplicationELB \
  --dimensions Name=LoadBalancer,Value=<alb-arn-suffix> \
  --statistic Sum \
  --period 60 \
  --evaluation-periods 2 \
  --threshold 10 \
  --comparison-operator GreaterThanThreshold \
  --alarm-actions arn:aws:sns:ap-northeast-1:<account-id>:project-viz-alerts
```

## ログ管理

### アプリケーションログ（CloudWatch Logs）

バックエンドはJSON形式でログを出力します（Uber Zap）。ECSのログドライバーで自動的にCloudWatch Logsに送信されます。

**ロググループ名**: `/ecs/project-viz-api-production`

#### メトリクスフィルターの設定

5xxエラーをカウントするメトリクスフィルターを作成します:

```bash
aws logs put-metric-filter \
  --log-group-name "/ecs/project-viz-api-production" \
  --filter-name "5xx-errors" \
  --filter-pattern '{ $.status >= 500 }' \
  --metric-transformations \
    metricName=5xxErrorCount,metricNamespace=ProjectViz/API,metricValue=1,defaultValue=0
```

レスポンスタイムの平均値を計測するメトリクスフィルター（将来実装用）:

```bash
# バックエンドのログにdurationフィールドを追加した場合
aws logs put-metric-filter \
  --log-group-name "/ecs/project-viz-api-production" \
  --filter-name "response-time" \
  --filter-pattern '{ $.msg = "Request processed" }' \
  --metric-transformations \
    metricName=ResponseTime,metricNamespace=ProjectViz/API,metricValue='$.duration',unit=Milliseconds
```

### ログの保持期間設定

```bash
# 90日間保持
aws logs put-retention-policy \
  --log-group-name "/ecs/project-viz-api-production" \
  --retention-in-days 90
```

## SNS通知の設定

### トピックの作成

```bash
aws sns create-topic --name project-viz-alerts
```

### Slackへの通知（AWS Chatbot）

1. AWS Chatbot コンソールでSlackワークスペースを連携
2. SNSトピック `project-viz-alerts` をChatbotに設定
3. 通知先のSlackチャンネルを設定（例: `#project-viz-alerts`）

### メール通知

```bash
aws sns subscribe \
  --topic-arn arn:aws:sns:ap-northeast-1:<account-id>:project-viz-alerts \
  --protocol email \
  --notification-endpoint alerts@your-company.com
```

## CloudWatch ダッシュボード

以下のウィジェットを含むダッシュボードを作成します。

```json
{
  "widgets": [
    {
      "type": "metric",
      "properties": {
        "title": "ECS CPU / Memory",
        "metrics": [
          ["ECS/ContainerInsights", "CPUUtilization", "ClusterName", "project-viz-prod"],
          ["ECS/ContainerInsights", "MemoryUtilization", "ClusterName", "project-viz-prod"]
        ],
        "period": 60,
        "stat": "Average"
      }
    },
    {
      "type": "metric",
      "properties": {
        "title": "ALB Request Count / 5xx Errors",
        "metrics": [
          ["AWS/ApplicationELB", "RequestCount", "LoadBalancer", "<alb-arn-suffix>"],
          ["AWS/ApplicationELB", "HTTPCode_Target_5XX_Count", "LoadBalancer", "<alb-arn-suffix>"]
        ],
        "period": 60,
        "stat": "Sum"
      }
    },
    {
      "type": "metric",
      "properties": {
        "title": "ALB Response Time (p95)",
        "metrics": [
          ["AWS/ApplicationELB", "TargetResponseTime", "LoadBalancer", "<alb-arn-suffix>"]
        ],
        "period": 60,
        "stat": "p95"
      }
    },
    {
      "type": "metric",
      "properties": {
        "title": "Aurora CPU / Connections",
        "metrics": [
          ["AWS/RDS", "CPUUtilization", "DBClusterIdentifier", "project-viz-prod"],
          ["AWS/RDS", "DatabaseConnections", "DBClusterIdentifier", "project-viz-prod"]
        ],
        "period": 60,
        "stat": "Average"
      }
    }
  ]
}
```

## ヘルスチェックエンドポイント

ALBのターゲットグループのヘルスチェックには `/health` を使用します。

| 設定項目 | 値 |
|---|---|
| プロトコル | HTTP |
| パス | `/health` |
| 正常しきい値 | 2 |
| 異常しきい値 | 3 |
| タイムアウト | 5秒 |
| 間隔 | 30秒 |
| 成功コード | 200 |

## バックアップ確認手順

```bash
# Aurora の自動バックアップ一覧確認
aws rds describe-db-cluster-snapshots \
  --db-cluster-identifier project-viz-prod \
  --snapshot-type automated \
  --query 'DBClusterSnapshots[*].{ID:DBClusterSnapshotIdentifier,Time:SnapshotCreateTime,Status:Status}'

# 手動スナップショットの作成（デプロイ前に実施）
aws rds create-db-cluster-snapshot \
  --db-cluster-identifier project-viz-prod \
  --db-cluster-snapshot-identifier "project-viz-prod-pre-deploy-$(date +%Y%m%d)"
```

## 定期確認タスク

| 頻度 | タスク |
|---|---|
| 毎日 | CloudWatchアラームのステータス確認 |
| 毎週 | ログのエラー傾向レビュー |
| 毎月 | DBストレージ使用量とバックアップの確認 |
| リリース前 | 手動スナップショットの取得 |
| 四半期 | パフォーマンステスト（k6）の実施 |
