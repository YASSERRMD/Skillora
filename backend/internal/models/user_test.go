package models_test

import (
	"testing"
	"time"

	"github.com/skillora/backend/internal/models"
)

func TestUserStruct(t *testing.T) {
	now := time.Now()
	avatarURL := "https://example.com/avatar.png"

	u := models.User{
		ID:        "uuid-1234",
		GoogleID:  "google-abc",
		Email:     "alice@example.com",
		FullName:  "Alice Smith",
		AvatarURL: &avatarURL,
		Bio:       "Go developer",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if u.ID != "uuid-1234" {
		t.Errorf("ID mismatch: %s", u.ID)
	}
	if u.Email != "alice@example.com" {
		t.Errorf("Email mismatch: %s", u.Email)
	}
	if u.AvatarURL == nil || *u.AvatarURL != avatarURL {
		t.Error("AvatarURL mismatch")
	}
}

func TestUserAvatarURLNullable(t *testing.T) {
	u := models.User{Email: "bob@example.com"}
	if u.AvatarURL != nil {
		t.Error("AvatarURL should be nil by default")
	}
}
