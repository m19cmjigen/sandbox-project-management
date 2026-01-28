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
- [x] バックエンドのユニットテスト作成（部分完了）
  - [x] **ドメイン層**: 8テストケース（DelayStatus, IssueFilter, Organization等）
  - [x] **リポジトリ層**: 20+テストケース
    - Organization Repository (8 tests)
    - Project Repository (7 tests)
    - Issue Repository (6 tests)
    - User Repository (6 tests)
  - [ ] ビジネスロジック層（TODO）
  - [ ] API層（TODO）
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

### テスト特徴
- データベース接続が利用不可の場合は自動スキップ
- CI環境ではPostgreSQLサービスコンテナで実行
- テスト間の独立性を保証（クリーンアップ処理）
- testify/assertを使用した読みやすいアサーション

### テスト実行方法
```bash
# 全テスト実行
make backend-test

# カバレッジレポート生成
make backend-coverage
```

## 残タスク
1. **ユースケース層テスト**: ビジネスロジックの単体テスト
2. **ハンドラー層テスト**: HTTPハンドラーのテスト
3. **フロントエンドテスト**: Reactコンポーネントのテスト
4. **カバレッジ目標達成**: 全体で80%以上

## 依存関係
各実装チケット

## 見積もり工数
10日（残り: 6日）

## 進捗
完了: 40% (ドメイン層100% + リポジトリ層100%)
