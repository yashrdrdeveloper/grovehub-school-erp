// Package service is a temporary compatibility facade. Use core/pipeline.
package service

import (
	"context"
	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
)

type Outcome[R any] = corepipeline.Outcome[R]
type Handler[S,R any] = corepipeline.Handler[S,R]
type Strategy[S,R any] = corepipeline.Strategy[S,R]
type Rule[S,R any] = corepipeline.Rule[S,R]
type Operation = corepipeline.Operation
type Recovery[S,R any] = corepipeline.Recovery[S,R]
var (
	ErrUnresolvedChain=corepipeline.ErrUnresolvedChain
	ErrNoMatchingRule=corepipeline.ErrNoMatchingRule
	ErrMissingUnitOfWork=corepipeline.ErrMissingUnitOfWork
)
func Next[R any]() Outcome[R]{return corepipeline.Next[R]()}
func Return[R any](value R,created bool)Outcome[R]{return corepipeline.Return(value,created)}
func Fail[R any](err error)Outcome[R]{return corepipeline.Fail[R](err)}
func Continue[S,R any](ctx context.Context,state *S)Outcome[R]{return corepipeline.Continue[S,R](ctx,state)}
func Chain[S,R any](ctx context.Context,state *S,handlers ...Handler[S,R])Outcome[R]{return corepipeline.Chain(ctx,state,handlers...)}
func Run[S,R any](ctx context.Context,state *S,handlers ...Handler[S,R])(R,bool,error){return corepipeline.RunService(ctx,state,handlers...)}
func RequireHandled[R any](outcome Outcome[R],err error)Outcome[R]{return corepipeline.RequireHandled(outcome,err)}
func When[S,R any](predicate func(*S)bool,handler Handler[S,R])Rule[S,R]{return corepipeline.When(predicate,handler)}
func Otherwise[S,R any](handler Handler[S,R])Rule[S,R]{return corepipeline.Otherwise(handler)}
func Decide[S,R any](ctx context.Context,state *S,rules ...Rule[S,R])Outcome[R]{return corepipeline.Decide(ctx,state,rules...)}
func Result[T,R any](value T,err error,success func(T)Outcome[R],failure func(error)Outcome[R])Outcome[R]{return corepipeline.ResolveResult(value,err,success,failure)}
func Lookup[T,R any](value T,err error,missing error,found func(T)Outcome[R],notFound func()Outcome[R],failed func(error)Outcome[R])Outcome[R]{return corepipeline.Lookup(value,err,missing,found,notFound,failed)}
func Sequence[R any](ctx context.Context,result R,created bool,operations ...Operation)Outcome[R]{return corepipeline.Sequence(ctx,result,created,operations...)}
func Transactional[S,R any](unit coredatabase.UnitOfWork,handler Handler[S,R])Handler[S,R]{return corepipeline.Transactional(unit,handler)}
func OnError[S,R any](target error,handler func(context.Context,*S,error)Outcome[R])Recovery[S,R]{return corepipeline.OnError(target,handler)}
func OtherwiseError[S,R any](handler func(context.Context,*S,error)Outcome[R])Recovery[S,R]{return corepipeline.OtherwiseError(handler)}
func Recover[S,R any](ctx context.Context,state *S,result Outcome[R],recoveries ...Recovery[S,R])Outcome[R]{return corepipeline.Recover(ctx,state,result,recoveries...)}
