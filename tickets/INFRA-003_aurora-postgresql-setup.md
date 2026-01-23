# INFRA-003: Amazon Aurora PostgreSQL セットアップ

## 優先度
高

## カテゴリ
Infrastructure

## 説明
大量のチケットデータと組織階層データを格納するためのAurora PostgreSQLクラスタを構築する。

## タスク
- [ ] Aurora PostgreSQLクラスタの作成（マルチAZ構成）
- [ ] インスタンスタイプの選定（開発環境は小規模、本番は要件に応じてスケール）
- [ ] パラメータグループの設定
- [ ] バックアップ設定（7日間保持）
- [ ] セキュリティグループ設定（Private Subnetからのみアクセス可）
- [ ] Secrets ManagerでDB認証情報を管理
- [ ] 接続テスト

## 受け入れ基準
- Aurora PostgreSQLクラスタが稼働していること
- Private Subnet内からのみアクセス可能な構成になっていること
- Secrets Managerで認証情報が管理されていること
- 接続テストが成功すること

## 依存関係
- INFRA-001
- INFRA-002

## 関連ドキュメント
SPEC.md - セクション2.1

## 見積もり工数
2日
