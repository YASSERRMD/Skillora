package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// OpenAIAdapter implements LLMEngine using the OpenAI REST API.
type OpenAIAdapter struct {
	client *http.Client
	apiKey string
	model  string
}

// NewOpenAIAdapter creates a new adapter with the decrypted API key and target model.
func NewOpenAIAdapter(apiKey, model string) *OpenAIAdapter {
	return &OpenAIAdapter{
		client: &http.Client{},
		apiKey: apiKey,
		model:  model,
	}
}

// GenerateJSON calls the Chat Completions API demanding JSON output.
func (a *OpenAIAdapter) GenerateJSON(ctx context.Context, prompt string) (string, error) {
	reqBody, _ := json.Marshal(map[string]any{
		"model": a.model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant. Always output valid JSON."},
			{"role": "user", "content": prompt},
		},
		"response_format": map[string]string{"type": "json_object"},
	})

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai text req: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("openai returned status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("openai JSON decode: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai returned 0 choices")
	}

	return result.Choices[0].Message.Content, nil
}

// GenerateEmbedding calls the Embeddings API to embed text.
func (a *OpenAIAdapter) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	reqBody, _ := json.Marshal(map[string]any{
		"model": "text-embedding-3-small",
		"input": text,
	})

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("openai embed req: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai embed status %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("openai embed decode: %w", err)
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("openai embed returned 0 data arrays")
	}

	return result.Data[0].Embedding, nil
}
