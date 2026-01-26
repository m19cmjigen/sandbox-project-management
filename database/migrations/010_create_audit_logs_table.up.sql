-- Create audit_logs table for tracking user actions
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    username VARCHAR(50),
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id VARCHAR(100),
    method VARCHAR(10) NOT NULL,
    path VARCHAR(500) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    request_body TEXT,
    response_status INT,
    response_body TEXT,
    error_message TEXT,
    duration_ms INT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for efficient querying
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_username ON audit_logs(username);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX idx_audit_logs_resource_id ON audit_logs(resource_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX idx_audit_logs_method ON audit_logs(method);

COMMENT ON TABLE audit_logs IS 'Audit trail of all user actions and API calls';
COMMENT ON COLUMN audit_logs.user_id IS 'ID of the user who performed the action (NULL for unauthenticated requests)';
COMMENT ON COLUMN audit_logs.action IS 'Action performed (e.g., CREATE, UPDATE, DELETE, VIEW, LOGIN)';
COMMENT ON COLUMN audit_logs.resource_type IS 'Type of resource (e.g., user, organization, project, issue)';
COMMENT ON COLUMN audit_logs.resource_id IS 'ID of the affected resource';
COMMENT ON COLUMN audit_logs.method IS 'HTTP method (GET, POST, PUT, DELETE, etc.)';
COMMENT ON COLUMN audit_logs.path IS 'Request path';
COMMENT ON COLUMN audit_logs.duration_ms IS 'Request duration in milliseconds';
