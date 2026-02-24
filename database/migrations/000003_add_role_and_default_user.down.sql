-- Remove default admin user
DELETE FROM users WHERE email = 'admin@example.com';

-- Revert role CHECK constraint to original (admin, viewer only)
ALTER TABLE users DROP CONSTRAINT users_role_check;
ALTER TABLE users ADD CONSTRAINT users_role_check
  CHECK (role IN ('admin', 'viewer'));
