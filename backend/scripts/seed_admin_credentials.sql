-- Script to seed default admin credentials for local development
-- This creates the default admin username/password for first-time setup

-- Default credentials:
-- Username: admin
-- Password: admin123

-- First, ensure the dev user has admin privileges
UPDATE users
SET is_admin = true
WHERE email = 'dev@skillora.local';

-- Insert admin credentials for the dev user
-- Password hash for "admin123" (bcrypt, cost 12)
INSERT INTO admin_credentials (username, password_hash, user_id)
VALUES (
    'admin',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzpLaEmc3i',
    (SELECT id FROM users WHERE email = 'dev@skillora.local')
)
ON CONFLICT (username) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    updated_at = NOW();

-- Verify the setup
SELECT
    u.id,
    u.email,
    u.is_admin,
    ac.username,
    ac.created_at,
    ac.is_active
FROM users u
LEFT JOIN admin_credentials ac ON u.id = ac.user_id
WHERE u.email = 'dev@skillora.local';

-- Expected output:
-- - email: dev@skillora.local
-- - is_admin: true
-- - username: admin
-- - is_active: true
