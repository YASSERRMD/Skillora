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
	"github.com/skillora/backend/internal/agents"
	"github.com/skillora/backend/internal/api"
	adminapi "github.com/skillora/backend/internal/api/admin"
	barterapi "github.com/skillora/backend/internal/api/barter"
	userapi "github.com/skillora/backend/internal/api/user"
	skillsapi "github.com/skillora/backend/internal/api/skills"
	matchingapi "github.com/skillora/backend/internal/api/matching"
	"github.com/skillora/backend/internal/auth"
	"github.com/skillora/backend/internal/config"
	"github.com/skillora/backend/internal/db"
	"github.com/skillora/backend/internal/llm"
	"github.com/skillora/backend/internal/repository"
	"github.com/skillora/backend/internal/ws"
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

	// Initialize context for setup
	ctxSetup := context.Background()

	// Initialize Database and Cache
	if err := db.InitPostgres(ctxSetup); err != nil {
		log.Fatalf("failed to init postgres: %v", err)
	}
	defer db.ClosePostgres()

	if err := db.InitRedis(ctxSetup); err != nil {
		log.Fatalf("failed to init redis: %v", err)
	}

	// Wait for PG Vector readiness in real environment...
	time.Sleep(1 * time.Second) // Tiny sleep ensuring connection pool settles

	// Repositories
	userRepo := repository.NewUserRepository(db.PG)
	skillRepo := repository.NewSkillRepository(db.PG)
	userSkillRepo := repository.NewUserSkillRepository(db.PG)
	llmRepo := repository.NewLLMRepository(db.PG)

	// LLM Pipeline & Orchestration
	llmManager := llm.NewManager(llmRepo)
	llmManager.StartBackgroundSync(ctxSetup)
	llmRouter := llm.NewRouter(llmManager)

	// Agents
	appraisalAgent := agents.NewAppraisalAgent(llmRouter)
	vectorRepo := repository.NewVectorRepository(db.PG)
	embeddingAgent := agents.NewEmbeddingAgent(llmRouter, vectorRepo)

	barterRepo := repository.NewBarterRepository(db.PG)
	milestoneAgent := agents.NewMilestoneAgent(llmRouter)
	notifRepo := repository.NewNotificationRepository(db.PG)
	notifHub := api.NewNotificationHub(notifRepo)
	wsHub := ws.NewHub()

	// Route Handlers
	oauthCfg := auth.NewGoogleOAuthConfig()
	authHandler := auth.NewHandler(oauthCfg, userRepo)
	userHandler := userapi.NewHandler(userRepo, userSkillRepo)
	adminHandler := adminapi.NewLLMHandler(llmRepo)
	skillsHandler := skillsapi.NewHandler(skillRepo, userSkillRepo, appraisalAgent)
	barterHandler := barterapi.NewHandler(barterRepo, milestoneAgent)
	matchingHandler := matchingapi.NewHandler(vectorRepo, embeddingAgent)

	// --- Routes Setup ---
	v1 := router.Group("/api/v1")
	{
		// Public Auth
		authGrp := v1.Group("/auth")
		{
			authGrp.GET("/google/login", authHandler.GoogleLogin)
			authGrp.GET("/google/callback", authHandler.GoogleCallback)
		}

		// Public Skills/Categories
		v1.GET("/categories", skillsHandler.GetCategories)
		v1.GET("/categories/:id/skills", skillsHandler.GetCategorySkills)

		// Protected User Routes
		userGrp := v1.Group("/users")
		userGrp.Use(api.RequireAuth())
		{
			userGrp.GET("/me", userHandler.GetMe)
			userGrp.PUT("/me", userHandler.UpdateMe)
			userGrp.GET("/skills", userHandler.GetMySkills)
			userGrp.GET("/:id/skills", userHandler.GetUserSkills)
		}

		// Protected Skill Appraisal Route
		skillGrp := v1.Group("/skills")
		skillGrp.Use(api.RequireAuth())
		{
			skillGrp.POST("/appraise", skillsHandler.PostAppraise)
		}

		// Matching Engine Route
		v1.GET("/match", api.RequireAuth(), matchingHandler.GetMatches)

		// Barter Economy Routes
		barterGrp := v1.Group("/barters")
		barterGrp.Use(api.RequireAuth())
		{
			barterGrp.POST("", barterHandler.PostPropose)
			barterGrp.GET("", barterHandler.GetMyBarters)
			barterGrp.GET("/balance", barterHandler.GetCreditBalance)
			barterGrp.PATCH("/:id/status", barterHandler.PatchBarterStatus)
			barterGrp.POST("/:id/complete", barterHandler.PostComplete)
			barterGrp.GET("/:id/milestones", barterHandler.GetMilestones)
		}

		// Milestone specific actions
		msGrp := v1.Group("/milestones")
		msGrp.Use(api.RequireAuth())
		{
			msGrp.POST("/:id/complete", barterHandler.PostMilestoneComplete)
			msGrp.POST("/:id/approve", barterHandler.PostMilestoneApprove)
		}

		// Internal Admin Routes (Basic Auth protected for internal management)
		adminGrp := v1.Group("/admin")
		adminGrp.Use(api.AdminBasicAuth())
		{
			adminGrp.GET("/llm-providers", adminHandler.GetLLMProviders)
			adminGrp.POST("/llm-providers", adminHandler.PostLLMProvider)
		}
		// Notifications
		notifGrp := v1.Group("/notifications")
		notifGrp.Use(api.RequireAuth())
		{
			notifGrp.GET("", api.GetNotificationsHandler(notifRepo))
			notifGrp.POST("/read", api.MarkNotificationsReadHandler(notifRepo))
			notifGrp.GET("/stream", notifHub.SSEHandler)
		}

		// Real-time WebSocket Hub
		v1.GET("/ws", api.RequireAuth(), api.WSHandler(wsHub))
	}

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
