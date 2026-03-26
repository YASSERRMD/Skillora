package llm_test

import (
	"os"
	"testing"

	"github.com/skillora/backend/internal/config"
	"github.com/skillora/backend/internal/llm"
	"github.com/skillora/backend/internal/models"
)

// We verify it instantiates and returns empty map when DB is un-reachable (mocked via nil struct).
// Because we aren't using an interface for the repository, we mostly rely on compiler checks
// and explicit edge-case behaviors (like a nil cache without panic).
// However, I can initialize the manager safely.

func init() {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("AES_MASTER_KEY", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	config.Load()
}

func TestManager_GetProviders(t *testing.T) {
	// Provide a nil repo -- it won't panic on NewManager because it allocates securely.
	mgr := llm.NewManager(nil)

	// Since SyncOnce hasn't run, getting providers for a use case returns empty slice, not panic.
	p := mgr.GetProviders(models.UseCaseGeneral)
	if len(p) != 0 {
		t.Errorf("Expected 0 providers initially, got %d", len(p))
	}
}

// In a real environment with a DB, we would do:
// 1. repo.InsertProvider(...)
// 2. mgr.SyncOnce(ctx)
// 3. providers := mgr.GetProviders(...)
// 4. Assert length == 1 and DecryptedAPIKey == "original key"
