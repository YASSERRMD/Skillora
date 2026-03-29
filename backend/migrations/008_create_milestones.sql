-- Migration: 008_create_milestones
-- Adds milestones support for staged releases of credits in a barter transaction.

CREATE TABLE IF NOT EXISTS milestones (
    id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    barter_id      UUID        NOT NULL REFERENCES barter_transactions(id) ON DELETE CASCADE,
    title          VARCHAR(100) NOT NULL,
    description    TEXT        NOT NULL,
    credit_portion INT         NOT NULL CHECK (credit_portion > 0),
    status         VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'approved')),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_milestones_barter_id ON milestones(barter_id);
CREATE INDEX IF NOT EXISTS idx_milestones_status    ON milestones(status);

-- Add updated_at trigger for milestones.
CREATE TRIGGER set_milestones_updated_at
    BEFORE UPDATE ON milestones
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
