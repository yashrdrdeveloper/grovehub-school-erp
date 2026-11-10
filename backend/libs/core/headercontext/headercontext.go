// Package headercontext exposes allowlisted, validated inbound HTTP metadata
// through context.Context without trusting arbitrary headers or trace IDs.
package headercontext

import (
	"context"
	"errors"
	"net/http"
	"strings"

	coreevents "grovehub/libs/core/events"
	httpheaders "grovehub/libs/core/http/headers"
	coreidempotency "grovehub/libs/core/idempotency"
)

const maxMetadataIDLength = 128

var ErrInvalidHeader = errors.New("invalid request metadata header")

type contextKey struct{}

// Values is a small immutable copy of trusted header metadata.
type Values struct {
	IdempotencyKey string
	CorrelationID string
	RequestID string
}

// WithValues is reserved for trusted in-process operations and tests.
func WithValues(ctx context.Context, values Values) context.Context {
	ctx = context.WithValue(ctx, contextKey{}, values)
	if values.CorrelationID != "" { ctx = coreevents.WithCorrelationID(ctx, values.CorrelationID) }
	return ctx
}

func FromContext(ctx context.Context) Values {
	if values, ok := ctx.Value(contextKey{}).(Values); ok { return values }
	return Values{}
}

func IdempotencyKey(ctx context.Context) string { return FromContext(ctx).IdempotencyKey }
func CorrelationID(ctx context.Context) string { return FromContext(ctx).CorrelationID }
func RequestID(ctx context.Context) string { return FromContext(ctx).RequestID }

// WithoutIdempotencyKey isolates downstream internal commands from the
// idempotency key belonging to the inbound parent HTTP operation.
func WithoutIdempotencyKey(ctx context.Context) context.Context {
	values := FromContext(ctx)
	values.IdempotencyKey = ""
	return WithValues(ctx, values)
}

// Middleware captures allowlisted metadata once at the HTTP edge.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		values, err := readHTTPHeaders(request.Header)
		if err != nil {
			http.Error(writer, "invalid request metadata header", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(writer, request.WithContext(WithValues(request.Context(), values)))
	})
}

func readHTTPHeaders(header http.Header) (Values, error) {
	idempotency, err := readOne(header, httpheaders.IdempotencyKey)
	if err != nil { return Values{}, err }
	if idempotency != "" {
		idempotency, err = coreidempotency.ValidateKey(idempotency)
		if err != nil { return Values{}, ErrInvalidHeader }
	}
	correlation, err := readOne(header, httpheaders.CorrelationID)
	if err != nil || !validID(correlation) { return Values{}, ErrInvalidHeader }
	requestID, err := readOne(header, httpheaders.RequestID)
	if err != nil || !validID(requestID) { return Values{}, ErrInvalidHeader }
	return Values{IdempotencyKey: idempotency, CorrelationID: correlation, RequestID: requestID}, nil
}

func readOne(header http.Header, name string) (string, error) {
	values := header.Values(name)
	if len(values) == 0 { return "", nil }
	if len(values) != 1 { return "", ErrInvalidHeader }
	value := strings.TrimSpace(values[0])
	if value == "" { return "", ErrInvalidHeader }
	return value, nil
}

func validID(value string) bool {
	if len(value) > maxMetadataIDLength { return false }
	for _, c := range value {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' { continue }
		switch c {
		case '-', '_', '.', ':', '/': continue
		default: return false
		}
	}
	return true
}
