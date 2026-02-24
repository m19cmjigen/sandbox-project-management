-- Jira connection settings table
CREATE TABLE jira_settings (
    id          BIGSERIAL PRIMARY KEY,
    jira_url    VARCHAR(500)  NOT NULL,
    email       VARCHAR(255)  NOT NULL,
    api_token   VARCHAR(500)  NOT NULL,
    created_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_jira_settings_updated_at
    BEFORE UPDATE ON jira_settings
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE  jira_settings             IS 'Jira API接続設定';
COMMENT ON COLUMN jira_settings.api_token   IS 'Jira APIトークン（平文）';

-- Add is_active flag to projects so admins can hide specific projects
ALTER TABLE projects ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT true;

COMMENT ON COLUMN projects.is_active IS '表示フラグ（false=一覧から除外）';
