# API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

現在のバージョンでは認証は未実装です。将来的にはJWT認証を追加予定です。

## Response Format

### Success Response

```json
{
  "data": { ... },
  "message": "Success message (optional)"
}
```

### Error Response

```json
{
  "error": "Error message"
}
```

## Endpoints

### Health Check

#### GET /health

ヘルスチェックエンドポイント

**Response:**
```json
{
  "status": "ok",
  "service": "project-visualization-api"
}
```

#### GET /ready

Readinessチェックエンドポイント

**Response:**
```json
{
  "status": "ready",
  "database": "connected"
}
```

---

## Organizations (組織)

### GET /organizations

全組織を取得

**Response:**
```json
{
  "organizations": [
    {
      "id": 1,
      "name": "開発部",
      "parent_id": null,
      "path": "1",
      "level": 0,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### GET /organizations/tree

組織をツリー構造で取得

**Response:**
```json
{
  "tree": [
    {
      "id": 1,
      "name": "開発部",
      "parent_id": null,
      "path": "1",
      "level": 0,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "children": [
        {
          "id": 2,
          "name": "フロントエンドチーム",
          "parent_id": 1,
          "path": "1.2",
          "level": 1,
          "created_at": "2024-01-01T00:00:00Z",
          "updated_at": "2024-01-01T00:00:00Z",
          "children": []
        }
      ]
    }
  ]
}
```

### GET /organizations/:id

組織詳細を取得

**Parameters:**
- `id` (path, required): 組織ID

**Response:**
```json
{
  "id": 1,
  "name": "開発部",
  "parent_id": null,
  "path": "1",
  "level": 0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### GET /organizations/:id/children

子組織を取得

**Parameters:**
- `id` (path, required): 組織ID

**Response:**
```json
{
  "children": [
    {
      "id": 2,
      "name": "フロントエンドチーム",
      "parent_id": 1,
      "path": "1.2",
      "level": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### POST /organizations

組織を作成

**Request Body:**
```json
{
  "name": "新しい組織",
  "parent_id": 1
}
```

**Response:**
```json
{
  "id": 3,
  "name": "新しい組織",
  "parent_id": 1,
  "path": "1.3",
  "level": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### PUT /organizations/:id

組織を更新

**Parameters:**
- `id` (path, required): 組織ID

**Request Body:**
```json
{
  "name": "更新された組織名",
  "parent_id": 1
}
```

**Response:**
```json
{
  "id": 3,
  "name": "更新された組織名",
  "parent_id": 1,
  "path": "1.3",
  "level": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### DELETE /organizations/:id

組織を削除

**Parameters:**
- `id` (path, required): 組織ID

**Response:**
```json
{
  "message": "Organization deleted successfully"
}
```

---

## Projects (プロジェクト)

### GET /projects

プロジェクト一覧を取得

**Query Parameters:**
- `with_stats` (boolean, optional): 統計情報を含めるか（デフォルト: false）
- `organization_id` (number, optional): 組織IDでフィルタ
- `unassigned` (boolean, optional): 未割り当てプロジェクトのみ

**Response (with_stats=false):**
```json
{
  "projects": [
    {
      "id": 1,
      "jira_project_id": "10001",
      "key": "PROJ",
      "name": "プロジェクト名",
      "lead_account_id": "user123",
      "lead_email": "user@example.com",
      "organization_id": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**Response (with_stats=true):**
```json
{
  "projects": [
    {
      "id": 1,
      "jira_project_id": "10001",
      "key": "PROJ",
      "name": "プロジェクト名",
      "lead_account_id": "user123",
      "lead_email": "user@example.com",
      "organization_id": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "total_issues": 100,
      "red_issues": 5,
      "yellow_issues": 10,
      "green_issues": 70,
      "done_issues": 15
    }
  ]
}
```

### GET /projects/:id

プロジェクト詳細を取得

**Parameters:**
- `id` (path, required): プロジェクトID

**Query Parameters:**
- `with_stats` (boolean, optional): 統計情報を含めるか

**Response:**
```json
{
  "id": 1,
  "jira_project_id": "10001",
  "key": "PROJ",
  "name": "プロジェクト名",
  "lead_account_id": "user123",
  "lead_email": "user@example.com",
  "organization_id": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### PUT /projects/:id/organization

プロジェクトを組織に割り当て

**Parameters:**
- `id` (path, required): プロジェクトID

**Request Body:**
```json
{
  "organization_id": 1
}
```

**Response:**
```json
{
  "message": "Project assigned to organization successfully"
}
```

### GET /projects/:id/issues

プロジェクトのIssue一覧を取得

**Parameters:**
- `id` (path, required): プロジェクトID

**Response:**
```json
{
  "issues": [
    {
      "id": 1,
      "jira_issue_id": "10001",
      "jira_issue_key": "PROJ-123",
      "project_id": 1,
      "summary": "Issue概要",
      "status": "In Progress",
      "status_category": "in_progress",
      "due_date": "2024-12-31T00:00:00Z",
      "assignee_name": "田中太郎",
      "assignee_account_id": "user123",
      "delay_status": "GREEN",
      "priority": "High",
      "issue_type": "Task",
      "last_updated_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

---

## Issues (チケット)

### GET /issues

Issue一覧を取得

**Query Parameters:**
- `project_id` (number, optional): プロジェクトIDでフィルタ
- `delay_status` (string, optional): 遅延ステータスでフィルタ（RED/YELLOW/GREEN）
- `assignee` (string, optional): 担当者名でフィルタ
- `status` (string, optional): ステータスでフィルタ

**Response:**
```json
{
  "issues": [
    {
      "id": 1,
      "jira_issue_id": "10001",
      "jira_issue_key": "PROJ-123",
      "project_id": 1,
      "summary": "Issue概要",
      "status": "In Progress",
      "status_category": "in_progress",
      "due_date": "2024-12-31T00:00:00Z",
      "assignee_name": "田中太郎",
      "assignee_account_id": "user123",
      "delay_status": "GREEN",
      "priority": "High",
      "issue_type": "Task",
      "last_updated_at": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### GET /issues/:id

Issue詳細を取得

**Parameters:**
- `id` (path, required): IssueID

**Response:**
```json
{
  "id": 1,
  "jira_issue_id": "10001",
  "jira_issue_key": "PROJ-123",
  "project_id": 1,
  "summary": "Issue概要",
  "status": "In Progress",
  "status_category": "in_progress",
  "due_date": "2024-12-31T00:00:00Z",
  "assignee_name": "田中太郎",
  "assignee_account_id": "user123",
  "delay_status": "GREEN",
  "priority": "High",
  "issue_type": "Task",
  "last_updated_at": "2024-01-01T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

## Dashboard (ダッシュボード)

### GET /dashboard/summary

全体サマリーを取得

**Response:**
```json
{
  "total_projects": 10,
  "total_issues": 500,
  "red_issues": 50,
  "yellow_issues": 100,
  "green_issues": 300,
  "done_issues": 50,
  "projects_by_status": [
    {
      "id": 1,
      "name": "プロジェクトA",
      "key": "PROJA",
      "total_issues": 50,
      "red_issues": 5,
      "yellow_issues": 10,
      "green_issues": 30,
      "done_issues": 5
    }
  ]
}
```

### GET /dashboard/organizations/:id

組織別サマリーを取得

**Parameters:**
- `id` (path, required): 組織ID

**Response:**
```json
{
  "organization_id": 1,
  "organization_name": "開発部",
  "total_projects": 5,
  "total_issues": 250,
  "red_issues": 25,
  "yellow_issues": 50,
  "green_issues": 150,
  "done_issues": 25,
  "projects": [
    {
      "id": 1,
      "name": "プロジェクトA",
      "key": "PROJA",
      "total_issues": 50,
      "red_issues": 5,
      "yellow_issues": 10,
      "green_issues": 30,
      "done_issues": 5
    }
  ]
}
```

### GET /dashboard/projects/:id

プロジェクト別サマリーを取得

**Parameters:**
- `id` (path, required): プロジェクトID

**Response:**
```json
{
  "project_id": 1,
  "project_name": "プロジェクトA",
  "project_key": "PROJA",
  "total_issues": 50,
  "red_issues": 5,
  "yellow_issues": 10,
  "green_issues": 30,
  "done_issues": 5,
  "issues_by_status": {
    "To Do": 10,
    "In Progress": 25,
    "Done": 15
  }
}
```

---

## Sync (Jira同期)

### POST /sync/trigger

Jira同期を手動トリガー

**Request Body:**
```json
{
  "organization_id": 1
}
```

**Response:**
```json
{
  "message": "Sync completed",
  "sync_log": {
    "id": 1,
    "started_at": "2024-01-01T00:00:00Z",
    "completed_at": "2024-01-01T00:05:00Z",
    "status": "COMPLETED",
    "projects_synced": 10,
    "issues_synced": 500,
    "error_count": 0,
    "error_message": null,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:05:00Z"
  }
}
```

### POST /sync/projects/:id

特定プロジェクトを同期

**Parameters:**
- `id` (path, required): プロジェクトID

**Response:**
```json
{
  "message": "Project sync completed",
  "project_id": 1
}
```

### GET /sync/logs

同期ログ一覧を取得

**Query Parameters:**
- `limit` (number, optional): 取得件数（デフォルト: 20、最大: 100）

**Response:**
```json
{
  "logs": [
    {
      "id": 1,
      "started_at": "2024-01-01T00:00:00Z",
      "completed_at": "2024-01-01T00:05:00Z",
      "status": "COMPLETED",
      "projects_synced": 10,
      "issues_synced": 500,
      "error_count": 0,
      "error_message": null,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:05:00Z"
    }
  ]
}
```

### GET /sync/logs/latest

最新の同期ログを取得

**Response:**
```json
{
  "id": 1,
  "started_at": "2024-01-01T00:00:00Z",
  "completed_at": "2024-01-01T00:05:00Z",
  "status": "COMPLETED",
  "projects_synced": 10,
  "issues_synced": 500,
  "error_count": 0,
  "error_message": null,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:05:00Z"
}
```

### GET /sync/logs/:id

同期ログ詳細を取得

**Parameters:**
- `id` (path, required): ログID

**Response:**
```json
{
  "id": 1,
  "started_at": "2024-01-01T00:00:00Z",
  "completed_at": "2024-01-01T00:05:00Z",
  "status": "COMPLETED",
  "projects_synced": 10,
  "issues_synced": 500,
  "error_count": 0,
  "error_message": null,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:05:00Z"
}
```

---

## Data Models

### Organization

| Field | Type | Description |
|-------|------|-------------|
| id | number | 組織ID |
| name | string | 組織名 |
| parent_id | number\|null | 親組織ID |
| path | string | 階層パス（例: "1.2.3"） |
| level | number | 階層レベル（0始まり） |
| created_at | string | 作成日時 |
| updated_at | string | 更新日時 |

### Project

| Field | Type | Description |
|-------|------|-------------|
| id | number | プロジェクトID |
| jira_project_id | string | JiraプロジェクトID |
| key | string | プロジェクトキー |
| name | string | プロジェクト名 |
| lead_account_id | string\|null | リーダーアカウントID |
| lead_email | string\|null | リーダーメール |
| organization_id | number\|null | 組織ID |
| created_at | string | 作成日時 |
| updated_at | string | 更新日時 |

### ProjectWithStats

Project + 以下の統計情報

| Field | Type | Description |
|-------|------|-------------|
| total_issues | number | 総Issue数 |
| red_issues | number | 遅延Issue数 |
| yellow_issues | number | 注意Issue数 |
| green_issues | number | 正常Issue数 |
| done_issues | number | 完了Issue数 |

### Issue

| Field | Type | Description |
|-------|------|-------------|
| id | number | IssueID |
| jira_issue_id | string | JiraのIssue ID |
| jira_issue_key | string | Jiraのキー（例: PROJ-123） |
| project_id | number | プロジェクトID |
| summary | string | Issue概要 |
| status | string | ステータス |
| status_category | string | ステータスカテゴリ |
| due_date | string\|null | 期日 |
| assignee_name | string\|null | 担当者名 |
| assignee_account_id | string\|null | 担当者アカウントID |
| delay_status | DelayStatus | 遅延ステータス |
| priority | string\|null | 優先度 |
| issue_type | string\|null | Issueタイプ |
| last_updated_at | string | 最終更新日時 |
| created_at | string | 作成日時 |
| updated_at | string | 更新日時 |

### DelayStatus

遅延ステータスの列挙型:
- `RED`: 遅延（期日超過）
- `YELLOW`: 注意（期日まで3日以内）
- `GREEN`: 正常

### SyncLog

| Field | Type | Description |
|-------|------|-------------|
| id | number | ログID |
| started_at | string | 開始日時 |
| completed_at | string\|null | 完了日時 |
| status | SyncStatus | 同期ステータス |
| projects_synced | number | 同期したプロジェクト数 |
| issues_synced | number | 同期したIssue数 |
| error_count | number | エラー数 |
| error_message | string\|null | エラーメッセージ |
| created_at | string | 作成日時 |
| updated_at | string | 更新日時 |

### SyncStatus

同期ステータスの列挙型:
- `RUNNING`: 実行中
- `COMPLETED`: 完了
- `COMPLETED_WITH_ERRORS`: 完了（エラーあり）
- `FAILED`: 失敗

---

## Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - リクエストが不正 |
| 404 | Not Found - リソースが見つからない |
| 500 | Internal Server Error - サーバーエラー |
| 503 | Service Unavailable - サービス利用不可 |
