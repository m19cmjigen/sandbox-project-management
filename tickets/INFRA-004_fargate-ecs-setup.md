# INFRA-004: AWS Fargate (ECS) セットアップ

## 優先度
高

## カテゴリ
Infrastructure

## 説明
Backend APIとBatch Workerを実行するためのECS/Fargateクラスタを構築する。

## タスク
- [ ] ECSクラスタの作成
- [ ] Fargateタスク定義テンプレートの作成
- [ ] ALB (Application Load Balancer) の構築（API用）
- [ ] ターゲットグループの設定
- [ ] ECS Serviceの作成（API用）
- [ ] CloudWatch Logsとの連携設定
- [ ] Auto Scaling設定（本番環境）

## 受け入れ基準
- ECS/Fargateクラスタが作成されていること
- ALBが構築され、ヘルスチェックが機能すること
- CloudWatch Logsにログが出力されること

## 依存関係
- INFRA-001
- INFRA-002

## 関連ドキュメント
SPEC.md - セクション2.1, 2.2

## 見積もり工数
3日
