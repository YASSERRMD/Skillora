-- Migration: Add admin role support to users table
-- This migration adds role-based access control for admin functionality

-- Add is_admin flag to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_admin BOOLEAN DEFAULT FALSE;

-- Create index for efficient admin lookups
CREATE INDEX IF NOT EXISTS idx_users_is_admin ON users(is_admin);

-- Add comment for documentation
COMMENT ON COLUMN users.is_admin IS 'Administrative access flag for role-based access control. Only true for users with admin privileges.';
