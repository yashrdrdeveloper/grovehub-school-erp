// Package corehttp is the public facade for generic HTTP capabilities.
package corehttp

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	coreerrors "grovehub/libs/core/errors"
	"grovehub/libs/core/http/internal/src/constants"
	"grovehub/libs/core/http/internal/src/contracts"
	"grovehub/libs/core/http/internal/src/helpers"
	"grovehub/libs/core/http/internal/src/operations/correlation"
	"grovehub/libs/core/http/internal/src/operations/endpoint"
	"grovehub/libs/core/http/internal/src/operations/readiness"
	"grovehub/libs/core/http/internal/src/operations/request"
	"grovehub/libs/core/http/internal/src/operations/response"
	"grovehub/libs/core/http/internal/src/operations/tracing"
	corerequest "grovehub/libs/core/request"
)

const (
	IdempotencyKey = constants.IdempotencyKey
	CorrelationID  = constants.CorrelationID
	RequestID      = constants.RequestID
	TraceID        = constants.TraceID
)

type Success[T any] = contracts.Success[T]
type Error[M any] = contracts.Error[M]
type Request[Path, Query, Body any] = endpoint.Request[Path, Query, Body]
type Check = readiness.Check

func NewSuccess[T any](data T) Success[T] {
	return contracts.NewSuccess(data)
}

func NewError[M any](code coreerrors.Code, message string, status int, metadata M) Error[M] {
	return contracts.NewError(code, message, status, metadata)
}

func WithNew(ctx context.Context) context.Context { return correlation.WithNew(ctx) }
func ID(ctx context.Context) string               { return correlation.ID(ctx) }
func Ensure(ctx context.Context) context.Context  { return correlation.Ensure(ctx) }
func Ingress(next http.Handler) http.Handler      { return correlation.Ingress(next) }

func Call[T any](ctx context.Context, handler http.Handler, method, route string, pathDTO, body any) (T, int, error) {
	return endpoint.Call[T](ctx, handler, method, route, pathDTO, body)
}

func CallRequest[Response, Path, Query, Body any](ctx context.Context, handler http.Handler, method, route string, input Request[Path, Query, Body]) (Response, int, error) {
	return endpoint.CallRequest[Response](ctx, handler, method, route, input)
}

func CallStruct[Response any, Input corerequest.Contract](ctx context.Context, handler http.Handler, method, route string, input Input) (Response, int, error) {
	return endpoint.CallStruct[Response](ctx, handler, method, route, input)
}

func LevelForStatus(status int) slog.Level { return helpers.LevelForStatus(status) }

func ReadinessHandler(draining *atomic.Bool, checks ...Check) http.HandlerFunc {
	return readiness.Handler(draining, checks...)
}

func ReadinessHTTP(origin, path string) (Check, error) {
	return readiness.HTTP(origin, path)
}

func DecodeJSON(req *http.Request, destination any) error {
	return request.DecodeJSON(req, destination)
}
func DecodeStrictJSON(req *http.Request, destination any) error {
	return request.DecodeStrictJSON(req, destination)
}

func WriteSuccess[T any](writer http.ResponseWriter, status int, data T) error {
	return response.WriteSuccess(writer, status, data)
}

func WriteError(writer http.ResponseWriter, status int, definition coreerrors.Definition) error {
	return response.WriteError(writer, status, definition)
}

func WriteErrorContext(ctx context.Context, writer http.ResponseWriter, status int, definition coreerrors.Definition) error {
	return response.WriteErrorContext(ctx, writer, status, definition)
}

func WriteErrorWithMetadata[M any](writer http.ResponseWriter, status int, definition coreerrors.Definition, metadata M) error {
	return response.WriteErrorWithMetadata(writer, status, definition, metadata)
}

func WriteErrorWithMetadataContext[M any](ctx context.Context, writer http.ResponseWriter, status int, definition coreerrors.Definition, metadata M) error {
	return response.WriteErrorWithMetadataContext(ctx, writer, status, definition, metadata)
}

func NewTransport(base http.RoundTripper) http.RoundTripper { return tracing.NewTransport(base) }
func NewServiceTransport(service string, base http.RoundTripper) http.RoundTripper {
	return tracing.NewServiceTransport(service, base)
}
func NewClient(timeout time.Duration) *http.Client { return tracing.NewClient(timeout) }
func NewServiceClient(service string, timeout time.Duration) *http.Client {
	return tracing.NewServiceClient(service, timeout)
}
func ServiceHandler(service, operation string, next http.Handler) http.Handler {
	return tracing.ServiceHandler(service, operation, next)
}
