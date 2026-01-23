# INFRA-002: Terraform / IaCセットアップ

## 優先度
高

## カテゴリ
Infrastructure

## 説明
インフラ構成をコード管理するため、TerraformまたはAWS CDKのセットアップを行う。

## タスク
- [ ] Terraform / AWS CDKの選定
- [ ] プロジェクト構造の設計（modules, environments等）
- [ ] State管理用のS3バケット作成
- [ ] DynamoDBによるState Lock設定
- [ ] 基本的なVPC/Subnet定義のコード化
- [ ] CI/CDパイプラインでのTerraform実行設定

## 受け入れ基準
- IaCツールが選定され、初期セットアップが完了していること
- State管理が適切に設定されていること
- VPC等の基本リソースがコード化されていること

## 依存関係
- INFRA-001

## 関連ドキュメント
SPEC.md - セクション2.1

## 見積もり工数
3日
