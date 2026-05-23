package middleware

import (
	"net/http"
	"CesarRodriguezPardo/template-go/infra/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultRate  = 100
	defaultBurst = 150
)

// bucket represents a token bucket for rate limiting.
type bucket struct {
	tokens    float64
	lastCheck time.Time
}

// RateLimiter is an in-memory token bucket rate limiter per IP.
type RateLimiter struct {
	rate   float64 // tokens added per second
	burst  float64 // maximum bucket size
	buckets map[string]*bucket
	mu      sync.RWMutex
}

// NewRateLimiter creates a new token bucket rate limiter.
func NewRateLimiter(requestsPerMinute, burst int) *RateLimiter {
	if requestsPerMinute <= 0 {
		requestsPerMinute = defaultRate
	}
	if burst <= 0 {
		burst = defaultBurst
	}

	return &RateLimiter{
		rate:    float64(requestsPerMinute) / 60.0,
		burst:   float64(burst),
		buckets: make(map[string]*bucket),
	}
}

// Allow checks if a request from the given key (IP) is allowed.
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[key]
	if !exists {
		b = &bucket{tokens: rl.burst - 1, lastCheck: now}
		rl.buckets[key] = b
		return true
	}

	elapsed := now.Sub(b.lastCheck).Seconds()
	b.tokens += elapsed * rl.rate
	if b.tokens > rl.burst {
		b.tokens = rl.burst
	}
	b.lastCheck = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}

	return false
}

// RateLimitMiddleware returns a gin middleware that applies rate limiting per client IP.
func RateLimitMiddleware(requestsPerMinute, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(requestsPerMinute, burst)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		if !limiter.Allow(clientIP) {
			response.JsonResponse(c, http.StatusTooManyRequests, "rate limit exceeded", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
