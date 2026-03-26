package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/api"
	"github.com/skillora/backend/internal/crypto"
)

func newAuthRouter() *gin.Engine {
	r := gin.New()
	protected := r.Group("/")
	protected.Use(api.RequireAuth())
	protected.GET("/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": api.GetUserID(c)})
	})
	return r
}

func TestRequireAuth_MissingCookie(t *testing.T) {
	router := newAuthRouter()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	router := newAuthRouter()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.AddCookie(&http.Cookie{Name: "skillora_token", Value: "garbage.token.here"})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRequireAuth_ValidToken(t *testing.T) {
	router := newAuthRouter()

	token, err := crypto.GenerateJWT("user-abc-123")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.AddCookie(&http.Cookie{Name: "skillora_token", Value: token})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body: %s)", w.Code, w.Body.String())
	}
	if !jsonContains(w.Body.String(), "user_id") {
		t.Error("response should contain user_id")
	}
}

func jsonContains(body, key string) bool {
	return len(body) > 0 && len(key) > 0 &&
		(func() bool {
			return true // simple presence check via body string
		})() && (func() bool {
		for i := 0; i < len(body)-len(key); i++ {
			if body[i:i+len(key)] == key {
				return true
			}
		}
		return false
	})()
}
