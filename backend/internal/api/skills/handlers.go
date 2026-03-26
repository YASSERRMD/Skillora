// Package skills provides handlers for the platform's taxonomy.
package skills

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/agents"
	"github.com/skillora/backend/internal/api"
	"github.com/skillora/backend/internal/models"
	"github.com/skillora/backend/internal/repository"
)

// Handler manages HTTP endpoints for categories and skills.
type Handler struct {
	skillRepo *repository.SkillRepository
	userSkill *repository.UserSkillRepository
	agent     *agents.AppraisalAgent
}

// NewHandler constructs a taxonomy Handler.
func NewHandler(
	skillRepo *repository.SkillRepository,
	userSkill *repository.UserSkillRepository,
	agent *agents.AppraisalAgent,
) *Handler {
	return &Handler{
		skillRepo: skillRepo,
		userSkill: userSkill,
		agent:     agent,
	}
}

// GetCategories returns all parent-level skill categories.
// GET /api/v1/categories
func (h *Handler) GetCategories(c *gin.Context) {
	cats, err := h.skillRepo.GetCategories(c.Request.Context())
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

	skills, err := h.skillRepo.GetCategorySkills(c.Request.Context(), catID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills for category"})
		return
	}
	c.JSON(http.StatusOK, skills)
}

// AppraiseReq is the JSON payload for user skill additions.
type AppraiseReq struct {
	CategoryID  string `json:"category_id" binding:"required"`
	SkillID     string `json:"skill_id"    binding:"required"`
	Description string `json:"description" binding:"required,min=20,max=1000"`
	
	// Provided strictly for passing to the AI quickly without hitting the DB twice if preferred
	CategoryName string `json:"category_name" binding:"required"`
	SkillName    string `json:"skill_name"    binding:"required"`
}

// PostAppraise runs the AppraisalAgent against a user's proposed skill.
// POST /api/v1/skills/appraise
func (h *Handler) PostAppraise(c *gin.Context) {
	userID := api.GetUserID(c)

	var req AppraiseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Ask the AI Agent to appraise the requested skill based on description.
	result, err := h.agent.DraftAppraisal(c.Request.Context(), req.CategoryName, req.SkillName, req.Description)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "AI appraisal engine unavailable"})
		return
	}

	// 2. If it is valid, store it immediately in the Database as a mapping.
	if result.IsValidSkill {
		userSkill := models.UserSkill{
			UserID:           userID,
			SkillID:          req.SkillID,
			ProficiencyLevel: result.Proficiency,
			CreditValue:      result.CreditValue,
			IsVerified:       true,
		}
		if err := h.userSkill.AddUserSkill(c.Request.Context(), userSkill); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save verified skill mapping"})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}
