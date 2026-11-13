package request_test

import (
	"testing"

	corerequest "grovehub/libs/core/request"
)

type onlyBody struct {
	corerequest.Request
	RequestBody struct{ Name string }
}

type onlyParams struct {
	corerequest.Request
	Params struct{ GID string }
}

type mixed struct {
	corerequest.Request
	Params struct{ GID string }
	Query struct{ Page int }
	RequestBody struct{ Name string }
}

type notARequest struct{ RequestBody string }

type unknownSection struct {
	corerequest.Request
	Foo string
}

func TestCommonRequestAutoOmittedSections(t *testing.T) {
	tests := []struct{
		name string
		value corerequest.Contract
		params, query, body bool
	}{
		{"body only",onlyBody{},false,false,true},
		{"params only",onlyParams{},true,false,false},
		{"all sections",mixed{},true,true,true},
		{"no sections",corerequest.Request{},false,false,false},
	}
	for _,tt := range tests {
		t.Run(tt.name,func(t *testing.T){
			s,err:=corerequest.Inspect(tt.value)
			if err!=nil {t.Fatal(err)}
			if s.HasParams!=tt.params || s.HasQuery!=tt.query || s.HasRequestBody!=tt.body {
				t.Fatalf("wrong sections: %+v",s)
			}
		})
	}
}

func TestRequestRequiresCommonContract(t *testing.T) {
	if _,err:=corerequest.Inspect(notARequest{});err==nil { t.Fatal("expected missing common Request marker rejection") }
	if _,err:=corerequest.Inspect(unknownSection{});err==nil { t.Fatal("expected unknown section rejection") }
	var ptr *onlyBody
	if _,err:=corerequest.Inspect(ptr);err==nil { t.Fatal("expected nil request rejection") }
}

func TestEmbeddedRequestTypeIsAddressable(t *testing.T) {
	request:=onlyBody{}
	parts,err:=corerequest.Inspect(&request)
	if err!=nil {t.Fatal(err)}
	if !parts.RequestBody.CanAddr() {t.Fatal("RequestBody must be writable for HTTP binder")}
	parts.RequestBody.FieldByName("Name").SetString("Set")
	if request.RequestBody.Name!="Set" {t.Fatal("binder did not mutate the request")}
}
