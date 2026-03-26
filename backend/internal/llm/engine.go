package llm

import "context"

// LLMEngine defines the standard contract required by Skillora for interacting
// with an underlying language model provider.
type LLMEngine interface {
	// GenerateJSON sends a prompt and expects a strict JSON response.
	GenerateJSON(ctx context.Context, prompt string) (string, error)

	// GenerateEmbedding converts text into a fixed-length float32 vector map.
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
}
