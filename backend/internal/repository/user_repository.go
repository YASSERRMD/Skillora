// Package repository provides data access functions for all domain entities.
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skillora/backend/internal/models"
)

// UserRepository handles all DB operations on the users table.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository constructs a UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// UpsertGoogleUser creates a new user from Google OAuth data or updates existing fields.
// Returns the user and a boolean indicating whether the user was newly created.
func (r *UserRepository) UpsertGoogleUser(
	ctx context.Context,
	googleID, email, fullName, avatarURL string,
) (*models.User, bool, error) {
	const q = `
		INSERT INTO users (google_id, email, full_name, avatar_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (google_id) DO UPDATE
			SET email      = EXCLUDED.email,
				full_name  = EXCLUDED.full_name,
				avatar_url = EXCLUDED.avatar_url,
				updated_at = NOW()
		RETURNING id, google_id, email, full_name, avatar_url, bio, is_admin, created_at, updated_at,
		          (xmax = 0) AS is_new
	`

	row := r.db.QueryRow(ctx, q, googleID, email, fullName, avatarURL)

	var u models.User
	var isNew bool
	err := row.Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.FullName, &u.AvatarURL,
		&u.Bio, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt, &isNew,
	)
	if err != nil {
		return nil, false, fmt.Errorf("upsert user: %w", err)
	}
	return &u, isNew, nil
}

// GetByID fetches a user by their primary key UUID.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	const q = `
		SELECT id, google_id, email, full_name, avatar_url, bio, is_admin, created_at, updated_at
		FROM users WHERE id = $1
	`
	row := r.db.QueryRow(ctx, q, id)
	var u models.User
	if err := row.Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.FullName, &u.AvatarURL,
		&u.Bio, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &u, nil
}

// UpdateBio patches the bio field for a user.
func (r *UserRepository) UpdateBio(ctx context.Context, id, bio string) error {
	const q = `UPDATE users SET bio = $1 WHERE id = $2`
	if _, err := r.db.Exec(ctx, q, bio, id); err != nil {
		return fmt.Errorf("update user bio: %w", err)
	}
	return nil
}
