package barter_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	barterapi "github.com/skillora/backend/internal/api/barter"
	"github.com/skillora/backend/internal/config"
)

func init() {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("AES_MASTER_KEY", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	config.Load()
	gin.SetMode(gin.TestMode)
}

func TestPostPropose_SelfBarter(t *testing.T) {
	h := barterapi.NewHandler(nil)
	r := gin.New()
	r.POST("/barters", func(c *gin.Context) {
		// Inject fake user_id matching the receiver_id to trigger self-barter check.
		c.Set("user_id", "same-user-uuid")
		h.PostPropose(c)
	})

	body := `{"receiver_id":"same-user-uuid","initiator_skill_id":"s1","receiver_skill_id":"s2","credit_amount":10}`
	req := httptest.NewRequest(http.MethodPost, "/barters", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for self-barter, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPostPropose_MissingFields(t *testing.T) {
	h := barterapi.NewHandler(nil)
	r := gin.New()
	r.POST("/barters", func(c *gin.Context) {
		c.Set("user_id", "initiator-uuid")
		h.PostPropose(c)
	})

	body := `{"receiver_id":"other-uuid"}` // Missing required fields
	req := httptest.NewRequest(http.MethodPost, "/barters", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", w.Code)
	}
}

func TestPatchBarterStatus_InvalidEnum(t *testing.T) {
	h := barterapi.NewHandler(nil)
	r := gin.New()
	r.PATCH("/barters/:id/status", func(c *gin.Context) {
		c.Set("user_id", "user-uuid")
		h.PatchBarterStatus(c)
	})

	body := `{"status":"hacked"}` // Not in oneof=accepted cancelled
	req := httptest.NewRequest(http.MethodPatch, "/barters/barter-id/status", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid status enum, got %d", w.Code)
	}
}
