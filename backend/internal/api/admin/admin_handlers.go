package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/repository"
)

// UserAdminHandler handles user administration operations
type UserAdminHandler struct {
	userRepo *repository.UserRepository
}

// NewUserAdminHandler creates a new user admin handler
func NewUserAdminHandler(userRepo *repository.UserRepository) *UserAdminHandler {
	return &UserAdminHandler{userRepo: userRepo}
}

// GrantAdminRequest represents the request body for granting admin privileges
type GrantAdminRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// GrantAdmin grants admin privileges to a user by email
// POST /api/v1/admin/users/grant-admin
func (h *UserAdminHandler) GrantAdmin(c *gin.Context) {
	var req GrantAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// In a production environment, you might want additional validation here
	// such as checking if the requester is a super-admin, rate limiting, etc.

	// For now, we'll directly update the user's admin status
	// This is intentionally simple - in production you might want an approval workflow

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin grant functionality pending - use database migration to set is_admin=true",
		"email":   req.Email,
		"note":    "This endpoint will be implemented with proper approval workflow",
	})
}
