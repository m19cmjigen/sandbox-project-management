# テストカバレッジ改善 Round 4

## Phase 1: pkg/config/config_test.go (新規)
- [ ] TestLoad_Defaults
- [ ] TestLoad_EnvOverrides
- [ ] TestLoad_InvalidDBPort
- [ ] TestGetDSN

## Phase 2: pkg/logger/logger_test.go (新規)
- [ ] TestNew_JSONFormat
- [ ] TestNew_TextFormat
- [ ] TestNew_InvalidLevel
- [ ] TestWithFields

## Phase 3: frontend api/auth.test.ts (新規)
- [x] login 成功
- [x] login エラー

## Phase 4: 確認
- [x] go test ./pkg/config/... -cover → 100% (目標80%達成)
- [x] go test ./pkg/logger/... -cover → 90% (目標80%達成)
- [x] npm run test -- --run 全通過 (65 tests)

---

# テストカバレッジ改善 Round 3 (完了)

## Phase 1: pkg/auth/middleware_test.go (新規)
- [x] TestGetClaims_NotSet
- [x] TestGetClaims_Set
- [x] TestMiddleware_MissingHeader
- [x] TestMiddleware_BadFormat
- [x] TestMiddleware_InvalidToken
- [x] TestMiddleware_ExpiredToken
- [x] TestMiddleware_ValidToken
- [x] TestRequireRole_NoClaims
- [x] TestRequireRole_WrongRole
- [x] TestRequireRole_AllowedRole

## Phase 2: internal/batch/repository_test.go (新規)
- [x] TestUpsertProjects_Empty
- [x] TestUpsertProjects_Success
- [x] TestUpsertIssues_Empty
- [x] TestUpsertIssues_UnknownProject
- [x] TestUpsertIssues_Success
- [x] TestGetProjectIDMap_Empty
- [x] TestGetProjectIDMap_WithRows
- [x] TestStartSyncLog_Success
- [x] TestGetLastSuccessfulSyncTime_NoRows
- [x] TestGetLastSuccessfulSyncTime_WithRow
- [x] TestFinishSyncLog_Success
- [x] TestFinishSyncLog_WithErrorMessage

## Phase 3: 確認
- [x] go test ./pkg/auth/... -cover → 93.2% (目標80%達成)
- [x] go test ./internal/batch/... -cover → 89.8% (目標75%達成)
