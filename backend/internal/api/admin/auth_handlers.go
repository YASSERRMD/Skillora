package admin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/crypto"
	"github.com/skillora/backend/internal/password"
	"github.com/skillora/backend/internal/repository"
)

// AdminAuthHandler handles admin authentication operations
type AdminAuthHandler struct {
	adminCredRepo *repository.AdminCredentialsRepository
	userRepo      *repository.UserRepository
}

// NewAdminAuthHandler creates a new admin auth handler
func NewAdminAuthHandler(
	adminCredRepo *repository.AdminCredentialsRepository,
	userRepo *repository.UserRepository,
) *AdminAuthHandler {
	return &AdminAuthHandler{
		adminCredRepo: adminCredRepo,
		userRepo:      userRepo,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  string `json:"user_id,omitempty"`
}

// Login handles admin username/password login
// POST /api/v1/admin/login
func (h *AdminAuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx := context.Background()

	// Fetch admin credentials
	cred, err := h.adminCredRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// Check if credential is active
	if !cred.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin account is disabled"})
		return
	}

	// Verify password
	if err := password.VerifyPassword(req.Password, cred.PasswordHash); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// Fetch user to verify admin status
	user, err := h.userRepo.GetByID(ctx, cred.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	// Verify user is actually an admin
	if !user.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "user does not have admin privileges"})
		return
	}

	// Update last login timestamp
	if err := h.adminCredRepo.UpdateLastLogin(ctx, cred.ID); err != nil {
		// Log error but don't fail the login
		// In production, you might want to log this to a monitoring system
	}

	// Generate JWT token
	jwtToken, err := crypto.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Set JWT as HttpOnly cookie
	c.SetCookie(
		"skillora_token",
		jwtToken,
		int((7 * 24 * 3600)), // 7 days
		"/",
		"",
		false, // Secure - set to true in production with HTTPS
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, LoginResponse{
		Success: true,
		Message: "Login successful",
		UserID:  user.ID,
	})
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// ChangePassword allows an admin to change their password
// POST /api/v1/admin/change-password
func (h *AdminAuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate new password
	if err := password.ValidatePassword(req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	userID := c.GetString("user_id")

	// Fetch admin credentials for this user
	cred, err := h.adminCredRepo.GetByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin credentials not found"})
		return
	}

	// Verify current password
	if err := password.VerifyPassword(req.CurrentPassword, cred.PasswordHash); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}

	// Hash new password
	newHash, err := password.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Update password in database
	if err := h.adminCredRepo.UpdatePassword(ctx, cred.ID, newHash); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed successfully",
	})
}
