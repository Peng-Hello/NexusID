package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nexus-id/backend/internal/cache"
	"go.uber.org/zap"
)

const (
	// MaxFailedAttempts is the number of failed attempts before lockout
	MaxFailedAttempts = 5
	// LockoutDuration is how long an account remains locked
	LockoutDuration = 30 * time.Minute
	// RateLimitWindow is the time window for rate limiting
	RateLimitWindow = 15 * time.Minute
)

// RateLimiter handles rate limiting for login attempts
type RateLimiter struct {
	cache    *cache.Cache
	logger   *zap.Logger
	mu       sync.Mutex
	attempts map[string]*AttemptInfo
}

// AttemptInfo tracks login attempts
type AttemptInfo struct {
	Count        int
	FirstAttempt time.Time
	LastAttempt  time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cache *cache.Cache, logger *zap.Logger) *RateLimiter {
	rl := &RateLimiter{
		cache:    cache,
		logger:   logger,
		attempts: make(map[string]*AttemptInfo),
	}

	// Clean up old entries periodically
	go rl.cleanup()

	return rl
}

// CheckRateLimit checks if the request should be rate limited
func (rl *RateLimiter) CheckRateLimit(key string) (bool, time.Duration, error) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	info, exists := rl.attempts[key]

	// Check cache for existing lockout
	cacheKey := fmt.Sprintf("ratelimit:%s", key)
	var lockedUntil time.Time
	if err := rl.cache.Get(context.Background(), cacheKey, &lockedUntil); err == nil {
		if now.Before(lockedUntil) {
			return true, lockedUntil.Sub(now), nil
		}
		// Lockout expired, remove from cache
		rl.cache.Invalidate(context.Background(), cacheKey)
	}

	if !exists {
		rl.attempts[key] = &AttemptInfo{
			Count:        1,
			FirstAttempt: now,
			LastAttempt:  now,
		}
		return false, 0, nil
	}

	// Reset count if window has expired
	if now.Sub(info.FirstAttempt) > RateLimitWindow {
		info.Count = 1
		info.FirstAttempt = now
		info.LastAttempt = now
		return false, 0, nil
	}

	// Increment count
	info.Count++
	info.LastAttempt = now

	// Check if limit exceeded
	if info.Count >= MaxFailedAttempts {
		// Lock the account
		lockUntil := now.Add(LockoutDuration)

		// Store in cache
		if err := rl.cache.Set(context.Background(), cacheKey, lockUntil, LockoutDuration); err != nil {
			rl.logger.Warn("Failed to cache lockout", zap.Error(err))
		}

		rl.logger.Warn("Rate limit exceeded",
			zap.String("key", key),
			zap.Int("attempts", info.Count),
			zap.Duration("lockout", LockoutDuration),
		)

		return true, LockoutDuration, nil
	}

	remainingAttempts := MaxFailedAttempts - info.Count
	retryAfter := RateLimitWindow - now.Sub(info.FirstAttempt)

	rl.logger.Debug("Login attempt",
		zap.String("key", key),
		zap.Int("count", info.Count),
		zap.Int("remaining", remainingAttempts),
	)

	// Return true if we should warn the user
	return false, retryAfter, nil
}

// Reset resets the rate limit counter for a key
func (rl *RateLimiter) Reset(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.attempts, key)

	cacheKey := fmt.Sprintf("ratelimit:%s", key)
	rl.cache.Invalidate(context.Background(), cacheKey)
}

// GetRemainingAttempts returns the number of remaining attempts
func (rl *RateLimiter) GetRemainingAttempts(key string) int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	info, exists := rl.attempts[key]
	if !exists {
		return MaxFailedAttempts
	}

	// Check if window has expired
	if time.Since(info.FirstAttempt) > RateLimitWindow {
		return MaxFailedAttempts
	}

	remaining := MaxFailedAttempts - info.Count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// cleanup periodically removes old entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, info := range rl.attempts {
			// Remove entries older than the lockout duration
			if now.Sub(info.LastAttempt) > LockoutDuration*2 {
				delete(rl.attempts, key)
			}
		}
		rl.mu.Unlock()
	}
}

// LoginRateLimitMiddleware creates a Gin middleware for rate limiting login attempts
func LoginRateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get IP address
		ip := c.ClientIP()

		// Check rate limit
		limited, retryAfter, err := limiter.CheckRateLimit(ip)
		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		if limited {
			c.JSON(429, gin.H{
				"error":       "Too many failed login attempts. Please try again later.",
				"retry_after": retryAfter.Seconds(),
			})
			c.Abort()
			return
		}

		// Get remaining attempts for warning headers
		remaining := limiter.GetRemainingAttempts(ip)
		if remaining < MaxFailedAttempts-2 {
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", MaxFailedAttempts))
		}

		c.Next()
	}
}
