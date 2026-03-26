package matching_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	matchingapi "github.com/skillora/backend/internal/api/matching"
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

func TestGetMatches_MissingQuery(t *testing.T) {
	h := matchingapi.NewHandler(nil, nil)
	r := gin.New()
	r.GET("/match", func(c *gin.Context) {
		c.Set("user_id", "user-uuid")
		h.GetMatches(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/match", nil) // No ?q=
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing query, got %d: %s", w.Code, w.Body.String())
	}
}
