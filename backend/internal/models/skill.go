package models

import "time"

// Category classifies a logical grouping of skills.
type Category struct {
	ID        string    `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Slug      string    `db:"slug"       json:"slug"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Skill represents a specific barterable competency.
type Skill struct {
	ID          string    `db:"id"          json:"id"`
	CategoryID  string    `db:"category_id" json:"category_id"`
	Name        string    `db:"name"        json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
}
