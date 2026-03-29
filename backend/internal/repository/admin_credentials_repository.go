package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skillora/backend/internal/models"
)

// AdminCredentialsRepository handles DB operations for admin credentials
type AdminCredentialsRepository struct {
	db *pgxpool.Pool
}

// NewAdminCredentialsRepository creates a new admin credentials repository
func NewAdminCredentialsRepository(db *pgxpool.Pool) *AdminCredentialsRepository {
	return &AdminCredentialsRepository{db: db}
}

// Create creates a new admin credential
func (r *AdminCredentialsRepository) Create(
	ctx context.Context,
	username, passwordHash, userID string,
) (*models.AdminCredential, error) {
	const q = `
		INSERT INTO admin_credentials (username, password_hash, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, username, password_hash, user_id, is_active, created_at, updated_at, last_login_at
	`

	row := r.db.QueryRow(ctx, q, username, passwordHash, userID)

	var cred models.AdminCredential
	err := row.Scan(
		&cred.ID, &cred.Username, &cred.PasswordHash, &cred.UserID,
		&cred.IsActive, &cred.CreatedAt, &cred.UpdatedAt, &cred.LastLoginAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create admin credential: %w", err)
	}

	return &cred, nil
}

// GetByUsername fetches admin credentials by username
func (r *AdminCredentialsRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*models.AdminCredential, error) {
	const q = `
		SELECT id, username, password_hash, user_id, is_active, created_at, updated_at, last_login_at
		FROM admin_credentials
		WHERE username = $1
	`

	row := r.db.QueryRow(ctx, q, username)

	var cred models.AdminCredential
	err := row.Scan(
		&cred.ID, &cred.Username, &cred.PasswordHash, &cred.UserID,
		&cred.IsActive, &cred.CreatedAt, &cred.UpdatedAt, &cred.LastLoginAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get admin credential by username: %w", err)
	}

	return &cred, nil
}

// UpdatePassword updates the password hash for an admin credential
func (r *AdminCredentialsRepository) UpdatePassword(
	ctx context.Context,
	credentialID, newPasswordHash string,
) error {
	const q = `
		UPDATE admin_credentials
		SET password_hash = $2,
		    updated_at = NOW()
		WHERE id = $1
	`

	if _, err := r.db.Exec(ctx, q, credentialID, newPasswordHash); err != nil {
		return fmt.Errorf("update admin password: %w", err)
	}

	return nil
}

// UpdateLastLogin updates the last_login_at timestamp
func (r *AdminCredentialsRepository) UpdateLastLogin(
	ctx context.Context,
	credentialID string,
) error {
	const q = `
		UPDATE admin_credentials
		SET last_login_at = NOW()
		WHERE id = $1
	`

	if _, err := r.db.Exec(ctx, q, credentialID); err != nil {
		return fmt.Errorf("update last login: %w", err)
	}

	return nil
}

// GetByUserID fetches admin credentials by user ID
func (r *AdminCredentialsRepository) GetByUserID(
	ctx context.Context,
	userID string,
) (*models.AdminCredential, error) {
	const q = `
		SELECT id, username, password_hash, user_id, is_active, created_at, updated_at, last_login_at
		FROM admin_credentials
		WHERE user_id = $1
	`

	row := r.db.QueryRow(ctx, q, userID)

	var cred models.AdminCredential
	err := row.Scan(
		&cred.ID, &cred.Username, &cred.PasswordHash, &cred.UserID,
		&cred.IsActive, &cred.CreatedAt, &cred.UpdatedAt, &cred.LastLoginAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get admin credential by user id: %w", err)
	}

	return &cred, nil
}
