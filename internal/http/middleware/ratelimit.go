package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"codedock.dev/codedock/internal/utils"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for range time.NewTicker(time.Minute).C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-rl.window)
		for ip, v := range rl.visitors {
			if v.lastSeen.Before(cutoff) {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		now := time.Now()

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists || v.lastSeen.Add(rl.window).Before(now) {
			rl.visitors[ip] = &visitor{count: 1, lastSeen: now}
			rl.mu.Unlock()
			return next(c)
		}
		v.count++
		v.lastSeen = now
		if v.count > rl.limit {
			rl.mu.Unlock()
			return c.JSON(http.StatusTooManyRequests, utils.RateLimitError{Message: "rate limit exceeded", RetryAfter: int(rl.window.Seconds())})
		}
		rl.mu.Unlock()
		return next(c)
	}
}
