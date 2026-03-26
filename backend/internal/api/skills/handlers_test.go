package skills_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/api/skills"
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

func TestGetCategories_Integration(t *testing.T) {
	// Nil repo translates to runtime panic only when calling DB.
	handler := skills.NewHandler(nil)
	router := gin.New()
	router.GET("/categories", handler.GetCategories)

	// Will panic as expected with nil repo if it tries to hit DB.
	// We'll skip deep testing here without a DI framework or integration test DB.
	// Just verifies route definition binds cleanly.
	t.Skip("skipping db-dependent test")
}

func TestGetCategorySkills_MissingParam(t *testing.T) {
	handler := skills.NewHandler(nil)
	router := gin.New()
	router.GET("/categories/:id/skills", handler.GetCategorySkills)
	// Missing id param (Gin usually prevents this by routing logic, but handler handles explicit empty check)
	req := httptest.NewRequest(http.MethodGet, "/categories//skills", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Request with missing 'id' hits the route but with empty Param if manually constructed,
	// or our handler intercepts it and returns 400.
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 from handler for missing param, got %d", w.Code)
	}
}
