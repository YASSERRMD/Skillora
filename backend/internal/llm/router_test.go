package llm

import (
	"errors"
	"testing"
)

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		err      error
		expected bool
	}{
		{nil, false},
		{errors.New("openai returned status 400"), false}, // Bad request
		{errors.New("openai returned status 404"), false}, // Not found
		{errors.New("deepseek returned status 500"), true}, // Internal Server Error
		{errors.New("anthropic returned status 429"), true}, // Rate limited
		{errors.New("anthropic returned status 503"), true}, // Overloaded
		{errors.New("fallback router required for embeddings"), true}, // Specific fallback hint
	}

	for _, tc := range tests {
		result := isRetryableError(tc.err)
		if result != tc.expected {
			t.Errorf("Expected %v for error '%v', got %v", tc.expected, tc.err, result)
		}
	}
}
