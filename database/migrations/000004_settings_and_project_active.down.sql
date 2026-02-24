ALTER TABLE projects DROP COLUMN IF EXISTS is_active;

DROP TRIGGER IF EXISTS update_jira_settings_updated_at ON jira_settings;
DROP TABLE IF EXISTS jira_settings;
