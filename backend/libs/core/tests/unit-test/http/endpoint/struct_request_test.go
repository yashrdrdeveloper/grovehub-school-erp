package endpoint_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	endpoint "grovehub/libs/core/http/endpoint"
	corerequest "grovehub/libs/core/request"
	httpresponse "grovehub/libs/core/http/response"
)

type onlyParamsRequest struct {
	corerequest.Request
	Params personPath
}

type paramsBodyRequest struct {
	corerequest.Request
	Params personPath
	RequestBody struct { Name string `json:"name"` }
}

type queryRequest struct {
	corerequest.Request
	Query queryInputs
}

func TestCallStructInfersAbsentSections(t *testing.T) {
	mux:=http.NewServeMux()
	mux.HandleFunc("GET /people/{gid}",func(w http.ResponseWriter,r *http.Request){
		if r.URL.RawQuery!=""||r.ContentLength!=0 {t.Fatalf("unexpected undeclared query/body")}
		if err:=httpresponse.WriteSuccess(w,200,personRecord{GID:r.PathValue("gid")});err!=nil{t.Fatal(err)}
	})
	res,status,err:=endpoint.CallStruct[personRecord](context.Background(),mux,http.MethodGet,"/people/{gid}",onlyParamsRequest{Params:personPath{GID:"valid-gid"}})
	if err!=nil||status!=200||res.GID!="valid-gid" {t.Fatalf("result=%+v status=%d err=%v",res,status,err)}
}

func TestCallStructMapsOnlyDeclaredParamsAndBody(t *testing.T) {
	mux:=http.NewServeMux()
	mux.HandleFunc("PATCH /people/{gid}",func(w http.ResponseWriter,r *http.Request){
		if r.URL.RawQuery!="" {t.Fatal("query was not declared")}
		var body struct{Name string `json:"name"`}
		if err:=json.NewDecoder(r.Body).Decode(&body);err!=nil{t.Fatal(err)}
		if body.Name!="Ada"||r.PathValue("gid")!="valid-gid"{t.Fatalf("bad binding")}
		if err:=httpresponse.WriteSuccess(w,200,personRecord{Name:body.Name});err!=nil{t.Fatal(err)}
	})
	input:=paramsBodyRequest{Params:personPath{GID:"valid-gid"}}
	input.RequestBody.Name="Ada"
	res,status,err:=endpoint.CallStruct[personRecord](context.Background(),mux,http.MethodPatch,"/people/{gid}",input)
	if err!=nil||status!=200||res.Name!="Ada"{t.Fatalf("result=%+v status=%d err=%v",res,status,err)}
}

func TestCallStructMapsOnlyDeclaredQuery(t *testing.T) {
	mux:=http.NewServeMux()
	mux.HandleFunc("GET /people",func(w http.ResponseWriter,r *http.Request){
		if r.URL.Query().Get("page")!="4" {t.Fatalf("query was not transmitted")}
		if err:=httpresponse.WriteSuccess(w,200,personRecord{Name:"query"});err!=nil{t.Fatal(err)}
	})
	_,status,err:=endpoint.CallStruct[personRecord](context.Background(),mux,http.MethodGet,"/people",queryRequest{Query:queryInputs{Page:4}})
	if err!=nil||status!=200{t.Fatalf("status=%d err=%v",status,err)}
}
