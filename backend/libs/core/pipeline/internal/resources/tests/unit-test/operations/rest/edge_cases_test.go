package rest_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	corerequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
	rest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/internal/src/operations/rest"
)

type numericQuery struct {
	Unsigned uint16 `query:"unsigned"`
	Fraction float64 `query:"fraction"`
	Ints []int `query:"int"`
}
type unsupportedQuery struct { Value map[string]string `query:"value"` }
type hiddenQuery struct { secret string `query:"secret"` }
type invalidShape struct { corerequest.Request; Extra string }
type invalidQueryShape struct { corerequest.Request; Query string }
type invalidParamsShape struct { corerequest.Request; Params string }
type hiddenQueryShape struct { corerequest.Request; Query hiddenQuery }
type noSections struct { corerequest.Request }

var edgeInvalidParams = coreerrors.NewDefinition("EDGE_PARAMS","Invalid params",400)
var edgeInvalidQuery = coreerrors.NewDefinition("EDGE_QUERY","Invalid query",400)
var edgeInvalidBody = coreerrors.NewDefinition("EDGE_BODY","Invalid body",400)
var edgeFallback = coreerrors.NewDefinition("EDGE_INTERNAL","Internal error",500)

func runQueryCase[Q any](url string, strict bool, service func(context.Context,rest.Input[rest.Empty,Q,rest.Empty])(rest.Result[string],error)) *httptest.ResponseRecorder {
	w:=httptest.NewRecorder()
	rest.Run(w,httptest.NewRequest("GET",url,nil),rest.Config[rest.Empty,Q,rest.Empty,string]{StrictQuery:strict,InvalidQuery:edgeInvalidQuery,Fallback:edgeFallback,Service:service})
	return w
}

func unexpectedQuery[Q any](context.Context,rest.Input[rest.Empty,Q,rest.Empty])(rest.Result[string],error){
	Fail("unexpected service execution")
	return rest.OK(""),nil
}

func runStructCase[Q corerequest.Contract](r *http.Request, strict bool, service func(context.Context,Q)(rest.Result[string],error), validate func(Q)error) *httptest.ResponseRecorder {
	w:=httptest.NewRecorder()
	rest.RunStruct(w,r,rest.StructConfig[Q,string]{StrictQuery:strict,InvalidParams:edgeInvalidParams,InvalidQuery:edgeInvalidQuery,InvalidRequestBody:edgeInvalidBody,Fallback:edgeFallback,Service:service,Validate:validate})
	return w
}

func unexpectedStruct[Q corerequest.Contract](context.Context,Q)(rest.Result[string],error){
	Fail("unexpected service execution")
	return rest.OK(""),nil
}

var _=Describe("REST binding edge cases",func(){
	It("binds unsigned and fractional scalars and numeric slices",func(){
		w:=runQueryCase("/?unsigned=12&fraction=1.25&int=2&int=3",true,func(_ context.Context,input rest.Input[rest.Empty,numericQuery,rest.Empty])(rest.Result[string],error){
			Expect(input.Query.Unsigned).To(Equal(uint16(12)))
			Expect(input.Query.Fraction).To(Equal(1.25))
			Expect(input.Query.Ints).To(Equal([]int{2,3}))
			return rest.OK("bound"),nil
		})
		Expect(w.Code).To(Equal(200))
		for _,url:=range []string{"/?unsigned=-1","/?fraction=bad","/?int=1&int=bad"}{
			Expect(runQueryCase(url,true,unexpectedQuery[numericQuery]).Code).To(Equal(400))
		}
	})
	It("rejects unsupported DTO types, hidden fields and missing path values",func(){
		Expect(runQueryCase("/?value=1",true,unexpectedQuery[string]).Code).To(Equal(400))
		Expect(runQueryCase("/",false,unexpectedQuery[string]).Code).To(Equal(400))
		Expect(runQueryCase("/?value=1",false,unexpectedQuery[unsupportedQuery]).Code).To(Equal(400))
		Expect(runQueryCase("/?secret=1",false,unexpectedQuery[hiddenQuery]).Code).To(Equal(400))
		w:=httptest.NewRecorder()
		rest.Run(w,httptest.NewRequest("GET","/",nil),rest.Config[personPath,rest.Empty,rest.Empty,string]{InvalidPath:edgeInvalidParams,Fallback:edgeFallback,Service:func(context.Context,rest.Input[personPath,rest.Empty,rest.Empty])(rest.Result[string],error){
			Fail("unexpected service");return rest.OK(""),nil
		}})
		Expect(w.Code).To(Equal(400))
		w=httptest.NewRecorder()
		rest.Run(w,httptest.NewRequest("GET","/",nil),rest.Config[rest.Empty,rest.Empty,rest.Empty,string]{Fallback:edgeFallback})
		Expect(w.Code).To(Equal(500))
	})
	It("returns safe typed error messages",func(){
		Expect(rest.BindError{Source:"query",Field:"count",Reason:"malformed"}.Error()).To(Equal("invalid query parameter count: malformed"))
		err:=&rest.DecodeError{Err:errors.New("invalid JSON")}
		Expect(err.Error()).To(Equal("decode request: invalid JSON"))
		Expect(err.Unwrap()).To(MatchError("invalid JSON"))
	})
})

var _=Describe("Structured REST failure paths",func(){
	It("rejects missing callbacks and malformed request contracts",func(){
		r:=httptest.NewRequest("GET","/",nil)
		Expect(runStructCase[noSections](r,false,nil,nil).Code).To(Equal(500))
		Expect(runStructCase(r,false,unexpectedStruct[invalidShape],nil).Code).To(Equal(500))
		Expect(runStructCase(r,true,unexpectedStruct[invalidQueryShape],nil).Code).To(Equal(400))
		Expect(runStructCase(r,false,unexpectedStruct[invalidQueryShape],nil).Code).To(Equal(400))
		Expect(runStructCase(r,false,unexpectedStruct[invalidParamsShape],nil).Code).To(Equal(400))
		Expect(runStructCase(r,false,unexpectedStruct[hiddenQueryShape],nil).Code).To(Equal(400))
		Expect(runStructCase(r,false,unexpectedStruct[operationGet],nil).Code).To(Equal(400))
	})
	It("rejects malformed undeclared queries and invalid JSON bodies",func(){
		r:=httptest.NewRequest("GET","/?value=%zz",nil)
		Expect(runStructCase(r,true,unexpectedStruct[noSections],nil).Code).To(Equal(400))
		r=httptest.NewRequest("PATCH","/",strings.NewReader("{"))
		r.SetPathValue("gid","id")
		Expect(runStructCase(r,false,unexpectedStruct[operationPatch],nil).Code).To(Equal(400))
	})
	It("classifies validation errors by the declared source",func(){
		r:=httptest.NewRequest("GET","/",nil)
		r.SetPathValue("gid"," ")
		w:=runStructCase(r,false,unexpectedStruct[operationGet],nil)
		Expect(w.Code).To(Equal(400));Expect(w.Body.String()).To(ContainSubstring("EDGE_PARAMS"))
		w=runStructCase(httptest.NewRequest("GET","/?page=0",nil),false,unexpectedStruct[operationList],nil)
		Expect(w.Code).To(Equal(400));Expect(w.Body.String()).To(ContainSubstring("EDGE_QUERY"))
		r=httptest.NewRequest("PATCH","/",strings.NewReader(`{"name":" "}`))
		r.SetPathValue("gid","id")
		w=runStructCase(r,false,unexpectedStruct[operationPatch],nil)
		Expect(w.Code).To(Equal(400));Expect(w.Body.String()).To(ContainSubstring("EDGE_BODY"))
	})
	It("executes cross-source validation before the service and handles service failure",func(){
		r:=httptest.NewRequest("GET","/",nil)
		w:=runStructCase(r,false,unexpectedStruct[noSections],func(noSections)error{return errors.New("validation failed")})
		Expect(w.Code).To(Equal(400))
		w=runStructCase(r,false,func(context.Context,noSections)(rest.Result[string],error){return rest.Result[string]{},errors.New("private failure")},func(noSections)error{return nil})
		Expect(w.Code).To(Equal(500))
		Expect(w.Body.String()).NotTo(ContainSubstring("private failure"))
	})
})
