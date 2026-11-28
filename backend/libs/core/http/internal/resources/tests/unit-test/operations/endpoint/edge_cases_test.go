package endpoint_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	headercontext "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext"
	corerequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/internal/src/operations/endpoint"
)

func expectQueryFailure[Q any](input Q) {
	handler:=http.HandlerFunc(func(http.ResponseWriter,*http.Request){Fail("unexpected dispatch")})
	_,_,err:=endpoint.CallRequest[string](context.Background(),handler,"GET","/",endpoint.Request[corerequest.Empty,Q,corerequest.Empty]{Query:input})
	Expect(err).To(HaveOccurred())
}

type edgeQuery struct {
	String string `query:"s"`
	Bool bool `query:"b"`
	Int int64 `query:"i"`
	Uint uint16 `query:"u"`
	Float float64 `query:"f"`
	Repeated []int `query:"r"`
	Optional *string `query:"o"`
	Absent *int `query:"absent"`
	Ignored string `query:"-"`
	Untagged string
}

type edgeInvalidRequest struct { corerequest.Request; Extra string }
type edgeInvalidQueryRequest struct { corerequest.Request; Query map[string]string }

var _=Describe("Endpoint input encoding",func(){
	It("encodes scalar, pointer and repeated query fields and trusted metadata",func(){
		note:="a & b"
		q:=edgeQuery{String:"Ada",Bool:true,Int:-3,Uint:9,Float:1.25,Repeated:[]int{1,2},Optional:&note}
		ctx:=headercontext.WithValues(context.Background(),headercontext.Values{IdempotencyKey:"key",CorrelationID:"correlation",RequestID:"request"})
		handler:=http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
			Expect(r.URL.Query().Get("s")).To(Equal("Ada"))
			Expect(r.URL.Query().Get("b")).To(Equal("true"))
			Expect(r.URL.Query().Get("i")).To(Equal("-3"))
			Expect(r.URL.Query().Get("u")).To(Equal("9"))
			Expect(r.URL.Query().Get("f")).To(Equal("1.25"))
			Expect(r.URL.Query()["r"]).To(Equal([]string{"1","2"}))
			Expect(r.URL.Query().Get("o")).To(Equal(note))
			Expect(r.URL.Query().Has("absent")).To(BeFalse())
			Expect(r.Header.Get("Idempotency-Key")).To(Equal("key"))
			Expect(r.Header.Get("X-Correlation-ID")).To(Equal("correlation"))
			Expect(r.Header.Get("X-Request-ID")).To(Equal("request"))
			Expect(r.Header.Get("Content-Type")).To(Equal("application/json"))
			w.WriteHeader(201)
			fmt.Fprint(w,`{"success":true,"data":"ok"}`)
		})
		_,status,err:=endpoint.CallRequest[string](ctx,handler,"POST","/",endpoint.Request[corerequest.Empty,*edgeQuery,string]{Query:&q,Body:"body"})
		Expect(err).NotTo(HaveOccurred());Expect(status).To(Equal(201))
		var absent *edgeQuery
		_,status,err=endpoint.CallRequest[string](ctx,http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){Expect(r.URL.RawQuery).To(BeEmpty());w.WriteHeader(204)}),"GET","/",endpoint.Request[corerequest.Empty,*edgeQuery,corerequest.Empty]{Query:absent})
		Expect(err).NotTo(HaveOccurred());Expect(status).To(Equal(204))
	})
	It("rejects unsupported and duplicate query fields",func(){
		expectQueryFailure("not a struct")
		expectQueryFailure(struct{hidden string `query:"hidden"`}{hidden:"x"})
		expectQueryFailure(struct{A string `query:"x"`;B string `query:"x"`}{A:"1",B:"2"})
		expectQueryFailure(struct{A []map[string]string `query:"a"`}{A:[]map[string]string{{"a":"b"}}})
		expectQueryFailure(struct{A map[string]string `query:"a"`}{A:map[string]string{"a":"b"}})
	})
	It("validates path DTO shape, scalar types, tags, and traversal protection",func(){
		handler:=http.HandlerFunc(func(http.ResponseWriter,*http.Request){Fail("unexpected dispatch")})
		var nilDTO *struct{ID string `path:"id"`}
		for _,item:=range []struct{route string; dto any}{
			{"bad",nil},{"/bad?x=1",nil},{"/bad/{1}",nil},{"/{id}",nil},{"/{id}",nilDTO},{"/{id}",1},
			{"/{id}",struct{hidden string `path:"id"`}{hidden:"a"}},
			{"/{id}",struct{A string `path:"id"`; B string `path:"id"`}{A:"a",B:"b"}},
			{"/{id}",struct{ID *string `path:"id"`}{}},
			{"/{id}",struct{ID bool `path:"id"`}{ID:true}},
			{"/{id}/tail",struct{}{}},
			{"/{id}",struct{ID string `path:"id"`}{ID:"."}},
		}{
			_,_,err:=endpoint.Call[string](context.Background(),handler,"GET",item.route,item.dto,nil)
			Expect(err).To(HaveOccurred(),item.route)
		}
		type params struct { Int *int `path:"i"`; Uint uint `path:"u"`; Omitted string `path:"-"` }
		number:=-2
		_,status,err:=endpoint.Call[string](context.Background(),http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){Expect(r.URL.Path).To(Equal("/-2/3"));w.WriteHeader(204)}),"GET","/{i}/{u}",&params{Int:&number,Uint:3},nil)
		Expect(err).NotTo(HaveOccurred());Expect(status).To(Equal(204))
	})
	It("rejects malformed endpoint requests and structured request contracts",func(){
		handler:=http.HandlerFunc(func(http.ResponseWriter,*http.Request){Fail("unexpected dispatch")})
		ctx,cancel:=context.WithCancel(context.Background())
		cancel()
		_,_,err:=endpoint.Call[string](ctx,handler,"GET","/",nil,nil)
		Expect(err).To(MatchError(context.Canceled))
		_,_,err=endpoint.Call[string](context.Background(),handler,"invalid method","/",nil,nil)
		Expect(err).To(HaveOccurred())
		_,_,err=endpoint.Call[string](context.Background(),handler,"POST","/",nil,make(chan int))
		Expect(err).To(HaveOccurred())
		_,_,err=endpoint.CallStruct[string](context.Background(),handler,"GET","/",edgeInvalidRequest{})
		Expect(err).To(HaveOccurred())
		_,_,err=endpoint.CallStruct[string](context.Background(),handler,"GET","/",edgeInvalidQueryRequest{})
		Expect(err).To(HaveOccurred())
	})
})

var _=Describe("Endpoint response classification",func(){
	It("handles empty data and rejects unexpected or mismatched envelopes",func(){
		for _,item:=range []struct{status int;body string;failure bool}{
			{0,"",true},{199,"",true},{600,"",true},
			{200,`{"success":true,"data":null}`,false},
			{200,`{"success":false}`,true},
			{200,`{"success":true,"data":12}`,true},
			{400,`{"success":false,"errorCode":"INVALID","httpErrorCode":404}`,true},
		}{
			h:=http.HandlerFunc(func(w http.ResponseWriter,_ *http.Request){if item.status!=0 {w.WriteHeader(item.status)};if item.body!="" {fmt.Fprint(w,item.body)}})
			value,_,err:=endpoint.Call[string](context.Background(),h,"GET","/",nil,nil)
			if item.failure {Expect(err).To(HaveOccurred())}else{Expect(err).NotTo(HaveOccurred());Expect(value).To(BeEmpty())}
		}
	})
	It("preserves structured application errors and their optional metadata",func(){
		for _,metadata:=range []string{"null","{}",`{"field":"name"}`}{
			h:=http.HandlerFunc(func(w http.ResponseWriter,_ *http.Request){w.WriteHeader(409);fmt.Fprint(w,`{"success":false,"errorCode":"CONFLICT","errorMessage":"Conflict","httpErrorCode":409,"errorMetaData":`+metadata+"}")})
			_,status,err:=endpoint.Call[string](context.Background(),h,"GET","/",nil,nil)
			Expect(status).To(Equal(409))
			var app *coreerrors.Error
			Expect(err).To(BeAssignableToTypeOf(app))
			app=err.(*coreerrors.Error)
			Expect(app.Definition().Code()).To(Equal(coreerrors.Code("CONFLICT")))
			if strings.Contains(metadata,"field"){Expect(app.Metadata()).To(Equal(map[string]any{"field":"name"}))}else{Expect(app.Metadata()).To(BeNil())}
		}
	})
})
