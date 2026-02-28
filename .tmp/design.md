# パスワード変更機能 設計書

## 概要

管理者がユーザーのパスワードをリセットできるエンドポイントとUIを追加する。

## 要件

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
