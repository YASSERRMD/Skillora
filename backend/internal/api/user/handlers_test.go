package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/api/user"
	"github.com/skillora/backend/internal/config"
)

// setupMockDB spins up a mock DB by using a miniredis server (just for config loads)
// and returns a mock handler (which still technically needs a valid pg connection
// to test the repo, but since we're writing atomic commits and don't have a live DB
// in unit tests here, we'll test the handler logic with a nil repo, which will panic.
// Let's refactor the handler to accept an interface, or since we are keeping it simple,
// we will just do an integration-style skip or basic binding test.)
// Wait, for atomic and robust Go code, I should pass an interface if I want to unit test handlers.
// Since the blueprint doesn't strictly dictate clean architecture interfaces vs struct methods,
// I'll skip deep repo testing here unless we have testcontainers, but I can test the JSON binding
// and RequireAuth pipeline.

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

func mockRequireAuth(mockUserID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", mockUserID)
		c.Next()
	}
}

// Without an interface for repo, we can't fully mock DB hits in handler tests easily
// without testcontainers. We'll test the binding and error paths.
func TestUpdateMe_BindingErrors(t *testing.T) {
	// A nil repo will panic if it reaches DB, so we only test validation failure.
	handler := user.NewHandler(nil)
	router := gin.New()
	router.PUT("/me", mockRequireAuth("test-id"), handler.UpdateMe)

	// Test 1: Empty body → 400
	req, _ := http.NewRequest("PUT", "/me", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for empty body, got %d", w.Code)
	}

	// Test 2: Bio too long → 400
	longBio := string(make([]byte, 501)) // > 500 chars limit
	body, _ := json.Marshal(map[string]string{"bio": longBio})
	req2, _ := http.NewRequest("PUT", "/me", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for too-long bio, got %d", w2.Code)
	}
}
