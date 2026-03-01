# DB-004: シードデータの作成 タスクリスト

## Phase 1: ファイル作成

- [x] `database/seeds/` ディレクトリを作成
- [x] `database/seeds/seed.sql` を作成
  - [x] Organizations（12件: 3本部 + 4部 + 5課）
  - [x] Projects（6件）
  - [x] Issues（44件: RED/YELLOW/GREEN混在）
  - [x] Sync Logs（3件: 条件付きINSERTで冪等性確保）
- [x] `database/seeds/apply.sh` を作成
- [x] Makefile に `db-seed` ターゲットを追加

## Phase 2: チケット更新

- [x] `tickets/DB-004_seed-data-creation.md` に `[完了]` を追加

## 確認

- [x] seed.sql 構造確認（44件issues, 6件projects, 12件organizations）
- [x] apply.sh 実行権限付与

## データ内容

| カテゴリ | 件数 |
|---------|-----|
| Organizations | 12 (本部3 + 部4 + 課5) |
| Projects | 6 (5割当済 + 1未分類) |
| Issues | 44 (RED: 12件, YELLOW: 12件, GREEN: 20件) |
| Sync Logs | 3 (SUCCESS×2, FAILURE×1) |
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
