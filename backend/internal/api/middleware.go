package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/config"
)

// CORSMiddleware allows the Next.js frontend to make credentialed requests.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := config.C.FrontendURL

		// Allow exact frontend origin or localhost variations during development.
		if origin == allowed || strings.HasPrefix(origin, "http://localhost:") {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Idempotency-Key")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// AdminBasicAuth protects admin routes using credentials from environment.
func AdminBasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		config.C.AdminUsername: config.C.AdminPassword,
	})
}
