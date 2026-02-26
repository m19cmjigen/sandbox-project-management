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
