# パスワード変更機能 タスクリスト

## Phase 1: Backend

- [x] `user_handlers.go` に `changePasswordHandlerWithDB` を追加
- [x] `router.go` に `users.PUT("/:id/password", ...)` を追加
- [x] `user_handlers_test.go` に5テスト追加
  - [x] TestChangePasswordHandler_Success (200)
  - [x] TestChangePasswordHandler_InvalidID (400)
  - [x] TestChangePasswordHandler_TooShort (400)
  - [x] TestChangePasswordHandler_UserNotFound (404)
  - [x] TestChangePasswordHandler_DBError (500)

## Phase 2: Frontend

- [x] `frontend/src/api/users.ts` に `changePassword()` を追加
- [x] `frontend/src/pages/UserManagement.tsx` に `ChangePasswordDialog` を追加
- [x] ユーザー一覧各行に `LockResetIcon` ボタンを追加

## 確認

- [x] `cd backend && go test ./internal/infrastructure/router/... -run TestChangePassword -v` → 5テスト全通過
- [x] `cd frontend && npm run type-check` → エラーなし
