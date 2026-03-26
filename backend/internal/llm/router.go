package llm

import (
	"context"
	"fmt"
	"log"
	"strings"
)

// Router handles executing AI prompts by attempting providers in priority order.
type Router struct {
	manager *Manager
}

// NewRouter constructs a new AI router securely hooked up to the Manager cache.
func NewRouter(manager *Manager) *Router {
	return &Router{manager: manager}
}

// buildEngine instantiates the correct adapter struct.
func buildEngine(p ProviderConfig) (LLMEngine, error) {
	switch p.ProviderName {
	case "openai":
		return NewOpenAIAdapter(p.DecryptedAPIKey, p.ModelName), nil
	case "anthropic":
		return NewAnthropicAdapter(p.DecryptedAPIKey, p.ModelName), nil
	case "deepseek":
		return NewDeepSeekAdapter(p.DecryptedAPIKey, p.ModelName), nil
	default:
		return nil, fmt.Errorf("unknown provider %s", p.ProviderName)
	}
}

// isRetryableError checks if the LLM provider failed due to rate limits or internal server errors.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	// E.g. "openai returned status 429" or "deepseek returned status 503"
	return strings.Contains(msg, "status 429") ||
		strings.Contains(msg, "status 500") ||
		strings.Contains(msg, "status 502") ||
		strings.Contains(msg, "status 503") ||
		strings.Contains(msg, "status 504") ||
		strings.Contains(msg, "fallback router required") // Specific condition to fallback embedding to OpenAI
}

// GenerateJSON attempts to fetch a JSON string from the providers of a given use case.
func (r *Router) GenerateJSON(ctx context.Context, useCase, prompt string) (string, error) {
	providers := r.manager.GetProviders(useCase)
	if len(providers) == 0 {
		return "", fmt.Errorf("no active LLM providers configured for use case: %s", useCase)
	}

	for _, p := range providers {
		engine, err := buildEngine(p)
		if err != nil {
			log.Printf("[router] skip %s (%s): %v", p.ProviderName, p.ModelName, err)
			continue
		}

		log.Printf("[router] attempting GenerateJSON via %s (priority %d)", p.ProviderName, p.Priority)
		result, err := engine.GenerateJSON(ctx, prompt)
		if err == nil {
			return result, nil // Success
		}

		log.Printf("[router] provider %s failed: %v", p.ProviderName, err)
		if !isRetryableError(err) {
			// If it's a 400 Bad Request, we don't retry. The prompt form is bad.
			return "", err
		}
		// Fall through to next provider logic natively because of the loop.
		log.Printf("[router] falling back to next provider...")
	}

	return "", fmt.Errorf("all providers exhausted for use case: %s", useCase)
}

// GenerateEmbedding attempts to vectorize text through prioritized providers.
func (r *Router) GenerateEmbedding(ctx context.Context, useCase, text string) ([]float32, error) {
	providers := r.manager.GetProviders(useCase)
	if len(providers) == 0 {
		return nil, fmt.Errorf("no active embedding providers for use case: %s", useCase)
	}

	for _, p := range providers {
		engine, err := buildEngine(p)
		if err != nil {
			continue
		}

		result, err := engine.GenerateEmbedding(ctx, text)
		if err == nil {
			return result, nil
		}

		if !isRetryableError(err) {
			return nil, err
		}
	}

	return nil, fmt.Errorf("all embedding providers exhausted for use case: %s", useCase)
}
