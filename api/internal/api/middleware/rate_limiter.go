package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterConfig holds configuration for rate limiting
type RateLimiterConfig struct {
	MaxRequests     int
	WindowDuration  time.Duration
	CleanupInterval time.Duration
}

// DefaultAuthRateLimiterConfig returns default config for auth endpoints
func DefaultAuthRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		MaxRequests:     100,
		WindowDuration:  1 * time.Minute,
		CleanupInterval: 5 * time.Minute,
	}
}

// RateLimiter provides IP-based rate limiting
type RateLimiter struct {
	config   RateLimiterConfig
	requests map[string]*requestInfo
	mu       sync.RWMutex
	stopCh   chan struct{}
}

type requestInfo struct {
	count     int
	windowEnd time.Time
}

// NewRateLimiter creates a new rate limiter with the given configuration
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		config:   config,
		requests: make(map[string]*requestInfo),
		stopCh:   make(chan struct{}),
	}

	go rl.cleanupLoop()

	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanup()
		case <-rl.stopCh:
			return
		}
	}
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, info := range rl.requests {
		if now.After(info.windowEnd) {
			delete(rl.requests, ip)
		}
	}
}

// Stop stops the cleanup goroutine
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}

func (rl *RateLimiter) isAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	info, exists := rl.requests[ip]

	if !exists || now.After(info.windowEnd) {
		rl.requests[ip] = &requestInfo{
			count:     1,
			windowEnd: now.Add(rl.config.WindowDuration),
		}
		return true
	}

	if info.count >= rl.config.MaxRequests {
		return false
	}

	info.count++
	return true
}

func (rl *RateLimiter) getRemainingTime(ip string) time.Duration {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	info, exists := rl.requests[ip]
	if !exists {
		return 0
	}

	remaining := time.Until(info.windowEnd)
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Middleware returns a Gin middleware that rate limits requests
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !rl.isAllowed(ip) {
			remaining := rl.getRemainingTime(ip)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":        "RATE_LIMIT_EXCEEDED",
					"message":     "Too many requests. Please try again later.",
					"retry_after": int(remaining.Seconds()),
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthRateLimiter is a pre-configured rate limiter for authentication endpoints
var AuthRateLimiter = NewRateLimiter(DefaultAuthRateLimiterConfig())

// RateLimitAuth returns a middleware that rate limits authentication attempts
func RateLimitAuth() gin.HandlerFunc {
	return AuthRateLimiter.Middleware()
}
