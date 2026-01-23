-- 全社プロジェクト進捗可視化プラットフォーム
-- 初期スキーママイグレーション (DOWN)

-- ビューの削除
DROP VIEW IF EXISTS organization_delay_summary;
DROP VIEW IF EXISTS project_delay_summary;

-- トリガーの削除
DROP TRIGGER IF EXISTS calculate_issue_delay_status ON issues;
DROP TRIGGER IF EXISTS update_issues_updated_at ON issues;
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;

-- トリガー関数の削除
DROP FUNCTION IF EXISTS calculate_delay_status();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- テーブルの削除（外部キー制約を考慮した順序）
DROP TABLE IF EXISTS sync_logs;
DROP TABLE IF EXISTS issues;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS organizations;
