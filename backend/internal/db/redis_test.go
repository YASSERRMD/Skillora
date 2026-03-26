package db_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/skillora/backend/internal/config"
	"github.com/skillora/backend/internal/db"
)

func init() {
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("GOOGLE_CLIENT_ID", "test-client-id")
	os.Setenv("GOOGLE_CLIENT_SECRET", "test-client-secret")
	os.Setenv("JWT_SECRET", "test-jwt-secret")
	os.Setenv("AES_MASTER_KEY", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	config.Load()
}

// startMiniRedis spins up an in-process Redis and wires db.RDB to it.
func startMiniRedis(t *testing.T) *miniredis.Miniredis {
	t.Helper()
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	t.Cleanup(mr.Close)

	db.RDB = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return mr
}

type payload struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

func TestSetJSON(t *testing.T) {
	startMiniRedis(t)

	p := payload{Name: "Alice", Score: 42}
	if err := db.SetJSON(context.Background(), "user:1", p, 5*time.Minute); err != nil {
		t.Fatalf("SetJSON error: %v", err)
	}

	var got payload
	if err := db.GetJSON(context.Background(), "user:1", &got); err != nil {
		t.Fatalf("GetJSON after SetJSON error: %v", err)
	}
	if got.Name != p.Name || got.Score != p.Score {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, p)
	}
}

func TestGetJSON(t *testing.T) {
	mr := startMiniRedis(t)

	p := payload{Name: "Bob", Score: 99}
	data, _ := json.Marshal(p)
	if err := mr.Set("user:2", string(data)); err != nil {
		t.Fatal(err)
	}

	var got payload
	if err := db.GetJSON(context.Background(), "user:2", &got); err != nil {
		t.Fatalf("GetJSON error: %v", err)
	}
	if got.Name != p.Name || got.Score != p.Score {
		t.Errorf("got %+v, want %+v", got, p)
	}
}

func TestGetJSON_MissingKey(t *testing.T) {
	startMiniRedis(t)

	var dest payload
	err := db.GetJSON(context.Background(), "missing:key", &dest)
	if err != redis.Nil {
		t.Errorf("expected redis.Nil for missing key, got: %v", err)
	}
}

func TestSetJSON_TTLExpiry(t *testing.T) {
	mr := startMiniRedis(t)

	p := payload{Name: "Expiry", Score: 1}
	if err := db.SetJSON(context.Background(), "expire:1", p, 1*time.Second); err != nil {
		t.Fatalf("SetJSON error: %v", err)
	}

	// Fast-forward miniredis clock past TTL.
	mr.FastForward(2 * time.Second)

	var dest payload
	err := db.GetJSON(context.Background(), "expire:1", &dest)
	if err != redis.Nil {
		t.Errorf("expected redis.Nil after TTL expiry, got: %v", err)
	}
}
