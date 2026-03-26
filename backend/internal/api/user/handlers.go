// Package user provides user-specific HTTP handlers.
package user

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/api"
	"github.com/skillora/backend/internal/repository"
)

// Handler holds user dependencies.
type Handler struct {
	repo          *repository.UserRepository
	userSkillRepo *repository.UserSkillRepository
}

// NewHandler constructs a user Handler.
func NewHandler(repo *repository.UserRepository, userSkillRepo *repository.UserSkillRepository) *Handler {
	return &Handler{repo: repo, userSkillRepo: userSkillRepo}
}

// GetMe returns the authenticated user's profile.
// GET /api/v1/users/me
func (h *Handler) GetMe(c *gin.Context) {
	userID := api.GetUserID(c)

	user, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateBioReq is the JSON body for PUT /api/v1/users/me
type UpdateBioReq struct {
	Bio string `json:"bio" binding:"required,max=500"`
}

// UpdateMe updates the authenticated user's profile (specifically, bio).
// PUT /api/v1/users/me
func (h *Handler) UpdateMe(c *gin.Context) {
	userID := api.GetUserID(c)

	var req UpdateBioReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sanitize slightly (trim whitespace).
	req.Bio = strings.TrimSpace(req.Bio)

	if err := h.repo.UpdateBio(c.Request.Context(), userID, req.Bio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user bio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "bio": req.Bio})
}

// GetMySkills returns all verified skills for the authenticated user.
// GET /api/v1/users/skills
func (h *Handler) GetMySkills(c *gin.Context) {
	userID := api.GetUserID(c)
	skills, err := h.userSkillRepo.GetUserSkills(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills"})
		return
	}
	c.JSON(http.StatusOK, skills)
}

// GetUserSkills returns all verified skills for a specific user ID.
// GET /api/v1/users/:id/skills
func (h *Handler) GetUserSkills(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
		return
	}
	skills, err := h.userSkillRepo.GetUserSkills(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load skills"})
		return
	}
	c.JSON(http.StatusOK, skills)
}
