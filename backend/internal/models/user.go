package models

import "time"

// User represents a Skillora user authenticated via Google OAuth2.
type User struct {
	ID        string    `db:"id"         json:"id"`
	GoogleID  string    `db:"google_id"  json:"-"`
	Email     string    `db:"email"      json:"email"`
	FullName  string    `db:"full_name"  json:"full_name"`
	AvatarURL *string   `db:"avatar_url" json:"avatar_url"`
	Bio       string     `db:"bio"        json:"bio"`
	IsAdmin   bool       `db:"is_admin"   json:"is_admin"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}
