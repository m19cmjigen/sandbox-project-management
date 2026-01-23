# BACK-003: 組織管理APIの実装

## 優先度
高

## カテゴリ
Backend

## 説明
組織の作成・編集・削除・階層取得を行うREST APIを実装する。

## タスク
- [ ] GET /api/organizations - 全組織階層の取得
- [ ] GET /api/organizations/:id - 特定組織の詳細取得
- [ ] POST /api/organizations - 組織の作成
- [ ] PUT /api/organizations/:id - 組織の更新
- [ ] DELETE /api/organizations/:id - 組織の削除
- [ ] バリデーション処理の実装
- [ ] エラーハンドリング
- [ ] APIドキュメント作成（Swagger/OpenAPI）
- [ ] 統合テスト

## 受け入れ基準
- すべてのエンドポイントが実装されていること
- APIドキュメントが作成されていること
- 統合テストが合格していること

## 依存関係
- BACK-002

## 関連ドキュメント
SPEC.md - セクション3.3

## 見積もり工数
3日
