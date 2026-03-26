-- Migration: 003_create_skills
-- Creates categories and skills reference tables for the barter economy platform.

CREATE TABLE IF NOT EXISTS categories (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL UNIQUE,
    slug       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index the categories slug for fast lookup and URL routing.
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories (slug);

CREATE TABLE IF NOT EXISTS skills (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID        NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL UNIQUE,
    description TEXT        DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index skills for querying directly by category.
CREATE INDEX IF NOT EXISTS idx_skills_category_id ON skills (category_id);
-- Text search index for skill mapping via LLM pipeline.
CREATE INDEX IF NOT EXISTS idx_skills_name ON skills (name);
