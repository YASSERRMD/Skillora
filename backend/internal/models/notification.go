package models

import "time"

// NotificationType characterises the reason for a notification.
type NotificationType string

const (
	NotifBarterProposed  NotificationType = "barter_proposed"
	NotifBarterAccepted  NotificationType = "barter_accepted"
	NotifBarterCompleted NotificationType = "barter_completed"
	NotifBarterCancelled NotificationType = "barter_cancelled"
	NotifSystem          NotificationType = "system"
)

// Notification represents a platform event pushed to a user.
type Notification struct {
	ID        string           `db:"id"         json:"id"`
	UserID    string           `db:"user_id"    json:"user_id"`
	Title     string           `db:"title"      json:"title"`
	Body      string           `db:"body"       json:"body"`
	Type      NotificationType `db:"type"       json:"type"`
	IsRead    bool             `db:"is_read"    json:"is_read"`
	Metadata  []byte           `db:"metadata"   json:"metadata,omitempty"` // Raw JSONB
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
}
