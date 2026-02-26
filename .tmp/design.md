# 通知機能 実装設計

## 概要

インアプリ通知センターを追加し、Jira同期の完了・失敗イベントをリアルタイムにユーザーへ通知する。

## 要件

- 通知の種類: インアプリのみ（メール・Slack対象外）
- 通知トリガー: Jira同期完了（成功・失敗）
- 受信者: 全アクティブユーザー（ブロードキャスト）
- UI: サイドバーにベルアイコン + 未読バッジ + 通知パネル（Popover）

## アーキテクチャ

### DB
- notifications テーブル（マイグレーション: 000005）
- user_id → users(id) ON DELETE CASCADE
- type: SYNC_COMPLETED | SYNC_FAILED
- is_read: 既読フラグ（デフォルトFALSE）
- 部分インデックス: notifications_user_unread ON (user_id, is_read) WHERE is_read = FALSE

### Backend
- GET /api/v1/notifications → 自分の通知一覧（未読優先、最大50件）
- PUT /api/v1/notifications/read-all → 全通知既読化
- PUT /api/v1/notifications/:id/read → 特定通知既読化
- triggerSyncHandler のゴルーチン末尾で broadcastSyncNotification を呼び出し

### Frontend
- Zustand store: useNotificationStore（persist なし）
- 30秒間隔ポーリング（Layout.tsx の useEffect で管理）
- NotificationBell: MUI Badge + Popover
- NotificationPanel: 通知一覧 (max-height: 400px)

---

# テストカバレッジ改善 Round 2 (追記)

## Round 2 対象

### Backend (router package) - 残りの低カバレッジハンドラー
| ハンドラー | 現在のカバレッジ | 追加するテスト |
|------------|---------------|--------------|
| createOrganizationHandlerWithDB | 23.1% | success(parent有/無), ParentNotFound, MaxDepth |
| updateOrganizationHandlerWithDB | 35.7% | NotFound, Success |
| deleteOrganizationHandlerWithDB | 55.2% | NotFound, Success |
| assignProjectToOrganizationHandlerWithDB | 25.0% | OrgNotFound, ProjectNotFound, Success, NullOrg |
| updateProjectHandlerWithDB | 0% | InvalidID, MissingIsActive, Success |
| testJiraConnectionHandler | 0% | NotConfigured |

### Frontend
| ファイル | 現状 | 追加するテスト |
|----------|------|--------------|
| `src/utils/permissions.ts` | テストなし | 全4関数の全ロール組み合わせ |
| `src/stores/authStore.ts` | テストなし | login/logout/isAuthenticated |

---

# routerパッケージ テストカバレッジ改善 (Round 1)

## 現状
`go test ./internal/infrastructure/router/... -cover` の結果: **32.3%**

0%カバレッジのハンドラーが多数ある:
- `loginHandler` / `meHandler` (auth_handlers.go)
- `getDashboardSummaryHandlerWithDB` / `getOrganizationSummaryHandlerWithDB`
- `listOrganizationsHandlerWithDB` / `getOrganizationHandlerWithDB` / `getChildOrganizationsHandlerWithDB`
- `listProjectsHandlerWithDB` / `getProjectHandlerWithDB` / `updateProjectHandlerWithDB`

また成功パスのテストが不足しているハンドラー:
- `createUserHandlerWithDB` (33.3%)
- `updateUserHandlerWithDB` (52.0%)

## 追加するテストファイル

### auth_handlers_test.go (新規)
| テスト名 | 期待コード |
|---------|-----------|
| TestLoginHandler_MissingFields | 400 |
| TestLoginHandler_ShortPassword | 401 |
| TestLoginHandler_UserNotFound | 401 |
| TestLoginHandler_DisabledAccount | 401 |
| TestLoginHandler_WrongPassword | 401 |
| TestLoginHandler_Success | 200 |
| TestMeHandler_NoClaims | 401 |
| TestMeHandler_Success | 200 |

### organization_test.go (追記)
| テスト名 | 期待コード |
|---------|-----------|
| TestListOrganizationsHandler_EmptyResult | 200 |
| TestListOrganizationsHandler_ReturnsList | 200 |
| TestGetOrganizationHandler_InvalidID | 400 |
| TestGetOrganizationHandler_NotFound | 404 |
| TestGetOrganizationHandler_Success | 200 |
| TestGetChildOrganizationsHandler_InvalidID | 400 |
| TestGetChildOrganizationsHandler_EmptyResult | 200 |

### dashboard_handlers_test.go (追記)
| テスト名 | 期待コード |
|---------|-----------|
| TestGetDashboardSummaryHandler_Success | 200 |
| TestGetOrganizationSummaryHandler_InvalidID | 400 |
| TestGetOrganizationSummaryHandler_NotFound | 404 |
| TestGetOrganizationSummaryHandler_Success | 200 |

### user_handlers_test.go (追記)
| テスト名 | 期待コード |
|---------|-----------|
| TestCreateUserHandler_Success | 201 |
| TestUpdateUserHandler_Success | 200 |

### project_handlers_test.go (新規)
| テスト名 | 期待コード |
|---------|-----------|
| TestListProjectsHandler_EmptyResult | 200 |
| TestGetProjectHandler_InvalidID | 400 |
| TestGetProjectHandler_NotFound | 404 |
| TestGetProjectHandler_Success | 200 |

## 実装上の注意
- `loginHandler_Success`: bcryptはMinCost(4)で高速なハッシュを生成
- `auth.NewTokenManager("test-secret")` を使用
- sqlmockのクエリマッチは `mock.ExpectQuery` で部分マッチ
