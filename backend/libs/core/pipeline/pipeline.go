// Package pipeline exposes typed REST execution and service orchestration.
package pipeline

import (
	"context"
	"net/http"

	coreannotation "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation"
	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	corerequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/internal/src/operations/rest"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/internal/src/operations/service"
)

type Empty = rest.Empty
type Input[P, Q, B any] = rest.Input[P, Q, B]
type Result[T any] = rest.Result[T]
type Config[P, Q, B, R any] = rest.Config[P, Q, B, R]
type StructConfig[Q, R any] = rest.StructConfig[Q, R]
type Writer[T any] = rest.Writer[T]
type Failure = rest.Failure
type BindError = rest.BindError
type DecodeError = rest.DecodeError

func OK[T any](data T) Result[T] { return rest.OK(data) }
func Created[T any](data T) Result[T] { return rest.Created(data) }
func CreatedOrOK[T any](data T, created bool) Result[T] { return rest.CreatedOrOK(data, created) }
func ReadBody[T any](request *http.Request, decode func(*http.Request, any) error, rules coreannotation.ValidationRules[T]) (T, error) {
	return rest.ReadBody(request, decode, rules)
}

func RunREST[P, Q, B, R any](writer http.ResponseWriter, request *http.Request, config Config[P, Q, B, R]) {
	rest.Run(writer, request, config)
}

func RunStruct[Q corerequest.Contract, R any](writer http.ResponseWriter, request *http.Request, config StructConfig[Q, R]) {
	rest.RunStruct(writer, request, config)
}

type Outcome[R any] = service.Outcome[R]
type Handler[S, R any] = service.Handler[S, R]
type Strategy[S, R any] = service.Strategy[S, R]
type Rule[S, R any] = service.Rule[S, R]
type Operation = service.Operation
type Recovery[S, R any] = service.Recovery[S, R]

var (
	ErrUnresolvedChain = service.ErrUnresolvedChain
	ErrNoMatchingRule = service.ErrNoMatchingRule
	ErrMissingUnitOfWork = service.ErrMissingUnitOfWork
)

func Next[R any]() Outcome[R] { return service.Next[R]() }
func Return[R any](value R, created bool) Outcome[R] { return service.Return(value, created) }
func Fail[R any](err error) Outcome[R] { return service.Fail[R](err) }
func Continue[S, R any](ctx context.Context, state *S) Outcome[R] { return service.Continue[S, R](ctx, state) }
func Chain[S, R any](ctx context.Context, state *S, handlers ...Handler[S, R]) Outcome[R] { return service.Chain(ctx, state, handlers...) }
func RunService[S, R any](ctx context.Context, state *S, handlers ...Handler[S, R]) (R, bool, error) {
	return service.Run(ctx, state, handlers...)
}
func RequireHandled[R any](outcome Outcome[R], err error) Outcome[R] { return service.RequireHandled(outcome, err) }
func When[S, R any](predicate func(*S) bool, handler Handler[S, R]) Rule[S, R] { return service.When(predicate, handler) }
func Otherwise[S, R any](handler Handler[S, R]) Rule[S, R] { return service.Otherwise(handler) }
func Decide[S, R any](ctx context.Context, state *S, rules ...Rule[S, R]) Outcome[R] { return service.Decide(ctx, state, rules...) }
func ResolveResult[T, R any](value T, err error, success func(T) Outcome[R], failure func(error) Outcome[R]) Outcome[R] {
	return service.Result(value, err, success, failure)
}
func Lookup[T, R any](value T, err error, missing error, found func(T) Outcome[R], notFound func() Outcome[R], failed func(error) Outcome[R]) Outcome[R] {
	return service.Lookup(value, err, missing, found, notFound, failed)
}
func Sequence[R any](ctx context.Context, result R, created bool, operations ...Operation) Outcome[R] {
	return service.Sequence(ctx, result, created, operations...)
}
func Transactional[S, R any](unit coredatabase.UnitOfWork, handler Handler[S, R]) Handler[S, R] { return service.Transactional(unit, handler) }
func OnError[S, R any](target error, handler func(context.Context, *S, error) Outcome[R]) Recovery[S, R] {
	return service.OnError(target, handler)
}
func OtherwiseError[S, R any](handler func(context.Context, *S, error) Outcome[R]) Recovery[S, R] { return service.OtherwiseError(handler) }
func Recover[S, R any](ctx context.Context, state *S, result Outcome[R], recoveries ...Recovery[S, R]) Outcome[R] {
	return service.Recover(ctx, state, result, recoveries...)
}
