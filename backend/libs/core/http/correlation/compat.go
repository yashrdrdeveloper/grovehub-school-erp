// Package correlation is a compatibility facade.
// Deprecated: use the root core/http package for new code.
package correlation

import ("context"; "net/http"; corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http")
const Header = corehttp.RequestID
func WithNew(ctx context.Context) context.Context { return corehttp.WithNew(ctx) }
func ID(ctx context.Context) string { return corehttp.ID(ctx) }
func Ensure(ctx context.Context) context.Context { return corehttp.Ensure(ctx) }
func Ingress(next http.Handler) http.Handler { return corehttp.Ingress(next) }
