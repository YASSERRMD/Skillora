-- Migration: 007_create_notifications
-- Push-style notification system for barter events and system messages.

CREATE TABLE IF NOT EXISTS notifications (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title       VARCHAR(100) NOT NULL,
    body        TEXT        NOT NULL,
    type        VARCHAR(20)  NOT NULL CHECK (type IN ('barter_proposed', 'barter_accepted', 'barter_completed', 'barter_cancelled', 'system')),
    is_read     BOOLEAN     NOT NULL DEFAULT false,
    metadata    JSONB       DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id     ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_unread      ON notifications(user_id) WHERE is_read = false;
