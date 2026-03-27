package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	DatabaseURL        string
	RedisURL           string
	GoogleClientID     string
	GoogleClientSecret string
	JWTSecret          string
	AESMasterKey       string
	FrontendURL        string
	Port               string
	AdminUsername      string
	AdminPassword      string
}

var C *Config

// Load reads environment variables (and an optional .env file) into the global Config.
func Load() {
	// Load .env if present; ignore error if not found (production uses real env).
	if err := godotenv.Load("../../.env"); err != nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("[config] no .env file found, reading from environment")
		}
	}

	C = &Config{
		DatabaseURL:        mustEnv("DATABASE_URL"),
		RedisURL:           mustEnv("REDIS_URL"),
		GoogleClientID:     mustEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: mustEnv("GOOGLE_CLIENT_SECRET"),
		JWTSecret:          mustEnv("JWT_SECRET"),
		AESMasterKey:       mustEnv("AES_MASTER_KEY"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		Port:               getEnv("PORT", "8080"),
		AdminUsername:      getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword:      getEnv("ADMIN_PASSWORD", "skillora_pass_2026"),
	}
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("[config] required environment variable %q is not set", key)
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
