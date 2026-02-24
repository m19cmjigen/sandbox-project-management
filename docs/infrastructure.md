# AWSインフラ設計書

## 概要

本プラットフォームはAWS上に構築する。インフラはTerraformでコード管理する（IaCコードは `infra/terraform/` 以下）。

---

## 1. 環境構成

| 環境 | 用途 | Stateファイルパス |
|------|------|------------------|
| staging | 開発・結合テスト | `staging/` |
| production | 本番運用 | `production/` |

---

## 2. ネットワーク構成 (VPC)

### CIDR設計

| 環境 | VPC CIDR |
|------|----------|
| staging | 10.0.0.0/16 |
| production | 10.1.0.0/16 |

### サブネット構成（各環境共通）

| サブネット | AZ-a | AZ-c | 用途 |
|-----------|------|------|------|
| Public | 10.x.0.0/24 | 10.x.1.0/24 | ALB、NAT GW |
| Private | 10.x.10.0/24 | 10.x.11.0/24 | ECS, Aurora |

- ECSタスク・AuroraはPrivate Subnetに配置し、外部から直接アクセスできない構成とする
- 各AZにNAT Gatewayを配置し、PrivateサブネットからのアウトバウンドはNAT GW経由とする

### セキュリティグループ

| SG名 | Inbound | 用途 |
|------|---------|------|
| alb_sg | 80, 443 from 0.0.0.0/0 | ALB |
| api_sg | 8080 from alb_sg | ECS APIタスク |
| db_sg | 5432 from api_sg | Aurora PostgreSQL |

バッチタスクはAPIと同じPrivate Subnetに配置し、api_sgを通じてDBにアクセスする。

---

## 3. IAMロール設計

| ロール名 | Assume Entity | 権限概要 |
|---------|--------------|---------|
| `ecs-task-execution` | ECS | ECRプル、CloudWatch Logs書き込み、Secrets Manager読み取り |
| `api-task` | ECS | Secrets Manager読み取り（DB認証情報・JWT_SECRET） |
| `batch-task` | ECS | Secrets Manager読み取り（DB認証情報・Jira認証情報） |
| `github-actions` | GitHub OIDC | ECRプッシュ、ECSデプロイ、S3同期、CloudFrontキャッシュ無効化 |

### GitHub Actions OIDC認証

長期クレデンシャル（Access Key）を使わず、OIDCによるフェデレーション認証を採用する。

```yaml
# GitHub Actionsワークフローでの利用例
- uses: aws-actions/configure-aws-credentials@v4
  with:
    role-to-assume: arn:aws:iam::<account-id>:role/project-viz-staging-github-actions
    aws-region: ap-northeast-1
```

---

## 4. CloudWatch Logsポリシー

| ロググループ | 保持期間（staging） | 保持期間（production） |
|------------|---------------------|------------------------|
| `/ecs/project-viz-staging-api` | 30日 | - |
| `/ecs/project-viz-production-api` | - | 90日 |
| `/ecs/project-viz-staging-batch` | 30日 | - |
| `/ecs/project-viz-production-batch` | - | 90日 |

---

## 5. Terraformディレクトリ構造

```
infra/terraform/
├── modules/
│   ├── vpc/          # VPC・サブネット・セキュリティグループ・CloudWatch Logs
│   ├── iam/          # IAMロール・ポリシー・GitHub OIDC
│   ├── aurora/       # Aurora PostgreSQL (INFRA-003)
│   ├── ecs/          # ECS Cluster・ALB・サービス (INFRA-004)
│   └── cloudfront/   # CloudFront + S3 (INFRA-005)
└── environments/
    ├── staging/
    │   ├── main.tf
    │   ├── variables.tf
    │   └── outputs.tf
    └── production/
        ├── main.tf
        ├── variables.tf
        └── outputs.tf
```

---

## 6. 初回セットアップ手順

### 前提

- AWS CLIが設定済みであること (`aws configure`)
- Terraform 1.6以上がインストールされていること

### Step 1: Terraformバックエンド用リソース作成 (INFRA-002)

```bash
# S3バケット（Stateファイル格納）
aws s3api create-bucket \
  --bucket project-viz-terraform-state \
  --region ap-northeast-1 \
  --create-bucket-configuration LocationConstraint=ap-northeast-1

# バージョニング有効化（誤削除対策）
aws s3api put-bucket-versioning \
  --bucket project-viz-terraform-state \
  --versioning-configuration Status=Enabled

# サーバーサイド暗号化
aws s3api put-bucket-encryption \
  --bucket project-viz-terraform-state \
  --server-side-encryption-configuration '{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}'

# DynamoDB（Stateロック）
aws dynamodb create-table \
  --table-name project-viz-terraform-lock \
  --attribute-definitions AttributeName=LockID,AttributeType=S \
  --key-schema AttributeName=LockID,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST \
  --region ap-northeast-1
```

### Step 2: GitHub OIDC Providerの作成（アカウントにつき1回）

```bash
aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com \
  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1
```

または、Terraformの `create_github_oidc_provider = true` で作成する（初回staging環境デプロイ時のみ）。

### Step 3: stagingのデプロイ

```bash
cd infra/terraform/environments/staging

terraform init
terraform plan -var="aws_account_id=<your-account-id>"
terraform apply -var="aws_account_id=<your-account-id>"
```

### Step 4: productionのデプロイ

```bash
cd infra/terraform/environments/production

terraform init
terraform plan -var="aws_account_id=<your-account-id>"
terraform apply -var="aws_account_id=<your-account-id>"
```

---

## 7. セキュリティ設計方針

- **最小権限原則**: 各IAMロールには必要最小限の権限のみを付与する
- **ネットワーク分離**: ECS・AuroraはPrivate Subnetに配置し、インターネットからの直接アクセスを遮断する
- **認証情報管理**: DB接続情報・JWTシークレット・Jira APIトークンはAWS Secrets Managerで管理し、ECSタスク起動時に環境変数として注入する
- **OIDC認証**: CI/CDはOIDCベースのロール認証を使用し、長期クレデンシャル（Access Key）を使わない
- **暗号化**: S3（Stateファイル）はAES256で暗号化する
