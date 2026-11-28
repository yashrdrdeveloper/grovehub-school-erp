// Package headercontext exposes allowlisted, validated inbound HTTP metadata
// through context.Context without trusting arbitrary headers or trace IDs.
package headercontext

import (
	"context"
	"net/http"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/operations/metadata"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/operations/middleware"
)

// Values is a small immutable copy of trusted header metadata.
type Values = contracts.Values

var ErrInvalidHeader = contracts.ErrInvalidHeader

// WithValues is reserved for trusted in-process operations and tests.
func WithValues(ctx context.Context, values Values) context.Context {
	return metadata.WithValues(ctx, values)
}

func FromContext(ctx context.Context) Values {
	return metadata.FromContext(ctx)
}

func IdempotencyKey(ctx context.Context) string {
	return metadata.FromContext(ctx).IdempotencyKey
}

func CorrelationID(ctx context.Context) string {
	return metadata.FromContext(ctx).CorrelationID
}

func RequestID(ctx context.Context) string {
	return metadata.FromContext(ctx).RequestID
}

// WithoutIdempotencyKey isolates downstream commands from the parent HTTP operation.
func WithoutIdempotencyKey(ctx context.Context) context.Context {
	return metadata.WithoutIdempotencyKey(ctx)
}

// Middleware captures allowlisted metadata once at the HTTP edge.
func Middleware(next http.Handler) http.Handler {
	return middleware.New(next)
}
