package models

import "time"

// AdminCredential represents local authentication credentials for admin users
type AdminCredential struct {
	ID           string     `db:"id"            json:"id"`
	Username     string     `db:"username"      json:"username"`
	PasswordHash string     `db:"password_hash" json:"-"`
	UserID       string     `db:"user_id"       json:"user_id"`
	IsActive     bool       `db:"is_active"     json:"is_active"`
	CreatedAt    time.Time  `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"    json:"updated_at"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at"`
}
