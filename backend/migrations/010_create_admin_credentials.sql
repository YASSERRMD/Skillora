-- Migration: Create admin_credentials table for local admin authentication
-- This allows admins to login with username/password in addition to OAuth

CREATE TABLE IF NOT EXISTS admin_credentials (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    username        TEXT        NOT NULL UNIQUE,
    password_hash   TEXT        NOT NULL,
    user_id         UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_active       BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_login_at   TIMESTAMPTZ
);

-- Create index for username lookups
CREATE INDEX idx_admin_credentials_username ON admin_credentials(username);
CREATE INDEX idx_admin_credentials_user_id ON admin_credentials(user_id);

-- Add comments for documentation
COMMENT ON TABLE admin_credentials IS 'Local authentication credentials for admin users. Allows username/password login alongside OAuth.';
COMMENT ON COLUMN admin_credentials.username IS 'Unique username for local admin authentication';
COMMENT ON COLUMN admin_credentials.password_hash IS 'Bcrypt hash of the admin password';
COMMENT ON COLUMN admin_credentials.user_id IS 'Reference to the users table. Must have is_admin=true';
COMMENT ON COLUMN admin_credentials.is_active IS 'Whether this credential is active. Can be disabled without deleting.';
COMMENT ON COLUMN admin_credentials.last_login_at IS 'Timestamp of the last successful login using these credentials';
