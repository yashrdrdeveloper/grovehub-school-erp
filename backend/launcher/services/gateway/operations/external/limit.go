package external

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// Limiter is a bounded per-process throttle. Production multi-replica
// deployments must add a shared/distributed limiter at the edge.
type Limiter struct {
	mu sync.Mutex
	entries map[string]attempt
}

type attempt struct {
	count int
	until time.Time
}

func NewLimiter() *Limiter {
	return &Limiter{entries: make(map[string]attempt)}
}

func (limiter *Limiter) Allow(pattern string, request *http.Request) bool {
	limit := 0
	switch pattern {
	case "/identity/public/auth/login":
		limit = 12
	case "/identity/public/auth/forgot-password":
		limit = 5
	case "/identity/public/auth/recovery/exchange",
		"/identity/public/auth/invitations/exchange":
		limit = 8
	case "/identity/public/auth/refresh":
		limit = 30
	default:
		return true
	}

	address, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		address = request.RemoteAddr
	}
	key := pattern + "|" + address
	now := time.Now()

	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	if len(limiter.entries) > 100000 {
		for existingKey, entry := range limiter.entries {
			if now.After(entry.until) {
				delete(limiter.entries, existingKey)
			}
		}
	}
	if len(limiter.entries) > 100000 {
		return false
	}

	entry := limiter.entries[key]
	if now.After(entry.until) {
		entry = attempt{until: now.Add(time.Minute)}
	}
	if entry.count >= limit {
		return false
	}

	entry.count++
	limiter.entries[key] = entry
	return true
}
