# Jira認証情報のセキュア管理

このドキュメントでは、Jira Cloud API の認証情報を安全に管理する方法を説明します。

## 認証方式の選定

**API Token 方式を採用**します（OAuth 2.0 3LO 非採用の理由:サーバー間通信のためユーザー認証フロー不要）。

| 比較項目 | API Token | OAuth 2.0 (3LO) |
|---------|-----------|-----------------|
| 用途 | サーバー間通信 | ユーザー代理アクセス |
| セットアップ | 簡単 | 複雑（リダイレクト設定等） |
| バッチ処理 | 適している | 不向き（インタラクティブ認証が必要）|
| ローテーション | Jira UIから手動 | OAuthフロー内で自動 |

## 環境変数一覧

| 変数名 | 説明 | 例 |
|--------|------|-----|
| `JIRA_BASE_URL` | Jira Cloud インスタンスのURL | `https://your-org.atlassian.net` |
| `JIRA_EMAIL` | Atlassian アカウントのメールアドレス | `batch@example.com` |
| `JIRA_API_TOKEN` | Jira API トークン | `ATATT3xFfGF0...` |

## ローカル開発

`.env` ファイルに記述します（`.gitignore` に必ず追加）:

```bash
cp backend/.env.example backend/.env
# .env を編集して Jira 認証情報を設定
```

`.env.example` の該当箇所:

```env
JIRA_BASE_URL=https://your-org.atlassian.net
JIRA_EMAIL=your-email@example.com
JIRA_API_TOKEN=your-api-token-here
```

## 本番環境: AWS Secrets Manager

### アーキテクチャ

```
[Secrets Manager] --inject env vars--> [ECS Task]
                                             |
                                        [Batch Process]
                                             |
                                       reads env vars only
                                       (no AWS SDK in app code)
```

アプリケーションコードは通常の環境変数を読むだけです。
ECS がタスク起動時に Secrets Manager から値を取得して環境変数に注入します。

### Secrets Manager へのシークレット登録

```bash
# Jira 認証情報を JSON 形式で登録
aws secretsmanager create-secret \
  --name "sandbox-pm/production/jira-credentials" \
  --description "Jira Cloud API credentials for batch sync" \
  --secret-string '{
    "base_url": "https://your-org.atlassian.net",
    "email": "batch@example.com",
    "api_token": "ATATT3xFfGF0..."
  }'
```

### ECS タスク定義への secrets 設定

```json
{
  "containerDefinitions": [{
    "name": "batch",
    "image": "YOUR_ECR_IMAGE",
    "secrets": [
      {
        "name": "JIRA_BASE_URL",
        "valueFrom": "arn:aws:secretsmanager:REGION:ACCOUNT:secret:sandbox-pm/production/jira-credentials:base_url::"
      },
      {
        "name": "JIRA_EMAIL",
        "valueFrom": "arn:aws:secretsmanager:REGION:ACCOUNT:secret:sandbox-pm/production/jira-credentials:email::"
      },
      {
        "name": "JIRA_API_TOKEN",
        "valueFrom": "arn:aws:secretsmanager:REGION:ACCOUNT:secret:sandbox-pm/production/jira-credentials:api_token::"
      }
    ]
  }]
}
```

### 環境別のシークレット管理

| 環境 | シークレット名 | 用途 |
|------|-------------|------|
| dev | `sandbox-pm/dev/jira-credentials` | 開発用Jiraワークスペース |
| staging | `sandbox-pm/staging/jira-credentials` | ステージングJiraワークスペース |
| production | `sandbox-pm/production/jira-credentials` | 本番Jiraワークスペース |

## IAM ロール設定（最小権限）

ECS タスクロールに以下のポリシーを付与します:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ReadJiraCredentials",
      "Effect": "Allow",
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Resource": [
        "arn:aws:secretsmanager:REGION:ACCOUNT:secret:sandbox-pm/ENVIRONMENT/jira-credentials-*"
      ]
    }
  ]
}
```

> **注意**: `secretsmanager:ListSecrets` や `secretsmanager:DescribeSecret` は不要です。

## 認証情報のローテーション

### 手動ローテーション手順

1. Jira の [API Token 管理画面](https://id.atlassian.com/manage-profile/security/api-tokens) で新しいトークンを発行
2. Secrets Manager でシークレット値を更新:
   ```bash
   aws secretsmanager put-secret-value \
     --secret-id "sandbox-pm/production/jira-credentials" \
     --secret-string '{"base_url":"...","email":"...","api_token":"NEW_TOKEN"}'
   ```
3. ECS タスクを再デプロイ（新しいシークレット値が注入される）
4. 古いトークンを Jira で無効化

### ローテーション推奨サイクル
- 本番環境: 90日ごと
- 担当者異動時: 即時ローテーション

## 監査ログの確認

AWS CloudTrail を有効にすると `GetSecretValue` の呼び出しが記録されます。

```bash
# 過去24時間のシークレットアクセスログを確認
aws cloudtrail lookup-events \
  --lookup-attributes AttributeKey=EventName,AttributeValue=GetSecretValue \
  --start-time $(date -u -v-24H +%Y-%m-%dT%H:%M:%SZ) \
  --query 'Events[*].{Time:EventTime,User:Username,IP:SourceIPAddress}' \
  --output table
```
