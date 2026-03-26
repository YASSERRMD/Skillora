package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skillora/backend/internal/models"
)

// SkillRepository manages categories and taxonomy logic.
type SkillRepository struct {
	db *pgxpool.Pool
}

// NewSkillRepository constructs the skills repository.
func NewSkillRepository(db *pgxpool.Pool) *SkillRepository {
	return &SkillRepository{db: db}
}

// GetCategories returns all parent-level taxonomy groupings sorted alphabetically.
func (r *SkillRepository) GetCategories(ctx context.Context) ([]models.Category, error) {
	const q = `SELECT id, name, slug, created_at FROM categories ORDER BY name ASC`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("GetCategories: %w", err)
	}
	defer rows.Close()

	var list []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetCategories scan: %w", err)
		}
		list = append(list, c)
	}

	if list == nil {
		list = make([]models.Category, 0)
	}
	return list, rows.Err()
}

// GetCategorySkills returns all skills belonging to a given category ID sorted alphabetically.
func (r *SkillRepository) GetCategorySkills(ctx context.Context, categoryID string) ([]models.Skill, error) {
	const q = `
		SELECT id, category_id, name, description, created_at
		FROM skills
		WHERE category_id = $1
		ORDER BY name ASC
	`
	rows, err := r.db.Query(ctx, q, categoryID)
	if err != nil {
		return nil, fmt.Errorf("GetCategorySkills for %s: %w", categoryID, err)
	}
	defer rows.Close()

	var list []models.Skill
	for rows.Next() {
		var s models.Skill
		if err := rows.Scan(&s.ID, &s.CategoryID, &s.Name, &s.Description, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetCategorySkills scan: %w", err)
		}
		list = append(list, s)
	}

	if list == nil {
		list = make([]models.Skill, 0)
	}
	return list, rows.Err()
}
