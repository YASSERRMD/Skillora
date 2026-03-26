package admin_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/api/admin"
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

func TestPostLLMProvider_Validation(t *testing.T) {
	// Nil repo translates to runtime panic if it reaches DB logic.
	// But it will panic only *after* validation if validation passes.
	handler := admin.NewLLMHandler(nil)
	router := gin.New()
	router.POST("/providers", handler.PostLLMProvider)

	tests := []struct {
		name     string
		payload  string
		expected int
	}{
		{
			"missing fields",
			`{"provider_name":"openai"}`,
			http.StatusBadRequest,
		},
		{
			"invalid provider enum",
			`{"provider_name":"gemini", "model_name":"pro", "api_key":"1234567890", "use_case":"general", "priority":1}`,
			http.StatusBadRequest, // Not in oneof
		},
		{
			"invalid use_case enum",
			`{"provider_name":"openai", "model_name":"pro", "api_key":"1234567890", "use_case":"hacker", "priority":1}`,
			http.StatusBadRequest,
		},
		{
			"short API key",
			`{"provider_name":"anthropic", "model_name":"claude", "api_key":"short", "use_case":"mediator", "priority":1}`,
			http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/providers", bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			if w.Code != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, w.Code)
			}
		})
	}
}
