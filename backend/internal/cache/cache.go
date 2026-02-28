package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nexus-id/backend/internal/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

var (
	ErrCacheMiss         = errors.New("cache: key not found")
	ErrNullCache         = errors.New("cache: null cache placeholder")
	NullCachePlaceholder = "__NULL__"
)

// Cache provides high-availability caching with singleflight and NULL caching
type Cache struct {
	client  *redis.Client
	group   *singleflight.Group
	logger  *zap.Logger
	nullTTL time.Duration
}

// New creates a new Cache instance
func New(cfg *config.Config, logger *zap.Logger) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConn,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("Redis connected successfully",
		zap.String("addr", fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)),
	)

	return &Cache{
		client:  client,
		group:   &singleflight.Group{},
		logger:  logger,
		nullTTL: 5 * time.Minute, // NULL cache expires in 5 minutes
	}, nil
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	// Check for NULL cache placeholder
	if val == NullCachePlaceholder {
		return ErrNullCache
	}

	return json.Unmarshal([]byte(val), dest)
}

// Set stores a value in cache
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Delete removes a key from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// DeletePattern removes keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}

// GetOrLoad retrieves a value from cache or loads it using the provided function
// This implements singleflight to prevent cache stampede and NULL caching for cache penetration
func (c *Cache) GetOrLoad(ctx context.Context, key string, dest interface{}, ttl time.Duration, loadFn func() (interface{}, error)) error {
	// Try to get from cache first
	err := c.Get(ctx, key, dest)
	if err == nil {
		return nil // Cache hit
	}

	if err == ErrNullCache {
		// NULL cache hit - return error to prevent DB query
		return ErrCacheMiss
	}

	if err != ErrCacheMiss {
		// Unexpected error
		return err
	}

	// Cache miss - use singleflight to prevent stampede
	result, err, shared := c.group.Do(key, func() (interface{}, error) {
		// Double-check cache in case another goroutine already populated it
		if err := c.Get(ctx, key, dest); err == nil {
			return nil, nil // Another goroutine cached it
		}

		// Load from data source
		value, err := loadFn()
		if err != nil {
			return nil, err
		}

		// Cache the value
		if value == nil {
			// NULL caching - store placeholder to prevent repeated queries for non-existent data
			if err := c.Set(ctx, key, NullCachePlaceholder, c.nullTTL); err != nil {
				c.logger.Warn("Failed to set NULL cache", zap.Error(err), zap.String("key", key))
			}
			return nil, ErrCacheMiss
		}

		// Marshal and cache the value
		data, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
			c.logger.Warn("Failed to set cache", zap.Error(err), zap.String("key", key))
		}

		return data, nil
	})

	if shared {
		c.logger.Debug("Shared singleflight result", zap.String("key", key))
	}

	if err != nil {
		return err
	}

	if result == nil {
		return ErrCacheMiss
	}

	// Unmarshal the result
	return json.Unmarshal(result.([]byte), dest)
}

// Invalidate removes a key from cache and singleflight group
func (c *Cache) Invalidate(ctx context.Context, key string) error {
	// Remove from singleflight group to allow reload
	c.group.Forget(key)
	return c.Delete(ctx, key)
}

// InvalidatePattern removes keys matching a pattern
func (c *Cache) InvalidatePattern(ctx context.Context, pattern string) error {
	// Note: singleflight group doesn't support pattern-based invalidation
	// Keys will naturally expire or be forgotten on next access
	return c.DeletePattern(ctx, pattern)
}

// Exists checks if a key exists in cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	return count > 0, err
}

// Close closes the Redis connection
func (c *Cache) Close() error {
	return c.client.Close()
}

// GetClient returns the underlying Redis client for advanced operations
func (c *Cache) GetClient() *redis.Client {
	return c.client
}
