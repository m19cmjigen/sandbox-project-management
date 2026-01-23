# データベーススキーマ設計書

## 概要

全社プロジェクト進捗可視化プラットフォームのデータベーススキーマ設計。
PostgreSQLを想定した設計。

## ER図

```
Organizations (組織マスタ)
  ├─ id (PK)
  ├─ name
  ├─ parent_id (FK → Organizations.id)
  └─ path (階層検索用)
      │
      ├── Projects (Jiraプロジェクト)
      │     ├─ id (PK)
      │     ├─ jira_project_id (UNIQUE)
      │     ├─ key
      │     ├─ name
      │     ├─ lead_account_id
      │     ├─ organization_id (FK → Organizations.id)
      │     └─ created_at, updated_at
      │         │
      │         └── Issues (チケット情報)
      │               ├─ id (PK)
      │               ├─ jira_issue_id (UNIQUE)
      │               ├─ project_id (FK → Projects.id)
      │               ├─ summary
      │               ├─ status
      │               ├─ status_category
      │               ├─ due_date
      │               ├─ assignee_name
      │               ├─ assignee_account_id
      │               ├─ delay_status (RED/YELLOW/GREEN)
      │               ├─ last_updated_at
      │               └─ created_at, updated_at
      │
      └── SyncLogs (バッチ実行ログ)
            ├─ id (PK)
            ├─ sync_type (FULL/DELTA)
            ├─ executed_at
            ├─ status (SUCCESS/FAILURE/RUNNING)
            ├─ projects_synced
            ├─ issues_synced
            ├─ error_message
            └─ completed_at
```

## テーブル詳細設計

### 1. Organizations (組織マスタ)

組織階層を管理するテーブル。本部-部-課の3階層を想定。

| カラム名 | データ型 | NULL | デフォルト値 | 説明 |
|---------|---------|------|------------|------|
| id | BIGSERIAL | NOT NULL | - | 主キー |
| name | VARCHAR(255) | NOT NULL | - | 組織名 |
| parent_id | BIGINT | NULL | NULL | 親組織ID（最上位はNULL） |
| path | VARCHAR(1000) | NOT NULL | - | 階層パス（例: /1/5/12/） |
| level | INTEGER | NOT NULL | 0 | 階層レベル（0:本部, 1:部, 2:課） |
| created_at | TIMESTAMP | NOT NULL | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | NOT NULL | CURRENT_TIMESTAMP | 更新日時 |

**主キー**: id
**外部キー**: parent_id → Organizations(id)
**インデックス**:
- `idx_organizations_parent_id` on parent_id
- `idx_organizations_path` on path (階層検索の高速化)

**制約**:
- name は空文字列不可
- path は '/' で開始・終了する形式（例: /1/5/12/）

### 2. Projects (Jiraプロジェクト)

Jiraから取得したプロジェクト情報を管理。

| カラム名 | データ型 | NULL | デフォルト値 | 説明 |
|---------|---------|------|------------|------|
| id | BIGSERIAL | NOT NULL | - | 主キー |
| jira_project_id | VARCHAR(100) | NOT NULL | - | JiraのプロジェクトID |
| key | VARCHAR(50) | NOT NULL | - | プロジェクトキー（例: PROJ） |
| name | VARCHAR(255) | NOT NULL | - | プロジェクト名 |
| lead_account_id | VARCHAR(100) | NULL | NULL | プロジェクトリーダーのJira Account ID |
| lead_email | VARCHAR(255) | NULL | NULL | リーダーのメールアドレス（組織推測用） |
| organization_id | BIGINT | NULL | NULL | 所属組織ID（未分類の場合はNULL） |
| created_at | TIMESTAMP | NOT NULL | CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | NOT NULL | CURRENT_TIMESTAMP | 更新日時 |

**主キー**: id
**外部キー**: organization_id → Organizations(id)
**ユニーク制約**: jira_project_id
**インデックス**:
- `idx_projects_jira_project_id` on jira_project_id (UNIQUE)
- `idx_projects_organization_id` on organization_id
- `idx_projects_key` on key

### 3. Issues (チケット情報)

Jiraのチケット（Issue）情報と遅延ステータスを管理。

| カラム名 | データ型 | NULL | デフォルト値 | 説明 |
|---------|---------|------|------------|------|
| id | BIGSERIAL | NOT NULL | - | 主キー |
| jira_issue_id | VARCHAR(100) | NOT NULL | - | JiraのIssue ID |
| jira_issue_key | VARCHAR(100) | NOT NULL | - | Issue Key（例: PROJ-123） |
| project_id | BIGINT | NOT NULL | - | プロジェクトID |
| summary | TEXT | NOT NULL | - | チケット概要 |
| status | VARCHAR(100) | NOT NULL | - | ステータス名 |
| status_category | VARCHAR(50) | NOT NULL | - | ステータスカテゴリ（To Do/In Progress/Done） |
| due_date | DATE | NULL | NULL | 納期 |
| assignee_name | VARCHAR(255) | NULL | NULL | 担当者名 |
| assignee_account_id | VARCHAR(100) | NULL | NULL | 担当者のJira Account ID |
| delay_status | VARCHAR(20) | NOT NULL | 'GREEN' | 遅延ステータス（RED/YELLOW/GREEN） |
| priority | VARCHAR(50) | NULL | NULL | 優先度 |
| issue_type | VARCHAR(100) | NULL | NULL | Issue Type（Bug/Task/Story等） |
| last_updated_at | TIMESTAMP | NOT NULL | - | Jira上の最終更新日時 |
| created_at | TIMESTAMP | NOT NULL | CURRENT_TIMESTAMP | レコード作成日時 |
| updated_at | TIMESTAMP | NOT NULL | CURRENT_TIMESTAMP | レコード更新日時 |

**主キー**: id
**外部キー**: project_id → Projects(id)
**ユニーク制約**: jira_issue_id
**インデックス**:
- `idx_issues_jira_issue_id` on jira_issue_id (UNIQUE)
- `idx_issues_jira_issue_key` on jira_issue_key
- `idx_issues_project_id` on project_id
- `idx_issues_delay_status` on delay_status
- `idx_issues_status_category` on status_category
- `idx_issues_due_date` on due_date
- `idx_issues_updated_at` on last_updated_at (Delta Sync用)

**制約**:
- delay_status は 'RED', 'YELLOW', 'GREEN' のいずれか
- status_category は 'To Do', 'In Progress', 'Done' のいずれか

### 4. SyncLogs (バッチ実行ログ)

Jira同期バッチの実行履歴を記録。

| カラム名 | データ型 | NULL | デフォルト値 | 説明 |
|---------|---------|------|------------|------|
| id | BIGSERIAL | NOT NULL | - | 主キー |
| sync_type | VARCHAR(20) | NOT NULL | - | 同期タイプ（FULL/DELTA） |
| executed_at | TIMESTAMP | NOT NULL | CURRENT_TIMESTAMP | 実行開始日時 |
| completed_at | TIMESTAMP | NULL | NULL | 実行完了日時 |
| status | VARCHAR(20) | NOT NULL | 'RUNNING' | ステータス（RUNNING/SUCCESS/FAILURE） |
| projects_synced | INTEGER | NULL | 0 | 同期したプロジェクト数 |
| issues_synced | INTEGER | NULL | 0 | 同期したIssue数 |
| error_message | TEXT | NULL | NULL | エラーメッセージ |
| duration_seconds | INTEGER | NULL | NULL | 実行時間（秒） |

**主キー**: id
**インデックス**:
- `idx_sync_logs_executed_at` on executed_at DESC
- `idx_sync_logs_status` on status

**制約**:
- sync_type は 'FULL' または 'DELTA'
- status は 'RUNNING', 'SUCCESS', 'FAILURE' のいずれか

## インデックス戦略

### パフォーマンス最適化のためのインデックス

1. **組織階層検索用**
   - `Organizations.path`: 階層クエリの高速化（LIKE '/1/5/%'）

2. **ダッシュボード集計用**
   - `Issues.project_id + delay_status`: プロジェクトごとの遅延チケット集計
   - `Issues.status_category`: ステータスカテゴリによる絞り込み

3. **Delta Sync用**
   - `Issues.last_updated_at`: 更新日時による差分取得

4. **外部キー参照**
   - すべての外部キーにインデックスを作成

## 遅延ステータス判定ロジック

Issue挿入・更新時にトリガーまたはアプリケーション側で以下のロジックを適用：

```sql
CASE
  WHEN status_category != 'Done' AND due_date < CURRENT_DATE THEN 'RED'
  WHEN status_category != 'Done' AND due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '3 days' THEN 'YELLOW'
  WHEN status_category != 'Done' AND due_date IS NULL THEN 'YELLOW'
  ELSE 'GREEN'
END
```

## 拡張性の考慮

### 将来の機能追加に備えた設計

1. **Issues テーブル**
   - `original_estimate`: 見積もり工数
   - `time_spent`: 実績工数
   - `remaining_estimate`: 残工数
   - これらのカラムは将来の工数予実機能のために追加可能

2. **Organizations テーブル**
   - `manager_name`: 組織責任者名
   - `cost_center`: コストセンター
   - 将来の管理機能拡張に対応

3. **Projects テーブル**
   - `start_date`: プロジェクト開始日
   - `end_date`: プロジェクト終了予定日
   - プロジェクト管理機能の拡張に対応

## データ整合性

### 外部キー制約

- `Organizations.parent_id` → `Organizations.id` (ON DELETE RESTRICT)
- `Projects.organization_id` → `Organizations.id` (ON DELETE SET NULL)
- `Issues.project_id` → `Projects.id` (ON DELETE CASCADE)

### 削除ポリシー

- 組織に子組織がある場合は削除不可
- 組織を削除する際、紐付けられたプロジェクトは未分類状態（organization_id = NULL）になる
- プロジェクトを削除すると、紐付くIssueも削除される

## パフォーマンス見積もり

想定データ量（1年運用時）:
- Organizations: ~100レコード
- Projects: ~500レコード
- Issues: ~50,000レコード
- SyncLogs: ~8,760レコード（1時間ごと）

上記規模であれば、適切なインデックスにより以下を達成可能:
- ダッシュボードサマリ取得: 500ms以内
- プロジェクト一覧取得: 100ms以内
- チケット詳細取得: 50ms以内

## 次のステップ

1. DB-002: マイグレーションツールのセットアップ
2. DB-003: 初期スキーマの実装（DDL作成）
3. DB-004: シードデータの作成
