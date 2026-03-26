package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skillora/backend/internal/api"
	"github.com/skillora/backend/internal/config"
)

func main() {
	// Load configuration from environment.
	config.Load()

	// Create Gin engine with default middleware (Logger + Recovery).
	router := gin.Default()

	// Apply global CORS middleware.
	router.Use(api.CORSMiddleware())

	// Health check endpoint.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "skillora-backend"})
	})

	// API v1 route group (populated by later phases).
	v1 := router.Group("/api/v1")
	_ = v1

	// Build HTTP server.
	srv := &http.Server{
		Addr:         ":" + config.C.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background goroutine.
	go func() {
		log.Printf("[server] listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[server] fatal: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT / SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[server] shutting down gracefully …")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[server] forced shutdown: %v", err)
	}

	log.Println("[server] stopped")
}
