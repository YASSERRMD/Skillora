package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/skillora/backend/internal/config"
)

var RDB *redis.Client

// InitRedis creates and verifies the Redis client connection.
func InitRedis(ctx context.Context) error {
	opt, err := redis.ParseURL(config.C.RedisURL)
	if err != nil {
		return fmt.Errorf("redis: parse URL: %w", err)
	}

	opt.PoolSize = 20
	opt.MinIdleConns = 5
	opt.DialTimeout = 5 * time.Second
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second

	RDB = redis.NewClient(opt)

	if _, err := RDB.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("redis: ping: %w", err)
	}

	log.Println("[db] Redis connected")
	return nil
}

// CloseRedis gracefully closes the Redis client.
func CloseRedis() {
	if RDB != nil {
		if err := RDB.Close(); err != nil {
			log.Printf("[db] redis close error: %v", err)
		}
		log.Println("[db] Redis connection closed")
	}
}

// SetJSON marshals v to JSON and stores it under key with the given TTL.
func SetJSON(ctx context.Context, key string, v any, ttl time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("redis SetJSON marshal: %w", err)
	}
	return RDB.Set(ctx, key, data, ttl).Err()
}

// GetJSON retrieves the JSON value stored under key and unmarshals it into dest.
// Returns redis.Nil if the key does not exist.
func GetJSON(ctx context.Context, key string, dest any) error {
	data, err := RDB.Get(ctx, key).Bytes()
	if err != nil {
		return err // callers check for redis.Nil
	}
	return json.Unmarshal(data, dest)
}
