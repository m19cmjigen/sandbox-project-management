# ER図

## エンティティ関連図

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ organizations                                                               │
│─────────────────────────────────────────────────────────────────────────── │
│ PK id           BIGSERIAL                                                   │
│    name         VARCHAR(255) NOT NULL CHECK (name != '')                    │
│ FK parent_id    BIGINT → organizations(id) ON DELETE RESTRICT (NULL=最上位) │
│    path         VARCHAR(1000) NOT NULL  例: /1/5/12/                        │
│    level        INTEGER [0..2]  0=本部, 1=部, 2=課                          │
│    created_at   TIMESTAMP                                                   │
│    updated_at   TIMESTAMP (トリガー自動更新)                                 │
└─────────────────────────────────────────────────────────────────────────────┘
         │ 1
         │ 自己参照 (parent_id)
         │ 0..N
         │
         │ 1
         ▼
┌──────────────────────────────────────────────────────────────────────────┐
│ projects                                                                 │
│──────────────────────────────────────────────────────────────────────── │
│ PK id               BIGSERIAL                                            │
│    jira_project_id  VARCHAR(100) NOT NULL UNIQUE                         │
│    key              VARCHAR(50) NOT NULL  例: PROJ                       │
│    name             VARCHAR(255) NOT NULL                                 │
│    lead_account_id  VARCHAR(100) NULL                                    │
│    lead_email       VARCHAR(255) NULL                                    │
│ FK organization_id  BIGINT → organizations(id) ON DELETE SET NULL        │
│    created_at       TIMESTAMP                                            │
│    updated_at       TIMESTAMP (トリガー自動更新)                          │
└──────────────────────────────────────────────────────────────────────────┘
         │ 1
         │
         │ 0..N
         ▼
┌──────────────────────────────────────────────────────────────────────────┐
│ issues                                                                   │
│──────────────────────────────────────────────────────────────────────── │
│ PK id                  BIGSERIAL                                         │
│    jira_issue_id       VARCHAR(100) NOT NULL UNIQUE                      │
│    jira_issue_key      VARCHAR(100) NOT NULL  例: PROJ-123               │
│ FK project_id          BIGINT → projects(id) ON DELETE CASCADE           │
│    summary             TEXT NOT NULL                                     │
│    status              VARCHAR(100) NOT NULL                             │
│    status_category     VARCHAR(50) IN ('To Do','In Progress','Done')     │
│    due_date            DATE NULL                                         │
│    assignee_name       VARCHAR(255) NULL                                 │
│    assignee_account_id VARCHAR(100) NULL                                 │
│    delay_status        VARCHAR(20) IN ('RED','YELLOW','GREEN') DEFAULT 'GREEN' │
│    priority            VARCHAR(50) NULL                                  │
│    issue_type          VARCHAR(100) NULL                                 │
│    last_updated_at     TIMESTAMP NOT NULL                                │
│    created_at          TIMESTAMP                                         │
│    updated_at          TIMESTAMP (トリガー自動更新)                       │
└──────────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────────┐
│ sync_logs                                                                │
│──────────────────────────────────────────────────────────────────────── │
│ PK id               BIGSERIAL                                            │
│    sync_type        VARCHAR(20) IN ('FULL','DELTA')                      │
│    executed_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP                  │
│    completed_at     TIMESTAMP NULL                                       │
│    status           VARCHAR(20) IN ('RUNNING','SUCCESS','FAILURE')       │
│    projects_synced  INTEGER DEFAULT 0                                    │
│    issues_synced    INTEGER DEFAULT 0                                    │
│    error_message    TEXT NULL                                            │
│    duration_seconds INTEGER NULL                                         │
└──────────────────────────────────────────────────────────────────────────┘
```

## ビュー

### project_delay_summary

プロジェクトごとのチケット遅延集計ビュー。

| カラム | 説明 |
|---|---|
| project_id | プロジェクトID |
| jira_project_id | JiraプロジェクトID |
| project_key | プロジェクトキー |
| project_name | プロジェクト名 |
| organization_id | 所属組織ID |
| total_issues | 総チケット数 |
| red_issues | 遅延チケット数 (RED) |
| yellow_issues | 要注意チケット数 (YELLOW) |
| green_issues | 正常チケット数 (GREEN) |
| open_issues | 未完了チケット数 |
| done_issues | 完了チケット数 |

### organization_delay_summary

組織ごとの遅延プロジェクト集計ビュー。

| カラム | 説明 |
|---|---|
| organization_id | 組織ID |
| organization_name | 組織名 |
| path | 階層パス |
| level | 階層レベル |
| total_projects | 総プロジェクト数 |
| delayed_projects | 遅延プロジェクト数 |
| total_red_issues | 遅延チケット総数 |
| total_yellow_issues | 要注意チケット総数 |
| total_green_issues | 正常チケット総数 |
| total_open_issues | 未完了チケット総数 |

## トリガー

### update_updated_at_column

organizations、projects、issuesテーブルの`updated_at`カラムをUPDATE時に自動更新します。

### calculate_issue_delay_status

issuesテーブルへのINSERT/UPDATE時に`delay_status`を自動計算します。

計算ロジック:
- `status_category != 'Done'` かつ `due_date < CURRENT_DATE` → `RED`
- `status_category != 'Done'` かつ `due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + 3 days` → `YELLOW`
- `status_category != 'Done'` かつ `due_date IS NULL` → `YELLOW`
- それ以外 → `GREEN`

## インデックス

| テーブル | インデックス | カラム |
|---|---|---|
| organizations | idx_organizations_parent_id | parent_id |
| organizations | idx_organizations_path | path |
| organizations | idx_organizations_level | level |
| projects | idx_projects_jira_project_id (UNIQUE) | jira_project_id |
| projects | idx_projects_organization_id | organization_id |
| projects | idx_projects_key | key |
| issues | idx_issues_jira_issue_id (UNIQUE) | jira_issue_id |
| issues | idx_issues_jira_issue_key | jira_issue_key |
| issues | idx_issues_project_id | project_id |
| issues | idx_issues_delay_status | delay_status |
| issues | idx_issues_status_category | status_category |
| issues | idx_issues_due_date | due_date |
| issues | idx_issues_updated_at | last_updated_at |
| issues | idx_issues_project_delay | (project_id, delay_status) |
| sync_logs | idx_sync_logs_executed_at | executed_at DESC |
| sync_logs | idx_sync_logs_status | status |
| sync_logs | idx_sync_logs_sync_type | sync_type |
