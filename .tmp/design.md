# DB-004: シードデータの作成 設計書

## 概要

開発・テスト用のシードデータを作成する。
組織階層・プロジェクト・チケット（RED/YELLOW/GREEN混在）・同期ログを投入する。
# パスワード変更機能 設計書

## 概要

管理者がユーザーのパスワードをリセットできるエンドポイントとUIを追加する。

## 作成するファイル

| ファイル | 内容 |
|---------|-----|
| `database/seeds/seed.sql` | シードデータSQL本体 |
| `database/seeds/apply.sh` | 実行シェルスクリプト |
| `Makefile` | `db-seed` ターゲット追加 |
| `tickets/DB-004_seed-data-creation.md` | 完了マーク追加 |

## データ設計

### Organizations (3階層)
- level 0 (本部): 技術本部, 営業本部, 管理本部
- level 1 (部): 開発部, インフラ部, 営業推進部, 人事部
- level 2 (課): Webシステム課, モバイル開発課, クラウド基盤課, 営業企画課, 人事企画課

### Projects (6件)
- 5件: 各課に紐付け
- 1件: 未分類（organization_id = NULL）

### Issues (約50件)
遅延ステータスはトリガーが自動計算（CURRENT_DATE基準）:
- RED: status_category != 'Done' AND due_date < CURRENT_DATE
- YELLOW: status_category != 'Done' AND due_date <= CURRENT_DATE+3 または due_date IS NULL
- GREEN: 完了済み または 十分余裕のある期限

### Sync Logs (3件)
- SUCCESS × 2, FAILURE × 1
- 対象ロール: Admin のみ
- 現在のパスワード入力は不要（管理者リセット操作のため）
- `auth.HashPassword()` / bcrypt (cost 12) を再利用

## Backend

### エンドポイント

`PUT /api/v1/users/:id/password`

**リクエストボディ:**
```json
{ "new_password": "newpass123" }
```

**レスポンス:**
- 200: `{"message": "password updated"}`
- 400: 不正ID or バリデーションエラー（8文字未満）
- 404: ユーザーが存在しない
- 500: DB エラー

### 実装箇所

- `backend/internal/infrastructure/router/user_handlers.go` に `changePasswordHandlerWithDB` を追加
- `backend/internal/infrastructure/router/router.go` に `users.PUT("/:id/password", ...)` を追加

## Frontend

### API

`frontend/src/api/users.ts` に `changePassword(id, newPassword)` を追加

### UI

`frontend/src/pages/UserManagement.tsx` に以下を追加:
- `ChangePasswordDialog` コンポーネント
- ユーザー一覧の各行に `LockResetIcon` ボタン
- クライアントバリデーション: 8文字以上 & 確認パスワード一致
