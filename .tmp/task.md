# テストカバレッジ改善 TEST-001 (80%+達成)

## Phase 1: issue_handlers_test.go (追記)
- [x] TestListIssuesHandler_WithProjectIDFilter
- [x] TestListIssuesHandler_WithDelayStatusFilter
- [x] TestListIssuesHandler_WithNoDueDateFilter
- [x] TestListIssuesHandler_WithStatusCategoryFilter
- [x] TestListIssuesHandler_WithAssigneeNameFilter
- [x] TestListIssuesHandler_WithDescOrderAndNonDefaultSort
- [x] TestListIssuesHandler_WithInvalidSortFallback
- [x] TestListIssuesHandler_DataSelectError
- [x] TestListProjectIssuesHandler_WithDelayStatusFilter
- [x] TestListProjectIssuesHandler_WithNoDueDateFilter
- [x] TestListProjectIssuesHandler_WithStatusCategoryAndAssigneeFilters
- [x] TestListProjectIssuesHandler_WithDescOrder
- [x] TestListProjectIssuesHandler_WithInvalidSortFallback
- [x] TestListProjectIssuesHandler_DataSelectError

## Phase 2: project_handlers_test.go (追記)
- [x] TestListProjectsHandler_WithUnassignedFilter
- [x] TestListProjectsHandler_WithOrganizationIDFilter
- [x] TestListProjectsHandler_WithDelayCountSort
- [x] TestListProjectsHandler_WithNameDescSort
- [x] TestListProjectsHandler_WithDelayStatusFilter (RED/YELLOW/GREEN)
- [x] TestListProjectsHandler_DataSelectError
- [x] TestListProjectsHandler_WithResults (delay status computation)

## Phase 3: 確認
- [x] go test ./... → 全パッケージ通過
- [x] 全体カバレッジ 80.4% (目標80%+達成)
  - router package: 75.2% → 81.6%
  - listIssuesHandlerWithDB: 70.1% → 100.0%
  - listProjectIssuesHandlerWithDB: 74.6% → 96.8%
  - listProjectsHandlerWithDB: 71.0% → 97.1%

---

# 通知機能実装

## Phase 1: DBマイグレーション
- [x] database/migrations/000005_create_notifications.up.sql 作成
- [x] database/migrations/000005_create_notifications.down.sql 作成

## Phase 2: バックエンド
- [x] notification_handlers.go 作成
  - listNotificationsHandlerWithDB (GET /api/v1/notifications)
  - readNotificationHandlerWithDB (PUT /api/v1/notifications/:id/read)
  - readAllNotificationsHandlerWithDB (PUT /api/v1/notifications/read-all)
  - broadcastSyncNotification (全ユーザーへ通知作成)
- [x] notification_handlers_test.go 作成 (8テスト全通過)
- [x] router.go 更新（通知ルート追加）
- [x] settings_handlers.go 更新（triggerSyncHandler にbroadcastSyncNotification追加）
- [x] go test ./... → 全パッケージ通過

## Phase 3: フロントエンド
- [x] src/api/notifications.ts 作成
- [x] src/stores/notificationStore.ts 作成
- [x] src/components/NotificationBell.tsx 作成
- [x] src/components/NotificationPanel.tsx 作成
- [x] src/components/Layout.tsx 更新（NotificationBell追加、ポーリング開始）
- [x] tsc --noEmit → エラーなし
- [x] npm run test -- --run → 85テスト全通過

---

# テストカバレッジ改善 Round 6

## Phase 1: frontend API unit tests (新規)
- [x] src/api/users.test.ts (getUsers/createUser/updateUser/deleteUser - 8テスト)
- [x] src/api/organizations.test.ts (getOrganizations/getOrganization/createOrganization/updateOrganization/deleteOrganization - 8テスト)
- [x] src/api/dashboard.test.ts (getDashboardSummary/getOrgSummary/getProjectSummary - 4テスト)

## Phase 2: 確認
- [x] npm run test -- --run → 85テスト全通過 (14 test files)

---

# テストカバレッジ改善 Round 5

## Phase 1: router_test.go (新規)
- [x] TestSecurityHeadersMiddleware
- [x] TestCORSMiddleware_Wildcard
- [x] TestCORSMiddleware_AllowedOrigin
- [x] TestCORSMiddleware_DisallowedOrigin
- [x] TestCORSMiddleware_Options
- [x] TestLoggerMiddleware_PassThrough
- [x] TestHealthCheckHandler
- [x] TestReadinessCheckHandler_DBUp
- [x] TestReadinessCheckHandler_DBDown

## Phase 2: settings_handlers_test.go (追記)
- [x] TestUpdateJiraSettingsHandler_InsertSuccess
- [x] TestUpdateJiraSettingsHandler_UpdatePath
- [x] TestTriggerSyncHandler_Success
- [x] TestListSyncLogsHandler_DBError

## Phase 3: 各ハンドラーDBエラーパス (追記)
- [x] TestListUsersHandler_DBError
- [x] TestDeleteUserHandler_DeleteSuccess
- [x] TestListProjectsHandler_DBError
- [x] TestListOrganizationsHandler_DBError
- [x] TestGetDashboardSummaryHandler_DBError
- [x] TestListIssuesHandler_DBError
- [x] TestListProjectIssuesHandler_DBError

## Phase 4: 確認
- [x] go test ./internal/infrastructure/router/... -cover → 75.0% (目標達成)
