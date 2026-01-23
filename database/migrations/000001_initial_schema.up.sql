-- 全社プロジェクト進捗可視化プラットフォーム
-- 初期スキーママイグレーション (UP)

-- ==============================================
-- 1. Organizations (組織マスタ)
-- ==============================================

CREATE TABLE organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL CHECK (name != ''),
    parent_id BIGINT REFERENCES organizations(id) ON DELETE RESTRICT,
    path VARCHAR(1000) NOT NULL,
    level INTEGER NOT NULL DEFAULT 0 CHECK (level >= 0 AND level <= 2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_organizations_parent_id ON organizations(parent_id);
CREATE INDEX idx_organizations_path ON organizations(path);
CREATE INDEX idx_organizations_level ON organizations(level);

COMMENT ON TABLE organizations IS '組織階層マスタ（本部-部-課）';
COMMENT ON COLUMN organizations.id IS '主キー';
COMMENT ON COLUMN organizations.name IS '組織名';
COMMENT ON COLUMN organizations.parent_id IS '親組織ID（最上位はNULL）';
COMMENT ON COLUMN organizations.path IS '階層パス（例: /1/5/12/）';
COMMENT ON COLUMN organizations.level IS '階層レベル（0:本部, 1:部, 2:課）';

-- ==============================================
-- 2. Projects (Jiraプロジェクト)
-- ==============================================

CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    jira_project_id VARCHAR(100) NOT NULL UNIQUE,
    key VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    lead_account_id VARCHAR(100),
    lead_email VARCHAR(255),
    organization_id BIGINT REFERENCES organizations(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_projects_jira_project_id ON projects(jira_project_id);
CREATE INDEX idx_projects_organization_id ON projects(organization_id);
CREATE INDEX idx_projects_key ON projects(key);

COMMENT ON TABLE projects IS 'Jiraプロジェクト情報';
COMMENT ON COLUMN projects.jira_project_id IS 'JiraのプロジェクトID';
COMMENT ON COLUMN projects.key IS 'プロジェクトキー（例: PROJ）';
COMMENT ON COLUMN projects.organization_id IS '所属組織ID（未分類の場合はNULL）';

-- ==============================================
-- 3. Issues (チケット情報)
-- ==============================================

CREATE TABLE issues (
    id BIGSERIAL PRIMARY KEY,
    jira_issue_id VARCHAR(100) NOT NULL UNIQUE,
    jira_issue_key VARCHAR(100) NOT NULL,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    summary TEXT NOT NULL,
    status VARCHAR(100) NOT NULL,
    status_category VARCHAR(50) NOT NULL CHECK (status_category IN ('To Do', 'In Progress', 'Done')),
    due_date DATE,
    assignee_name VARCHAR(255),
    assignee_account_id VARCHAR(100),
    delay_status VARCHAR(20) NOT NULL DEFAULT 'GREEN' CHECK (delay_status IN ('RED', 'YELLOW', 'GREEN')),
    priority VARCHAR(50),
    issue_type VARCHAR(100),
    last_updated_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_issues_jira_issue_id ON issues(jira_issue_id);
CREATE INDEX idx_issues_jira_issue_key ON issues(jira_issue_key);
CREATE INDEX idx_issues_project_id ON issues(project_id);
CREATE INDEX idx_issues_delay_status ON issues(delay_status);
CREATE INDEX idx_issues_status_category ON issues(status_category);
CREATE INDEX idx_issues_due_date ON issues(due_date);
CREATE INDEX idx_issues_updated_at ON issues(last_updated_at);
CREATE INDEX idx_issues_project_delay ON issues(project_id, delay_status);

COMMENT ON TABLE issues IS 'Jiraチケット（Issue）情報と遅延ステータス';
COMMENT ON COLUMN issues.delay_status IS '遅延ステータス（RED/YELLOW/GREEN）';

-- ==============================================
-- 4. SyncLogs (バッチ実行ログ)
-- ==============================================

CREATE TABLE sync_logs (
    id BIGSERIAL PRIMARY KEY,
    sync_type VARCHAR(20) NOT NULL CHECK (sync_type IN ('FULL', 'DELTA')),
    executed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'RUNNING' CHECK (status IN ('RUNNING', 'SUCCESS', 'FAILURE')),
    projects_synced INTEGER DEFAULT 0,
    issues_synced INTEGER DEFAULT 0,
    error_message TEXT,
    duration_seconds INTEGER
);

CREATE INDEX idx_sync_logs_executed_at ON sync_logs(executed_at DESC);
CREATE INDEX idx_sync_logs_status ON sync_logs(status);
CREATE INDEX idx_sync_logs_sync_type ON sync_logs(sync_type);

COMMENT ON TABLE sync_logs IS 'Jira同期バッチの実行履歴';

-- ==============================================
-- 5. トリガー関数
-- ==============================================

-- updated_at 自動更新関数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 遅延ステータス自動計算関数
CREATE OR REPLACE FUNCTION calculate_delay_status()
RETURNS TRIGGER AS $$
BEGIN
    NEW.delay_status := CASE
        WHEN NEW.status_category != 'Done' AND NEW.due_date < CURRENT_DATE THEN 'RED'
        WHEN NEW.status_category != 'Done' AND NEW.due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '3 days' THEN 'YELLOW'
        WHEN NEW.status_category != 'Done' AND NEW.due_date IS NULL THEN 'YELLOW'
        ELSE 'GREEN'
    END;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- ==============================================
-- 6. トリガー設定
-- ==============================================

-- Organizations
CREATE TRIGGER update_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Projects
CREATE TRIGGER update_projects_updated_at
    BEFORE UPDATE ON projects
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Issues
CREATE TRIGGER update_issues_updated_at
    BEFORE UPDATE ON issues
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER calculate_issue_delay_status
    BEFORE INSERT OR UPDATE ON issues
    FOR EACH ROW
    EXECUTE FUNCTION calculate_delay_status();

-- ==============================================
-- 7. ビュー
-- ==============================================

-- プロジェクト遅延サマリ
CREATE OR REPLACE VIEW project_delay_summary AS
SELECT
    p.id AS project_id,
    p.jira_project_id,
    p.key AS project_key,
    p.name AS project_name,
    p.organization_id,
    COUNT(i.id) AS total_issues,
    COUNT(CASE WHEN i.delay_status = 'RED' THEN 1 END) AS red_issues,
    COUNT(CASE WHEN i.delay_status = 'YELLOW' THEN 1 END) AS yellow_issues,
    COUNT(CASE WHEN i.delay_status = 'GREEN' THEN 1 END) AS green_issues,
    COUNT(CASE WHEN i.status_category != 'Done' THEN 1 END) AS open_issues,
    COUNT(CASE WHEN i.status_category = 'Done' THEN 1 END) AS done_issues
FROM
    projects p
    LEFT JOIN issues i ON p.id = i.project_id
GROUP BY
    p.id, p.jira_project_id, p.key, p.name, p.organization_id;

-- 組織遅延サマリ
CREATE OR REPLACE VIEW organization_delay_summary AS
SELECT
    o.id AS organization_id,
    o.name AS organization_name,
    o.path,
    o.level,
    COUNT(DISTINCT p.id) AS total_projects,
    COUNT(DISTINCT CASE WHEN pds.red_issues > 0 THEN p.id END) AS delayed_projects,
    SUM(pds.red_issues) AS total_red_issues,
    SUM(pds.yellow_issues) AS total_yellow_issues,
    SUM(pds.green_issues) AS total_green_issues,
    SUM(pds.open_issues) AS total_open_issues
FROM
    organizations o
    LEFT JOIN projects p ON o.id = p.organization_id
    LEFT JOIN project_delay_summary pds ON p.id = pds.project_id
GROUP BY
    o.id, o.name, o.path, o.level;
