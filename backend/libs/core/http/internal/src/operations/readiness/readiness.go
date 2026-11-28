// Package readiness provides bounded, dependency-neutral HTTP health checks.
package readiness

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

type Check func(context.Context) error

// Handler returns only aggregate status and never discloses dependency secrets.
func Handler(draining *atomic.Bool, checks ...Check) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		if draining == nil || draining.Load() || len(checks) == 0 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		for _, check := range checks {
			if check == nil || check(ctx) != nil || ctx.Err() != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
		}
		if draining.Load() {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// HTTP checks a fixed configured service origin, without credentials,
// environment proxies, redirects or user-controlled request parameters.
func HTTP(origin, path string) (Check, error) {
	u, err := url.Parse(origin)
	if err != nil || u.Hostname() == "" || u.User != nil || u.ForceQuery || u.RawQuery != "" || u.Fragment != "" ||
		(u.Path != "" && u.Path != "/") || u.RawPath != "" ||
		(u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "127.0.0.1" || u.Hostname() == "::1"))) {
		return nil, errors.New("invalid dependency origin")
	}
	if path != "/ready" {
		return nil, errors.New("invalid dependency readiness path")
	}
	u.Path = path
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	client := &http.Client{Transport: transport, Timeout: 2*time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil { return errors.New("dependency probe failed") }
		res, err := client.Do(req)
		if err != nil { return errors.New("dependency unavailable") }
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
			return errors.New("dependency not ready")
		}
		return nil
	}, nil
}
