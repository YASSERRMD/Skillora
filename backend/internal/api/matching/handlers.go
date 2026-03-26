package matching

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/agents"
	"github.com/skillora/backend/internal/api"
	"github.com/skillora/backend/internal/repository"
)

// Handler manages AI-driven skill matching endpoints.
type Handler struct {
	vectorRepo     *repository.VectorRepository
	embeddingAgent *agents.EmbeddingAgent
}

// NewHandler constructs the matching handler.
func NewHandler(repo *repository.VectorRepository, agent *agents.EmbeddingAgent) *Handler {
	return &Handler{vectorRepo: repo, embeddingAgent: agent}
}

// GetMatches finds the top-k skills most semantically similar to the user's query.
// GET /api/v1/match?q=<query_text>&limit=<n>
func (h *Handler) GetMatches(c *gin.Context) {
	userID := api.GetUserID(c)
	queryText := c.Query("q")
	if queryText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	// Step 1: Vectorize the search query.
	queryVec, err := h.embeddingAgent.EmbedQuery(c.Request.Context(), queryText)
	if err != nil {
		log.Printf("[matching] embed query error: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "embedding engine unavailable"})
		return
	}

	// Step 2: ANN search excluding the querying user's own skills.
	matches, err := h.vectorRepo.FindSimilarSkills(c.Request.Context(), queryVec, userID, limit)
	if err != nil {
		log.Printf("[matching] vector search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "matching engine error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   queryText,
		"results": matches,
	})
}
