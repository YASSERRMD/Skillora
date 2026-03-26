// Package skills provides handlers for the platform's taxonomy.
package skills

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/repository"
)

// Handler manages HTTP endpoints for categories and skills.
type Handler struct {
	repo *repository.SkillRepository
}

// NewHandler constructs a taxonomy Handler.
func NewHandler(repo *repository.SkillRepository) *Handler {
	return &Handler{repo: repo}
}

// GetCategories returns all parent-level skill categories.
// GET /api/v1/categories
func (h *Handler) GetCategories(c *gin.Context) {
	cats, err := h.repo.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load categories"})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// GetCategorySkills returns all skills within a specific category.
// GET /api/v1/categories/:id/skills
func (h *Handler) GetCategorySkills(c *gin.Context) {
	catID := c.Param("id")
	if catID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing category id"})
		return
	}

	skills, err := h.repo.GetCategorySkills(c.Request.Context(), catID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills for category"})
		return
	}
	c.JSON(http.StatusOK, skills)
}
