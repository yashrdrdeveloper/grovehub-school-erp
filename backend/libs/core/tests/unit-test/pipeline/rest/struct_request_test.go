package rest_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	coreerrors "grovehub/libs/core/errors"
	corepipeline "grovehub/libs/core/pipeline"
	corerequest "grovehub/libs/core/request"
)

type operationParams struct { GID string `path:"gid" normalize:"trim" validate:"notblank"` }
type operationQuery struct { Page int `query:"page" validate:"min=1"` }
type operationBody struct { Name string `json:"name" normalize:"trim" validate:"notblank"` }

type operationPatch struct {
	corerequest.Request
	Params operationParams `validate:"valid"`
	RequestBody operationBody `validate:"valid"`
}

type operationGet struct {
	corerequest.Request
	Params operationParams `validate:"valid"`
}

type operationList struct {
	corerequest.Request
	Query operationQuery `validate:"valid"`
}

func TestRunStructInfersOmittedFields(t *testing.T) {
	t.Run("params and body with no query",func(t *testing.T){
		req:=httptest.NewRequest(http.MethodPatch,"/people/valid-gid",strings.NewReader(`{"name":"  Ada  "}`))
		req.SetPathValue("gid","valid-gid")
		rec:=httptest.NewRecorder()
		called:=false
		corepipeline.RunStruct(rec,req,corepipeline.StructConfig[operationPatch,operationBody]{
			InvalidParams:coreerrors.NewDefinition("TEST_INVALID_PARAMS","Invalid params",400),
			InvalidRequestBody:coreerrors.NewDefinition("TEST_INVALID_BODY","Invalid body",400),
			Fallback:coreerrors.NewDefinition("TEST_FALLBACK","Failed",500),
			Service:func(_ context.Context, input operationPatch)(corepipeline.Result[operationBody],error){
				called=true
				if input.Params.GID!="valid-gid"||input.RequestBody.Name!="Ada" {t.Fatalf("binding failed: %#v",input)}
				return corepipeline.OK(input.RequestBody),nil
			},
		})
		if !called||rec.Code!=200||!strings.Contains(rec.Body.String(),"Ada"){
			t.Fatalf("unexpected response: called=%v status=%d body=%s",called,rec.Code,rec.Body.String())
		}
	})

	t.Run("params only with no body",func(t *testing.T){
		req:=httptest.NewRequest(http.MethodGet,"/people/valid-gid",nil)
		req.SetPathValue("gid","valid-gid")
		rec:=httptest.NewRecorder()
		corepipeline.RunStruct(rec,req,corepipeline.StructConfig[operationGet,operationBody]{
			Fallback:coreerrors.NewDefinition("TEST_FALLBACK","Failed",500),
			Service:func(_ context.Context, input operationGet)(corepipeline.Result[operationBody],error){
				return corepipeline.OK(operationBody{Name:input.Params.GID}),nil
			},
		})
		if rec.Code!=200 {t.Fatalf("expected 200, got %d: %s",rec.Code,rec.Body.String())}
	})

	t.Run("query only with no params or body",func(t *testing.T){
		req:=httptest.NewRequest(http.MethodGet,"/people?page=2",nil)
		rec:=httptest.NewRecorder()
		corepipeline.RunStruct(rec,req,corepipeline.StructConfig[operationList,operationBody]{
			StrictQuery:true,
			Fallback:coreerrors.NewDefinition("TEST_FALLBACK","Failed",500),
			Service:func(_ context.Context, input operationList)(corepipeline.Result[operationBody],error){
				if input.Query.Page!=2 {t.Fatalf("query was not bound: %#v",input)}
				return corepipeline.OK(operationBody{Name:"ok"}),nil
			},
		})
		if rec.Code!=200 {t.Fatalf("expected 200, got %d: %s",rec.Code,rec.Body.String())}
	})

	t.Run("undeclared query is rejected in strict mode",func(t *testing.T){
		req:=httptest.NewRequest(http.MethodGet,"/people/valid-gid?page=2",nil)
		req.SetPathValue("gid","valid-gid")
		rec:=httptest.NewRecorder()
		corepipeline.RunStruct(rec,req,corepipeline.StructConfig[operationGet,operationBody]{
			StrictQuery:true,
			InvalidQuery:coreerrors.NewDefinition("TEST_INVALID_QUERY","Invalid query",400),
			Fallback:coreerrors.NewDefinition("TEST_FALLBACK","Failed",500),
			Service:func(_ context.Context, _ operationGet)(corepipeline.Result[operationBody],error){
				t.Fatal("unexpected service call")
				return corepipeline.OK(operationBody{}),nil
			},
		})
		if rec.Code!=400||!strings.Contains(rec.Body.String(),"TEST_INVALID_QUERY"){
			t.Fatalf("expected rejected query, got %d: %s",rec.Code,rec.Body.String())
		}
	})
}
