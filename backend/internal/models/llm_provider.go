package models

import "time"

// LLMUseCase constants define the types of tasks an LLM provider handles.
const (
	UseCaseGeneral          = "general"
	UseCaseEmbedding        = "embedding"
	UseCaseCourseGeneration = "course_generation"
	UseCaseMediator         = "mediator"
)

// LLMProvider represents an AI text or embedding model configuration.
type LLMProvider struct {
	ID              string    `db:"id"               json:"id"`
	ProviderName    string    `db:"provider_name"    json:"provider_name"` // "openai", "anthropic", "deepseek"
	ModelName       string    `db:"model_name"       json:"model_name"`
	EncryptedAPIKey []byte    `db:"encrypted_api_key" json:"-"`            // Never exposed to frontend JSON
	UseCase         string    `db:"use_case"         json:"use_case"`
	Priority        int       `db:"priority"         json:"priority"`      // 1=Primary, 2=Fallback
	IsActive        bool      `db:"is_active"        json:"is_active"`
	CreatedAt       time.Time `db:"created_at"       json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"       json:"updated_at"`
}
