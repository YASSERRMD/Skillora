package repository

import (
	"context"
	"fmt"

	pgvector "github.com/pgvector/pgvector-go"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skillora/backend/internal/models"
)

// VectorRepository handles embedding storage and ANN similarity search.
type VectorRepository struct {
	db *pgxpool.Pool
}

// NewVectorRepository constructs the vector repo.
func NewVectorRepository(db *pgxpool.Pool) *VectorRepository {
	return &VectorRepository{db: db}
}

// UpsertSkillEmbedding stores or updates the embedding vector for a skill.
func (r *VectorRepository) UpsertSkillEmbedding(ctx context.Context, skillID string, embedding []float32) error {
	vec := pgvector.NewVector(embedding)
	const q = `UPDATE skills SET embedding = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, q, vec, skillID)
	if err != nil {
		return fmt.Errorf("UpsertSkillEmbedding: %w", err)
	}
	return nil
}

// FindSimilarSkills returns the top-k skills nearest to the given embedding vector,
// excluding skills already offered by the querying user.
func (r *VectorRepository) FindSimilarSkills(ctx context.Context, queryVec []float32, excludeUserID string, topK int) ([]models.SkillMatch, error) {
	vec := pgvector.NewVector(queryVec)
	const q = `
		SELECT
			s.id,
			s.name,
			c.name AS category_name,
			us.user_id AS owner_id,
			us.proficiency_level,
			us.credit_value,
			1 - (s.embedding <=> $1) AS similarity_score
		FROM skills s
		JOIN categories c ON s.category_id = c.id
		JOIN user_skills us ON s.id = us.skill_id
		WHERE s.embedding IS NOT NULL
		  AND us.is_verified = true
		  AND us.user_id != $2
		ORDER BY s.embedding <=> $1
		LIMIT $3
	`
	rows, err := r.db.Query(ctx, q, vec, excludeUserID, topK)
	if err != nil {
		return nil, fmt.Errorf("FindSimilarSkills: %w", err)
	}
	defer rows.Close()

	var list []models.SkillMatch
	for rows.Next() {
		var m models.SkillMatch
		if err := rows.Scan(
			&m.SkillID, &m.SkillName, &m.CategoryName,
			&m.OwnerID, &m.ProficiencyLevel, &m.CreditValue, &m.SimilarityScore,
		); err != nil {
			return nil, fmt.Errorf("FindSimilarSkills scan: %w", err)
		}
		list = append(list, m)
	}
	if list == nil {
		list = make([]models.SkillMatch, 0)
	}
	return list, rows.Err()
}
