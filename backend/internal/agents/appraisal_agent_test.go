package agents_test

import (
	"context"
	"testing"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/agents"
	"github.com/skillora/backend/internal/config"
	"github.com/skillora/backend/internal/llm"
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

func TestAppraisalAgent_Initialization(t *testing.T) {
	mgr := llm.NewManager(nil)
	router := llm.NewRouter(mgr)

	agent := agents.NewAppraisalAgent(router)
	if agent == nil {
		t.Fatal("expected agent to be initialized")
	}

	// Because we can't reliably test a live LLM router without mocking an interface
	// for the engine, we verify structural initialization and context setup.
	// Production testing would mock the external HTTP calls.
	_ = context.Background()
}
