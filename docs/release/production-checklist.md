# 本番デプロイ前チェックリスト

本番環境への初回デプロイ時に使用するチェックリストです。
各項目を確認後、担当者がサインオフしてから次のステップに進んでください。

---

## Phase 1: インフラ確認（INFRA チームが担当）

### AWS環境

- [ ] AWS アカウントと権限設定が完了している (INFRA-001)
- [ ] VPC・サブネット・セキュリティグループが設定済み (INFRA-001)
- [ ] Amazon Aurora PostgreSQL クラスターが起動している (INFRA-003)
  - [ ] Multi-AZ 構成になっている
  - [ ] 自動バックアップが有効（保持期間 7日以上）
  - [ ] メンテナンスウィンドウが設定されている
- [ ] AWS Fargate (ECS) クラスター・サービスが設定済み (INFRA-004)
  - [ ] タスク定義（CPU/メモリ・環境変数・ヘルスチェック）が正しく設定されている
  - [ ] ALB（ロードバランサー）とターゲットグループが設定済み
  - [ ] ACM証明書がALBにアタッチされている（HTTPS）
- [ ] CloudFront + S3 が設定済み (INFRA-005)
  - [ ] S3バケットのパブリックアクセスがブロックされている
  - [ ] CloudFrontディストリビューションがS3オリジンに接続されている
  - [ ] HTTPS リダイレクトが設定されている
- [ ] Route 53 で DNS が設定されている
- [ ] ECRリポジトリが作成されている

---

## Phase 2: セキュリティ確認（Security チームが担当）

- [ ] AWS Secrets Manager に本番認証情報が登録されている
  - [ ] DB パスワード
  - [ ] Jira API トークン（将来実装用）
- [ ] ECSタスクロール（IAM）が最小権限で設定されている
- [ ] セキュリティグループがバックエンドポート（8080）のみ許可している
- [ ] RDS がプライベートサブネットに配置されている（外部から直接アクセス不可）
- [ ] セキュリティ監査・脆弱性診断が完了している (SEC-003)
- [ ] APIの認証・認可が実装済みである (SEC-002) ※v1.0.0では未実装のため注意

---

## Phase 3: 環境変数・設定確認

バックエンドのECSタスク定義に以下の環境変数が設定されていることを確認します。

- [ ] `PORT=8080`
- [ ] `GIN_MODE=release`
- [ ] `DB_HOST` — AuroraクラスターのWriterエンドポイント
- [ ] `DB_PORT=5432`
- [ ] `DB_USER` — DB接続ユーザー名
- [ ] `DB_PASSWORD` — Secrets Managerからの参照
- [ ] `DB_NAME` — データベース名
- [ ] `DB_SSLMODE=require`
- [ ] `LOG_LEVEL=info`
- [ ] `LOG_FORMAT=json`

GitHub Secrets に以下が設定されていることを確認します（CDパイプライン用）。

- [ ] `AWS_ACCESS_KEY_ID_PROD`
- [ ] `AWS_SECRET_ACCESS_KEY_PROD`
- [ ] `SUBNET_IDS_PROD`
- [ ] `SG_ID_PROD`
- [ ] `CLOUDFRONT_DISTRIBUTION_PRODUCTION`
- [ ] `PRODUCTION_API_BASE_URL`
- [ ] `JIRA_BASE_URL`

フロントエンドのビルド環境変数を確認します。

- [ ] `VITE_API_BASE_URL` — 本番APIのURL
- [ ] `VITE_JIRA_BASE_URL` — JiraのベースURL（省略可）

---

## Phase 4: データベース準備

- [ ] バックアップが正常に取得できることを確認している
- [ ] マイグレーションをステージング環境で事前テスト済み
- [ ] マイグレーションファイル (`database/migrations/`) に破壊的変更がないことを確認
- [ ] 本番DBに接続してマイグレーションを実行する

```bash
# Aurora PostgreSQLへのマイグレーション実行
DATABASE_URL="postgres://<user>:<pass>@<aurora-writer-endpoint>:5432/<db>?sslmode=require" \
  migrate -path database/migrations -database "$DATABASE_URL" up

# バージョン確認
DATABASE_URL="..." migrate -path database/migrations -database "$DATABASE_URL" version
```

- [ ] マイグレーション後にバージョンが正しいことを確認

---

## Phase 5: デプロイ実行

### バックエンドデプロイ

- [ ] GitHub Actions の CD ワークフローを `workflow_dispatch` で実行（environment: production）
- [ ] 承認者が GitHub Environments の承認を実施
- [ ] ECSデプロイ完了後、新しいタスクが起動していることを確認

```bash
aws ecs describe-services \
  --cluster project-viz-prod \
  --services project-viz-api-production \
  --query 'services[0].{desired:desiredCount,running:runningCount,pending:pendingCount}'
```

### フロントエンドデプロイ

- [ ] S3バケットへのデプロイが完了している
- [ ] CloudFrontキャッシュが無効化されている

---

## Phase 6: 動作確認

### APIヘルスチェック

```bash
# ヘルスチェック
curl https://<api-domain>/health
# 期待値: {"status":"ok","service":"project-visualization-api"}

# レディネスチェック（DB接続確認）
curl https://<api-domain>/ready
# 期待値: {"status":"ready","database":"connected"}

# 主要APIエンドポイントの確認
curl https://<api-domain>/api/v1/dashboard/summary
curl https://<api-domain>/api/v1/organizations
curl https://<api-domain>/api/v1/projects
curl https://<api-domain>/api/v1/issues
```

### フロントエンド動作確認

- [ ] `https://<frontend-domain>/` でダッシュボードが表示される
- [ ] 組織一覧ページが正常に表示される
- [ ] プロジェクト一覧ページが正常に表示される
- [ ] チケット一覧ページが正常に表示される（フィルタ・ソート動作確認）
- [ ] 組織管理ページが正常に表示される
- [ ] HTTPS接続になっている（ブラウザのアドレスバーに鍵マーク）

### パフォーマンス確認

```bash
# 本番環境でスモークテストを実行
BASE_URL=https://<api-domain> ./performance/run.sh smoke
# p(95) < 500ms、エラー率 < 1% を確認
```

---

## Phase 7: 監視・アラート確認

- [ ] CloudWatch ダッシュボードが設定されている
- [ ] ECSコンテナのCPU・メモリアラートが設定されている
- [ ] ALBの5xxエラー率アラートが設定されている
- [ ] Aurora PostgreSQLのDB接続数・CPU・ストレージアラートが設定されている
- [ ] アラート通知先（Slack・メール等）が設定されている

詳細は `docs/release/monitoring.md` を参照してください。

---

## Phase 8: バックアップ確認

- [ ] Aurora自動バックアップが有効になっている（保持期間 7日以上）
- [ ] スナップショット（手動）を取得して、リストア手順を確認している
- [ ] S3バケットのバージョニングが有効になっている（フロントエンドのロールバック用）

---

## Phase 9: リリース完了

- [ ] リリースノート (`docs/release/v1.0.0-release-notes.md`) を最終確認
- [ ] GitHub Release を作成し、タグ `v1.0.0` を付与する
- [ ] 関係者へリリース完了を通知する
- [ ] 運用チームへ引き継ぎを実施する

```bash
# GitHub Release の作成
gh release create v1.0.0 \
  --title "v1.0.0 - 初回リリース" \
  --notes-file docs/release/v1.0.0-release-notes.md
```

---

## ロールバック判断基準

デプロイ後に以下の症状が発生した場合は即座にロールバックを実施します。

| 症状 | 対応 |
|---|---|
| ヘルスチェック (`/health`) が失敗する | バックエンドをロールバック |
| ECSタスクが起動しない | タスク定義を前のリビジョンに戻す |
| APIの5xxエラー率 > 5% | バックエンドをロールバック |
| フロントエンドが真っ白になる | S3を前のビルドに戻す + CloudFront無効化 |
| DBマイグレーションが失敗する | マイグレーションをロールバック後、イメージのロールバック |

ロールバック手順は `docs/deploy.md` の「ロールバック手順」を参照してください。
