package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skillora/backend/internal/models"
)

// UserSkillRepository maps users to their respective competencies.
type UserSkillRepository struct {
	db *pgxpool.Pool
}

// NewUserSkillRepository constructs the dependency.
func NewUserSkillRepository(db *pgxpool.Pool) *UserSkillRepository {
	return &UserSkillRepository{db: db}
}

// AddUserSkill stores a new mapping, typically after an LLM has appraised it.
func (r *UserSkillRepository) AddUserSkill(ctx context.Context, us models.UserSkill) error {
	const q = `
		INSERT INTO user_skills (user_id, skill_id, proficiency_level, credit_value, is_verified)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, skill_id) DO UPDATE 
			SET proficiency_level = EXCLUDED.proficiency_level,
			    credit_value      = EXCLUDED.credit_value,
			    is_verified       = EXCLUDED.is_verified
	`
	_, err := r.db.Exec(ctx, q, us.UserID, us.SkillID, us.ProficiencyLevel, us.CreditValue, us.IsVerified)
	if err != nil {
		return fmt.Errorf("AddUserSkill: %w", err)
	}
	return nil
}

// GetUserSkills returns all verified taxonomy details for a specific user ID.
func (r *UserSkillRepository) GetUserSkills(ctx context.Context, userID string) ([]models.UserSkillDetail, error) {
	const q = `
		SELECT 
			us.user_id, us.skill_id, us.proficiency_level, us.credit_value, us.is_verified, us.created_at,
			s.name as skill_name,
			c.name as category_name
		FROM user_skills us
		JOIN skills s ON us.skill_id = s.id
		JOIN categories c ON s.category_id = c.id
		WHERE us.user_id = $1
		ORDER BY s.name ASC
	`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserSkills: %w", err)
	}
	defer rows.Close()

	var list []models.UserSkillDetail
	for rows.Next() {
		var u models.UserSkillDetail
		if err := rows.Scan(
			&u.UserID, &u.SkillID, &u.ProficiencyLevel, &u.CreditValue, &u.IsVerified, &u.CreatedAt,
			&u.SkillName, &u.CategoryName,
		); err != nil {
			return nil, fmt.Errorf("GetUserSkills scan: %w", err)
		}
		list = append(list, u)
	}
	
	if list == nil {
		list = make([]models.UserSkillDetail, 0)
	}
	return list, rows.Err()
}
