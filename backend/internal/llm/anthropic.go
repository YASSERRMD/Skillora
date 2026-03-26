package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// AnthropicAdapter implements LLMEngine using the Anthropic API.
type AnthropicAdapter struct {
	client *http.Client
	apiKey string
	model  string
}

// NewAnthropicAdapter constructs the Anthropic Adapter.
func NewAnthropicAdapter(apiKey, model string) *AnthropicAdapter {
	return &AnthropicAdapter{
		client: &http.Client{},
		apiKey: apiKey,
		model:  model,
	}
}

// GenerateJSON requests a strictly formatted JSON string from Claude.
func (a *AnthropicAdapter) GenerateJSON(ctx context.Context, prompt string) (string, error) {
	// Claude doesn't have a JSON format native specifer like OpenAI,
	// so we use system prompting + explicit instructions.
	systemPrompt := "You areSkillora Appraiser. Parse the input and ONLY output valid JSON. Do not include markdown codeblocks (no ```json text), only raw JSON."

	reqBody, _ := json.Marshal(map[string]any{
		"model":      a.model, // e.g. "claude-3-5-sonnet-20240620"
		"max_tokens": 1024,
		"system":     systemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBody))
	req.Header.Set("x-api-key", a.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("anthropic req: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("anthropic status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("anthropic decode: %w", err)
	}
	if len(result.Content) == 0 {
		return "", errors.New("anthropic returned 0 content blocks")
	}

	return result.Content[0].Text, nil
}

// GenerateEmbedding throws an error because Anthropic currently does not offer
// a native public embeddings API compatible with text-embedding-3 vectors natively for this stack.
// In Phase 14, the router will fallback to OpenAI for embeddings if an anthropic instance is chosen for text.
func (a *AnthropicAdapter) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	return nil, errors.New("anthropic does not support vector embeddings natively; fallback router should handle this")
}
