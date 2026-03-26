package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skillora/backend/internal/config"
)

var PG *pgxpool.Pool

// InitPostgres creates the pgxpool connection, enables pgvector, and verifies connectivity.
func InitPostgres(ctx context.Context) error {
	cfg, err := pgxpool.ParseConfig(config.C.DatabaseURL)
	if err != nil {
		return fmt.Errorf("postgres: parse config: %w", err)
	}

	cfg.MaxConns = 25
	cfg.MinConns = 5
	cfg.MaxConnLifetime = 1 * time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("postgres: connect: %w", err)
	}

	// Verify connection.
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("postgres: ping: %w", err)
	}

	// Enable pgvector extension.
	if _, err := pool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector;"); err != nil {
		return fmt.Errorf("postgres: enable vector extension: %w", err)
	}

	PG = pool
	log.Println("[db] PostgreSQL connected and pgvector extension enabled")
	return nil
}

// ClosePostgres gracefully closes the pool.
func ClosePostgres() {
	if PG != nil {
		PG.Close()
		log.Println("[db] PostgreSQL connection pool closed")
	}
}
