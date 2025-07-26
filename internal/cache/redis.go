package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/redis/go-redis/v9"
)

// RedisClient wraps the Redis client with caching functionality
type RedisClient struct {
    client *redis.Client
    ctx    context.Context
}

// NewRedisClient creates a new Redis client using REDIS_URL environment variable
func NewRedisClient() *RedisClient {
    redisURL := os.Getenv("REDIS_URL")
    if redisURL == "" {
        redisURL = "localhost:6379"
    }

    rdb := redis.NewClient(&redis.Options{
        Addr: redisURL,
    })

    return &RedisClient{
        client: rdb,
        ctx:    context.Background(),
    }
}

// Get retrieves a value from Redis and unmarshals it into the provided interface
func (r *RedisClient) Get(key string, dest interface{}) error {
    val, err := r.client.Get(r.ctx, key).Result()
    if err != nil {
        return err
    }

    return json.Unmarshal([]byte(val), dest)
}

// Set stores a value in Redis with default expiration
func (r *RedisClient) Set(key string, value interface{}) error {
    return r.SetWithTTL(key, value, 0)
}

// SetWithTTL stores a value in Redis with specified TTL
func (r *RedisClient) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("failed to marshal value: %w", err)
    }

    return r.client.Set(r.ctx, key, data, ttl).Err()
}

// Incr increments a counter and returns the new value
func (r *RedisClient) Incr(key string) (int64, error) {
    return r.client.Incr(r.ctx, key).Result()
}

// Expire sets TTL on an existing key
func (r *RedisClient) Expire(key string, ttl time.Duration) error {
    return r.client.Expire(r.ctx, key, ttl).Err()
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
    return r.client.Close()
}