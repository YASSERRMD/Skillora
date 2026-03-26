-- Migration: 004_create_user_skills
-- Maps users to their verified capabilities.

CREATE TABLE IF NOT EXISTS user_skills (
    user_id           UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_id          UUID        NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    proficiency_level INT         NOT NULL CHECK (proficiency_level >= 1 AND proficiency_level <= 5),
    credit_value      INT         NOT NULL CHECK (credit_value > 0),
    is_verified       BOOLEAN     NOT NULL DEFAULT false, -- True if AI appraised successfully
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- A user can only have a specific skill mapped once
    PRIMARY KEY (user_id, skill_id)
);

CREATE INDEX IF NOT EXISTS idx_user_skills_user_id ON user_skills (user_id);
CREATE INDEX IF NOT EXISTS idx_user_skills_skill_id ON user_skills (skill_id);
