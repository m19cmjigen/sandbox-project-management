# テストカバレッジ改善 Round 2

## Phase 1: management_test.go (追記)
- [x] TestCreateOrganizationHandler_ParentNotFound
- [x] TestCreateOrganizationHandler_MaxDepthExceeded
- [x] TestCreateOrganizationHandler_Success (no parent)
- [x] TestCreateOrganizationHandler_SuccessWithParent
- [x] TestUpdateOrganizationHandler_NotFound
- [x] TestUpdateOrganizationHandler_Success
- [x] TestDeleteOrganizationHandler_NotFound
- [x] TestDeleteOrganizationHandler_Success
- [x] TestAssignProjectHandler_OrgNotFound
- [x] TestAssignProjectHandler_ProjectNotFound
- [x] TestAssignProjectHandler_Success
- [x] TestAssignProjectHandler_SuccessNullOrg

## Phase 2: project_handlers_test.go (追記)
- [x] TestUpdateProjectHandler_InvalidID
- [x] TestUpdateProjectHandler_MissingIsActive
- [x] TestUpdateProjectHandler_Success

## Phase 3: settings_handlers_test.go (追記)
- [x] TestTestJiraConnectionHandler_NotConfigured

## Phase 4: frontend tests (新規)
- [x] src/utils/permissions.test.ts (全4関数: 12テスト)
- [x] src/stores/authStore.test.ts (login/logout/isAuthenticated: 6テスト)

## Phase 5: 確認
- [x] go test ./internal/infrastructure/router/... -cover → 65.6% (目標65%達成)
- [x] npm run test -- --run → 63テスト全通過
