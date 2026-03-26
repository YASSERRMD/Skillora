package admin

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/crypto"
	"github.com/skillora/backend/internal/repository"
)

// LLMHandler manages LLM routing settings for admins.
type LLMHandler struct {
	repo *repository.LLMRepository
}

// NewLLMHandler constructs the LLM settings handler.
func NewLLMHandler(repo *repository.LLMRepository) *LLMHandler {
	return &LLMHandler{repo: repo}
}

// AddProviderReq is the JSON payload for configuring a new AI provider.
type AddProviderReq struct {
	ProviderName string `json:"provider_name" binding:"required,oneof=openai anthropic deepseek"`
	ModelName    string `json:"model_name"    binding:"required"`
	APIKey       string `json:"api_key"       binding:"required,min=10"` // raw key
	UseCase      string `json:"use_case"      binding:"required,oneof=general embedding course_generation mediator"`
	Priority     int    `json:"priority"      binding:"required,min=1"`
}

// PostLLMProvider securely encrypts the provided API key and registers the LLM model.
// POST /api/v1/admin/llm-providers
func (h *LLMHandler) PostLLMProvider(c *gin.Context) {
	var req AddProviderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedKey, err := crypto.Encrypt(strings.TrimSpace(req.APIKey))
	if err != nil {
		log.Printf("[admin] encryption failed for provider %s: %v", req.ProviderName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to secure provider key"})
		return
	}

	provider, err := h.repo.InsertProvider(
		c.Request.Context(), req.ProviderName, req.ModelName, encryptedKey, req.UseCase, req.Priority,
	)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "a provider with this priority already exists for this use case"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save provider"})
		return
	}

	c.JSON(http.StatusCreated, provider)
}

// GetLLMProviders lists all configured LLM providers (keys omitted by repository query).
// GET /api/v1/admin/llm-providers
func (h *LLMHandler) GetLLMProviders(c *gin.Context) {
	providers, err := h.repo.GetAllProviders(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list providers"})
		return
	}

	// Safe to return empty array instead of null
	if providers == nil {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	c.JSON(http.StatusOK, providers)
}
