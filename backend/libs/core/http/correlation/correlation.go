// Package correlation manages server-generated request correlation identifiers.
// Correlation is diagnostic metadata, never authentication or authorization.
package correlation

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type key struct{}

const Header = "X-Request-ID"

func WithNew(ctx context.Context) context.Context {
	return context.WithValue(ctx, key{}, uuid.NewString())
}

func ID(ctx context.Context) string {
	id, _ := ctx.Value(key{}).(string)
	return id
}

func Ensure(ctx context.Context) context.Context {
	if ID(ctx) != "" {
		return ctx
	}
	return WithNew(ctx)
}

// Ingress deliberately ignores untrusted client-supplied IDs.
func Ingress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := WithNew(r.Context())
		w.Header().Set(Header, ID(ctx))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
