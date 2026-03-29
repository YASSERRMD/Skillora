-- Script to grant admin privileges to a user
-- Usage: psql -U skillora -d skillora -f backend/scripts/grant_admin.sql

-- Replace 'admin@example.com' with the email address of the user you want to make admin
UPDATE users
SET is_admin = true,
    updated_at = NOW()
WHERE email = 'admin@example.com';

-- Verify the change
SELECT id, email, full_name, is_admin, updated_at
FROM users
WHERE is_admin = true;
