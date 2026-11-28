// Package rest is a temporary compatibility facade. Use core/pipeline.
package rest

import (
	"net/http"
	coreannotation "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation"
	corerequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
	corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
)

type Empty = corepipeline.Empty
type Input[P,Q,B any] = corepipeline.Input[P,Q,B]
type Result[T any] = corepipeline.Result[T]
type Config[P,Q,B,R any] = corepipeline.Config[P,Q,B,R]
type StructConfig[Q,R any] = corepipeline.StructConfig[Q,R]
type Writer[T any] = corepipeline.Writer[T]
type Failure = corepipeline.Failure
type BindError = corepipeline.BindError
type DecodeError = corepipeline.DecodeError

func OK[T any](data T) Result[T] { return corepipeline.OK(data) }
func Created[T any](data T) Result[T] { return corepipeline.Created(data) }
func CreatedOrOK[T any](data T, created bool) Result[T] { return corepipeline.CreatedOrOK(data, created) }
func ReadBody[T any](r *http.Request, decode func(*http.Request,any)error, rules coreannotation.ValidationRules[T])(T,error){return corepipeline.ReadBody(r,decode,rules)}
func Run[P,Q,B,R any](w http.ResponseWriter,r *http.Request,c Config[P,Q,B,R]){corepipeline.RunREST(w,r,c)}
func RunStruct[Q corerequest.Contract,R any](w http.ResponseWriter,r *http.Request,c StructConfig[Q,R]){corepipeline.RunStruct(w,r,c)}
