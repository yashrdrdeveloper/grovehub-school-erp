package impl

import (
	"grovehub/launcher/internal/src/services/gateway/operations/ratelimit"
	"net"
	"net/http"
	"sync"
	"time"
)

// Limiter is a bounded per-process throttle. Production multi-replica
// deployments must add a shared/distributed limiter at the edge.
// @Service
type limiter struct {
	mu      sync.Mutex
	entries map[string]attempt
	limits  map[string]int
	clock   func() time.Time
}

type attempt struct {
	count int
	until time.Time
}

func New(limits map[string]int) ratelimit.Limiter {
	return NewWithClock(limits, time.Now)
}

// NewWithClock constructs the same limiter with an explicit clock.
// The Launcher uses New; deterministic unit tests can advance time without sleeping.
func NewWithClock(limits map[string]int, clock func() time.Time) ratelimit.Limiter {
	configured := make(map[string]int, len(limits))
	for pattern, count := range limits {
		if count > 0 {
			configured[pattern] = count
		}
	}

	return &limiter{entries: make(map[string]attempt), limits: configured, clock: clock}
}

func (limiter *limiter) Allow(pattern string, request *http.Request) bool {
	limit, configured := limiter.limits[request.Method+" "+pattern]
	if !configured || limit < 1 {
		return true
	}

	address, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		address = request.RemoteAddr
	}

	key := pattern + "|" + address
	now := limiter.clock()

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

var _ ratelimit.Limiter = (*limiter)(nil)
