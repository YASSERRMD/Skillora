-- Migration: 002_create_llm_providers
-- Creates the llm_providers table for dynamic AI Engine routing.

CREATE TABLE IF NOT EXISTS llm_providers (
    id                UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    provider_name     VARCHAR(50) NOT NULL, -- e.g., 'openai', 'anthropic', 'deepseek'
    model_name        VARCHAR(50) NOT NULL, -- e.g., 'gpt-4o', 'claude-3-5-sonnet'
    encrypted_api_key BYTEA       NOT NULL, -- AES-256-GCM encrypted key
    use_case          VARCHAR(50) NOT NULL, -- e.g., 'general', 'embedding', 'course_generation'
    priority          INT         NOT NULL DEFAULT 1, -- 1=Primary, 2=Fallback
    is_active         BOOLEAN     NOT NULL DEFAULT true,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Ensure priority ranks are unique per use case to prevent routing ambiguity
    UNIQUE (use_case, priority)
);

CREATE INDEX IF NOT EXISTS idx_llm_providers_use_case_active ON llm_providers (use_case) WHERE is_active = true;

-- Trigger to auto-update updated_at timestamp
DROP TRIGGER IF EXISTS trg_llm_providers_updated_at ON llm_providers;
CREATE TRIGGER trg_llm_providers_updated_at
    BEFORE UPDATE ON llm_providers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
