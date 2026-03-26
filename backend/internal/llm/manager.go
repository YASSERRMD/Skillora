// Package llm manages AI model routing and API key decryption.
package llm

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/skillora/backend/internal/crypto"
	"github.com/skillora/backend/internal/models"
	"github.com/skillora/backend/internal/repository"
)

// ProviderConfig represents an in-memory routed LLM configuration with a decrypted API key.
// It is intended for ephemeral use by the router during a request.
type ProviderConfig struct {
	models.LLMProvider
	DecryptedAPIKey string
}

// Manager handles in-memory caching of LLM providers to avoid database hits on every AI request.
type Manager struct {
	repo *repository.LLMRepository
	mu   sync.RWMutex
	// map[useCase][]ProviderConfig
	cache map[string][]ProviderConfig
}

// NewManager creates an LLM Manager and initializes the empty cache.
func NewManager(repo *repository.LLMRepository) *Manager {
	return &Manager{
		repo:  repo,
		cache: make(map[string][]ProviderConfig),
	}
}

// StartBackgroundSync runs a goroutine that syncs providers from DB every 60 seconds.
func (m *Manager) StartBackgroundSync(ctx context.Context) {
	// First immediate sync
	if err := m.SyncOnce(ctx); err != nil {
		log.Printf("[llm-manager] initial sync failed: %v", err)
	}

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("[llm-manager] background sync stopping")
				return
			case <-ticker.C:
				if err := m.SyncOnce(ctx); err != nil {
					log.Printf("[llm-manager] background sync error: %v", err)
				}
			}
		}
	}()
}

// SyncOnce reads all active providers from DB and rebuilds the memory cache.
func (m *Manager) SyncOnce(ctx context.Context) error {
	// Query all known use cases (we could query distinct, but hardcoding them
	// per the spec is fast enough).
	useCases := []string{
		models.UseCaseGeneral,
		models.UseCaseEmbedding,
		models.UseCaseCourseGeneration,
		models.UseCaseMediator,
	}

	newCache := make(map[string][]ProviderConfig)

	for _, uc := range useCases {
		dbProviders, err := m.repo.GetActiveProvidersByUseCase(ctx, uc)
		if err != nil {
			return err
		}

		var active []ProviderConfig
		for _, dbP := range dbProviders {
			decrypted, err := crypto.Decrypt(dbP.EncryptedAPIKey)
			if err != nil {
				// Don't kill the whole sync if one key is malformed, just log and skip it.
				log.Printf("[llm-manager] failed to decrypt key for provider %s: %v", dbP.ProviderName, err)
				continue
			}

			active = append(active, ProviderConfig{
				LLMProvider:     dbP,
				DecryptedAPIKey: decrypted,
			})
		}
		newCache[uc] = active
	}

	m.mu.Lock()
	m.cache = newCache
	m.mu.Unlock()

	return nil
}

// GetProviders returns the prioritized list of providers for a given use case.
// Returns an empty slice if no providers are configured or cached.
func (m *Manager) GetProviders(useCase string) []ProviderConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy of the slice (which holds copies of the structs, ok since they don't contain pointers)
	// to prevent race conditions on slice mutation, though the router shouldn't mutate.
	cached := m.cache[useCase]
	out := make([]ProviderConfig, len(cached))
	copy(out, cached)
	return out
}
