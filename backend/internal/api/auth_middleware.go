package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/crypto"
	"github.com/skillora/backend/internal/repository"
)

const UserIDKey = "user_id"

// RequireAuth validates the JWT cookie and sets user_id in the Gin context.
// Returns 401 if the cookie is missing or the token is invalid/expired.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("skillora_token")
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		claims, err := crypto.ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

// GetUserID retrieves the authenticated user's ID from the Gin context.
// Panics if called outside of RequireAuth — by design.
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get(UserIDKey)
	return userID.(string)
}

// RequireAdmin validates that the authenticated user has admin privileges.
// Returns 401 if not authenticated, 403 if authenticated but not an admin.
func RequireAdmin(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure user is authenticated
		tokenStr, err := c.Cookie("skillora_token")
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		claims, err := crypto.ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Fetch user from database to check admin status
		ctx := context.Background()
		user, err := userRepo.GetByID(ctx, claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to verify admin status"})
			return
		}

		// Check if user has admin privileges
		if !user.IsAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin privileges required"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}
