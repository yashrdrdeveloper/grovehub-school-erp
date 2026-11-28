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
	limits map[string]int
}

type attempt struct {
	count int
	until time.Time
}

func NewLimiter(limits map[string]int) *Limiter {
	configured := make(map[string]int, len(limits))
	for pattern, count := range limits {
		if count > 0 {
			configured[pattern] = count
		}
	}
	return &Limiter{entries: make(map[string]attempt), limits: configured}
}

func (limiter *Limiter) Allow(pattern string, request *http.Request) bool {
	limit, configured := limiter.limits[request.Method+" "+pattern]
	if !configured || limit < 1 {
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
