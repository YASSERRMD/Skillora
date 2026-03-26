package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skillora/backend/internal/models"
)

// LLMRepository handles all DB operations on the llm_providers table.
type LLMRepository struct {
	db *pgxpool.Pool
}

// NewLLMRepository constructs an LLMRepository.
func NewLLMRepository(db *pgxpool.Pool) *LLMRepository {
	return &LLMRepository{db: db}
}

// GetActiveProvidersByUseCase fetches all active providers for a specific use case,
// ordered by priority (1=Highest). Useful for fallback routing.
func (r *LLMRepository) GetActiveProvidersByUseCase(ctx context.Context, useCase string) ([]models.LLMProvider, error) {
	const q = `
		SELECT id, provider_name, model_name, encrypted_api_key, use_case, priority, is_active, created_at, updated_at
		FROM llm_providers
		WHERE use_case = $1 AND is_active = true
		ORDER BY priority ASC
	`
	rows, err := r.db.Query(ctx, q, useCase)
	if err != nil {
		return nil, fmt.Errorf("query llm_providers for %s: %w", useCase, err)
	}
	defer rows.Close()

	var providers []models.LLMProvider
	for rows.Next() {
		var p models.LLMProvider
		if err := rows.Scan(
			&p.ID, &p.ProviderName, &p.ModelName, &p.EncryptedAPIKey,
			&p.UseCase, &p.Priority, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan llm_provider row for %s: %w", useCase, err)
		}
		providers = append(providers, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error for %s: %w", useCase, err)
	}

	return providers, nil
}

// InsertProvider securely stores a new LLM provider config.
func (r *LLMRepository) InsertProvider(
	ctx context.Context, providerName, modelName string, encryptedAPIKey []byte, useCase string, priority int,
) (*models.LLMProvider, error) {
	const q = `
		INSERT INTO llm_providers (provider_name, model_name, encrypted_api_key, use_case, priority, is_active)
		VALUES ($1, $2, $3, $4, $5, true)
		RETURNING id, provider_name, model_name, encrypted_api_key, use_case, priority, is_active, created_at, updated_at
	`
	row := r.db.QueryRow(ctx, q, providerName, modelName, encryptedAPIKey, useCase, priority)

	var p models.LLMProvider
	if err := row.Scan(
		&p.ID, &p.ProviderName, &p.ModelName, &p.EncryptedAPIKey,
		&p.UseCase, &p.Priority, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("insert llm_provider: %w", err)
	}
	return &p, nil
}
