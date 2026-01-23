# INFRA-005: CloudFront + S3 セットアップ (Frontend配信)

## 優先度
中

## カテゴリ
Infrastructure

## 説明
React SPAをホスティングするためのS3バケットとCloudFront Distributionを構築する。

## タスク
- [ ] S3バケットの作成（静的Webホスティング用）
- [ ] バケットポリシーの設定
- [ ] CloudFront Distributionの作成
- [ ] OAI (Origin Access Identity) の設定
- [ ] HTTPSリダイレクトの設定
- [ ] カスタムエラーレスポンス設定（SPA対応）
- [ ] キャッシュポリシーの設定

## 受け入れ基準
- S3バケットにReactアプリをアップロードできること
- CloudFront経由でHTTPSアクセスできること
- SPAのルーティングが正常に動作すること

## 依存関係
- INFRA-001
- INFRA-002

## 関連ドキュメント
SPEC.md - セクション2.2

## 見積もり工数
2日
