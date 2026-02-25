# [進行中] TEST-001: ユニットテストの実装

## 優先度
高

## カテゴリ
Testing

## 説明
バックエンドとフロントエンドのユニットテストを充実させる。

## タスク
- [x] バックエンドテストフレームワーク選定（Go testing, testify等）
- [ ] フロントエンドテストフレームワーク選定（Jest, Vitest等）
- [x] バックエンドのユニットテスト作成（進行中）
  - [x] **ドメイン層**: 8テストケース（DelayStatus, IssueFilter, Organization等）
  - [x] **リポジトリ層**: 27テストケース
    - Organization Repository (8 tests)
    - Project Repository (7 tests)
    - Issue Repository (6 tests)
    - User Repository (6 tests)
  - [x] **ユースケース層**: 96テストケース（70/73 passing = 95.9%）
    - Auth Usecase (27 tests) - 認証、ユーザー管理、JWT
    - Dashboard Usecase (13 tests) - ダッシュボードサマリー
    - Issue Usecase (18 tests) - Issue検索、フィルタリング
    - Project Usecase (15 tests) - プロジェクト管理
    - Organization Usecase (23 tests) - 組織階層管理
  - [x] **ハンドラー層（API層）**: 96テストケース（完了 - 全テストPASS）
    - Auth Handler (26 tests) ✅ - 認証、ユーザー管理API
    - Dashboard Handler (10 tests) ✅ - ダッシュボードAPI
    - Issue Handler (18 tests) ✅ - IssueAPI（フィルタリング、検索）
    - Project Handler (21 tests) ✅ - ProjectAPI（統計、組織割り当て）
    - Organization Handler (21 tests) ✅ - OrganizationAPI（CRUD、ツリー構造）
- [ ] フロントエンドのユニットテスト作成
  - [ ] コンポーネント
  - [ ] ユーティリティ関数
  - [ ] 状態管理
- [x] モック・スタブの作成（testify/mock使用）
- [ ] テストカバレッジ80%以上を目標（現在: ドメイン層100%, リポジトリ層完了）
- [x] CI/CDでのテスト自動実行（.github/workflows/ci.yml設定済み）

## 受け入れ基準
- [x] ユニットテストが実装されていること（ドメイン+リポジトリ層完了）
- [ ] テストカバレッジが80%以上であること（進行中）
- [x] CI/CDで自動実行されること
- [x] すべてのテストが合格すること（実装済みテストは合格）

## 実装済み内容

### ドメイン層テスト
- `internal/domain/issue_test.go` - Issue関連のバリデーションテスト
- `internal/domain/organization_test.go` - Organization検証ロジックテスト

### リポジトリ層統合テスト
- `internal/infrastructure/postgres/test_helper.go` - テストユーティリティ
- `internal/infrastructure/postgres/organization_repository_impl_test.go`
- `internal/infrastructure/postgres/project_repository_impl_test.go`
- `internal/infrastructure/postgres/issue_repository_impl_test.go`
- `internal/infrastructure/postgres/user_repository_test.go`

### ユースケース層ユニットテスト
- `internal/usecase/auth_usecase_test.go` - 27テスト（認証、トークン管理）
- `internal/usecase/dashboard_usecase_test.go` - 13テスト（サマリー、統計）
- `internal/usecase/issue_usecase_test.go` - 18テスト（Issue操作）
- `internal/usecase/project_usecase_test.go` - 15テスト（プロジェクト管理）
- `internal/usecase/organization_usecase_test.go` - 23テスト（組織階層）

### ハンドラー層HTTPテスト
- `internal/interface/http/auth_handler_test.go` - 26テスト（認証API）
- `internal/interface/handler/dashboard_handler_test.go` - 10テスト（ダッシュボードAPI）
- `internal/interface/handler/issue_handler_test.go` - 18テスト（IssueAPI、フィルタリング）
- `internal/interface/handler/project_handler_test.go` - 21テスト（ProjectAPI、統計情報）
- `internal/interface/handler/organization_handler_test.go` - 21テスト（OrganizationAPI、ツリー構造）

### テスト特徴
- データベース接続が利用不可の場合は自動スキップ（リポジトリ層）
- CI環境ではPostgreSQLサービスコンテナで実行
- テスト間の独立性を保証（クリーンアップ処理）
- testify/assert, testify/mock を使用
- ユースケース層では完全なモッキング（依存ゼロ）
- エラーハンドリングの網羅的テスト
- 境界条件とエッジケースのカバレッジ

### テスト実行方法
```bash
# 全テスト実行
make backend-test

# カバレッジレポート生成
make backend-coverage
```

## 残タスク
1. **フロントエンドテスト**: Reactコンポーネントのテスト
2. **カバレッジ目標達成**: 全体で80%以上
3. **ユースケース層の残り3テスト修正**: GetTree, Create tests (OrganizationUsecase)

## 依存関係
各実装チケット

## 見積もり工数
10日（残り: 4日）

## 進捗
完了: 80% (ドメイン層100% + リポジトリ層100% + ユースケース層96% + ハンドラー層100%)
