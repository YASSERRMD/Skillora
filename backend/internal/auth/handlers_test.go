package auth_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/skillora/backend/internal/auth"
	"github.com/skillora/backend/internal/config"
	"github.com/skillora/backend/internal/db"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("AES_MASTER_KEY", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	os.Setenv("FRONTEND_URL", "http://localhost:3000")
	config.Load()
	gin.SetMode(gin.TestMode)
}

func newMiniRedis(t *testing.T) {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(mr.Close)
	db.RDB = redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

func newOAuthCfg() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "test-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/api/v1/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func TestGoogleLogin_RedirectsToGoogle(t *testing.T) {
	newMiniRedis(t)

	handler := auth.NewHandler(newOAuthCfg(), nil)
	router := gin.New()
	router.GET("/api/v1/auth/google/login", handler.GoogleLogin)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/login", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307 redirect, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if loc == "" {
		t.Fatal("expected Location header for Google redirect")
	}
	// Should redirect to accounts.google.com
	if len(loc) < 20 {
		t.Errorf("Location header looks too short: %s", loc)
	}
}

func TestGoogleCallback_InvalidState(t *testing.T) {
	newMiniRedis(t)

	handler := auth.NewHandler(newOAuthCfg(), nil)
	router := gin.New()
	router.GET("/api/v1/auth/google/callback", handler.GoogleCallback)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/callback?state=bad-state&code=any", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Invalid state → 400
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid state, got %d", w.Code)
	}
}
