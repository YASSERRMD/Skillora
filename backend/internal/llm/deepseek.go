package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DeepSeekAdapter implements LLMEngine using DeepSeek's OpenAI-compatible API.
type DeepSeekAdapter struct {
	client *http.Client
	apiKey string
	model  string
}

// NewDeepSeekAdapter constructs the DeepSeek adapter.
func NewDeepSeekAdapter(apiKey, model string) *DeepSeekAdapter {
	return &DeepSeekAdapter{
		client: &http.Client{},
		apiKey: apiKey,
		model:  model,
	}
}

// GenerateJSON calls DeepSeek Chat Completions demanding JSON output.
func (a *DeepSeekAdapter) GenerateJSON(ctx context.Context, prompt string) (string, error) {
	reqBody, _ := json.Marshal(map[string]any{
		"model": a.model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant. Always output valid JSON."},
			{"role": "user", "content": prompt},
		},
		"response_format": map[string]string{"type": "json_object"},
	})

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.deepseek.com/chat/completions", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("deepseek text req: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		// Return specific error format so Router can detect 429/500
		return "", fmt.Errorf("deepseek returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("deepseek JSON decode: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("deepseek returned 0 choices")
	}

	return result.Choices[0].Message.Content, nil
}

// GenerateEmbedding calls the DeepSeek embedding endpoint (if they support it), or throws error.
// Currently DeepSeek does not natively expose public embeddings on all models. 
// Fallback router will handle this if the use_case relies on vector indexing.
func (a *DeepSeekAdapter) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// DeepSeek does not currently have an OpenAI-parity vector embedding endpoint guaranteed for production.
	return nil, fmt.Errorf("deepseek embeddings not natively supported, fallback router required")
}
