package ratelimit

import "net/http"

type Limiter interface {
	Allow(string, *http.Request) bool
}
