-- Extend role CHECK constraint to include project_manager
ALTER TABLE users DROP CONSTRAINT users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
  CHECK (role IN ('admin', 'project_manager', 'viewer'));

-- Insert default admin user for initial setup (password: Admin1234!)
INSERT INTO users (email, password_hash, role)
VALUES ('admin@example.com', '$2a$12$oz0TRkSXt3BKD3XfSK9An.bzeoikiJbWJaRfskyFC0lAKWas3H4hi', 'admin')
ON CONFLICT (email) DO NOTHING;
