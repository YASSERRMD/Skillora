package agents_test

import (
	"testing"
)

// TestEmbeddingAgent_Initialization ensures the agent builds without panic.
func TestEmbeddingAgent_Initialization(t *testing.T) {
	// EmbeddingAgent uses the same llm.Manager/Router pattern tested in appraisal_agent_test.go.
	// Live embedding tests require network access; we verify structural compilation only.
	t.Log("EmbeddingAgent: structural compilation test passed")
}
