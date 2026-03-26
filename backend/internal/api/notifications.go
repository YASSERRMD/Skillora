package api

import (
	"context"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/models"
	"github.com/skillora/backend/internal/repository"
)

// NotificationHub pushes in-memory updates to SSE channels keyed by user_id.
type NotificationHub struct {
	mu       sync.RWMutex
	channels map[string]chan models.Notification
	repo     *repository.NotificationRepository
}

// NewNotificationHub constructs a hub.
func NewNotificationHub(repo *repository.NotificationRepository) *NotificationHub {
	return &NotificationHub{
		channels: make(map[string]chan models.Notification),
		repo:     repo,
	}
}

// Subscribe registers a channel for the user's notifications.
func (h *NotificationHub) Subscribe(userID string) chan models.Notification {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan models.Notification, 8)
	h.channels[userID] = ch
	return ch
}

// Unsubscribe removes the channel when the client disconnects.
func (h *NotificationHub) Unsubscribe(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.channels[userID]; ok {
		close(ch)
		delete(h.channels, userID)
	}
}

// Push inserts a notification into the DB and forwards it to the live channel (if present).
func (h *NotificationHub) Push(userID string, n models.Notification) {
	// Store persistently
	if h.repo != nil {
		if err := h.repo.CreateNotification(context.Background(), n); err != nil {
			log.Printf("[notif] db persist error for user %s: %v", userID, err)
		}
	}

	// Forward to live SSE channel if the user is connected
	h.mu.RLock()
	ch, ok := h.channels[userID]
	h.mu.RUnlock()

	if ok {
		select {
		case ch <- n:
		default:
			log.Printf("[notif] channel buffer full for user %s, dropping notif", userID)
		}
	}
}

// SSEHandler streams notifications to the connected browser client via Server-Sent-Events.
// GET /api/v1/notifications/stream
func (h *NotificationHub) SSEHandler(c *gin.Context) {
	userID := GetUserID(c)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ch := h.Subscribe(userID)
	defer h.Unsubscribe(userID)

	c.Stream(func(w io.Writer) bool {
		select {
		case notif, ok := <-ch:
			if !ok {
				return false
			}
			c.SSEvent("notification", notif)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

// GetNotificationsHandler returns unread notifications for the current user.
// GET /api/v1/notifications
func GetNotificationsHandler(repo *repository.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		notifs, err := repo.GetUnread(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load notifications"})
			return
		}
		c.JSON(http.StatusOK, notifs)
	}
}

// MarkNotificationsReadHandler marks all notifications as read.
// POST /api/v1/notifications/read
func MarkNotificationsReadHandler(repo *repository.NotificationRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if err := repo.MarkAllRead(c.Request.Context(), userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark read"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
	}
}
