package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skillora/backend/internal/models"
)

// NotificationRepository stores and retrieves user notification records.
type NotificationRepository struct {
	db *pgxpool.Pool
}

// NewNotificationRepository constructs the dependency.
func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}

// CreateNotification inserts a new event notification for a user.
func (r *NotificationRepository) CreateNotification(ctx context.Context, n models.Notification) error {
	const q = `
		INSERT INTO notifications (user_id, title, body, type, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`
	metadata := n.Metadata
	if metadata == nil {
		metadata = []byte(`{}`)
	}
	_, err := r.db.Exec(ctx, q, n.UserID, n.Title, n.Body, string(n.Type), metadata)
	if err != nil {
		return fmt.Errorf("CreateNotification: %w", err)
	}
	return nil
}

// GetUnread returns all unread notifications for a user sorted by creation time DESC.
func (r *NotificationRepository) GetUnread(ctx context.Context, userID string) ([]models.Notification, error) {
	const q = `
		SELECT id, user_id, title, body, type, is_read, metadata, created_at
		FROM notifications
		WHERE user_id = $1 AND is_read = false
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUnread: %w", err)
	}
	defer rows.Close()

	var list []models.Notification
	for rows.Next() {
		var n models.Notification
		var notifType string
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &notifType, &n.IsRead, &n.Metadata, &n.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetUnread scan: %w", err)
		}
		n.Type = models.NotificationType(notifType)
		list = append(list, n)
	}
	if list == nil {
		list = make([]models.Notification, 0)
	}
	return list, rows.Err()
}

// MarkAllRead marks all unread notifications for a user as read.
func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID string) error {
	const q = `UPDATE notifications SET is_read = true WHERE user_id = $1 AND is_read = false`
	_, err := r.db.Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("MarkAllRead: %w", err)
	}
	return nil
}
