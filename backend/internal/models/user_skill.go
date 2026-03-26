package models

import "time"

// UserSkill maps a user to a tested or proposed skill on their profile.
type UserSkill struct {
	UserID           string    `db:"user_id"           json:"user_id"`
	SkillID          string    `db:"skill_id"          json:"skill_id"`
	ProficiencyLevel int       `db:"proficiency_level" json:"proficiency_level"` // 1-5 scale
	CreditValue      int       `db:"credit_value"      json:"credit_value"`
	IsVerified       bool      `db:"is_verified"       json:"is_verified"`       // Approved by AI or Oracle
	CreatedAt        time.Time `db:"created_at"        json:"created_at"`
}

// UserSkillDetail includes the taxonomy details inside the user skill association.
type UserSkillDetail struct {
	UserSkill
	SkillName    string `json:"skill_name"`
	CategoryName string `json:"category_name"`
}
