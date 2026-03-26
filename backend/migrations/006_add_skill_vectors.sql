-- Migration: 006_add_skill_vectors
-- Adds pgvector embedding column to skills for nearest-neighbour matching.

-- Enable the pgvector extension (must be installed on the server).
CREATE EXTENSION IF NOT EXISTS vector;

-- Add a 1536-dimension embedding column matching text-embedding-3-small output.
ALTER TABLE skills
    ADD COLUMN IF NOT EXISTS embedding vector(1536);

-- HNSW index for approximate nearest-neighbour cosine similarity searches.
-- This provides O(log n) ANN queries which scales to millions of skills.
CREATE INDEX IF NOT EXISTS idx_skills_embedding_hnsw
    ON skills
    USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);
