package agents

import (
	"context"
	"fmt"
	"log"

	"github.com/skillora/backend/internal/llm"
	"github.com/skillora/backend/internal/models"
	"github.com/skillora/backend/internal/repository"
)

// EmbeddingAgent embeds skills into vector space for matching engine queries.
type EmbeddingAgent struct {
	router    *llm.Router
	vectorRepo *repository.VectorRepository
}

// NewEmbeddingAgent constructs the embedding pipeline agent.
func NewEmbeddingAgent(router *llm.Router, repo *repository.VectorRepository) *EmbeddingAgent {
	return &EmbeddingAgent{router: router, vectorRepo: repo}
}

// EmbedAndStoreSkill generates a vector embedding for a skill and stores it in the DB.
func (a *EmbeddingAgent) EmbedAndStoreSkill(ctx context.Context, skill models.Skill, categoryName string) error {
	// Compose the semantic text: category + skill name + description
	text := fmt.Sprintf("%s: %s. %s", categoryName, skill.Name, skill.Description)

	embedding, err := a.router.GenerateEmbedding(ctx, models.UseCaseEmbedding, text)
	if err != nil {
		return fmt.Errorf("EmbedAndStoreSkill generate: %w", err)
	}

	if err := a.vectorRepo.UpsertSkillEmbedding(ctx, skill.ID, embedding); err != nil {
		return fmt.Errorf("EmbedAndStoreSkill upsert: %w", err)
	}

	log.Printf("[embedding] stored vector for skill %s (%s)", skill.Name, skill.ID)
	return nil
}

// EmbedQuery generates a temporary embedding for a free-text query (for matching searches).
func (a *EmbeddingAgent) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	embedding, err := a.router.GenerateEmbedding(ctx, models.UseCaseEmbedding, text)
	if err != nil {
		return nil, fmt.Errorf("EmbedQuery: %w", err)
	}
	return embedding, nil
}
