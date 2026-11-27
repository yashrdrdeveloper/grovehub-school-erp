package authentication

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestGeneratedBusinessPreservesAccountPatch(t *testing.T) {
	const id = "e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
	const input = `{"username":null,"roleAssignment":{"primaryRoleId":"`+id+`","secondaryRoleIds":[]},"version":9007199254740993}`
	const response = `{"success":false,"errorCode":"FORBIDDEN","errorMessage":"denied","httpErrorCode":403,"errorMetaData":{}}`
	client := newTestClient(t,func(w http.ResponseWriter,r *http.Request){
		if r.Method!="PATCH" || r.URL.Path!="/identity/restricted/accounts/"+id {t.Error("wrong operation")}
		if r.Header.Get("X-GH-Identity-Proof")!="proof" || r.Header.Get("Cookie")!="" {t.Error("wrong credential")}
		raw,err:=io.ReadAll(r.Body)
		if err!=nil || string(raw)!=input || r.ContentLength!=int64(len(input)) {t.Error("PATCH payload altered")}
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(http.StatusForbidden)
		_,_=io.WriteString(w,response)
	})
	result,err:=client.Business(context.Background(),"identity.accounts.update",id,json.RawMessage(input),"proof","")
	if err!=nil || result.Status!=403 || string(result.Body)!=response {t.Fatalf("rejection not preserved: %v",err)}
}

func TestGeneratedBusinessForwardsIdempotency(t *testing.T) {
	const input = `{"person":{"firstName":"Test","middleName":null,"lastName":"User"},"email":"test@example.com"}`
	const response = `{"success":false,"errorCode":"CONFLICT","errorMessage":"conflict","httpErrorCode":409,"errorMetaData":{}}`
	client:=newTestClient(t,func(w http.ResponseWriter,r *http.Request){
		if r.Method!="POST" || r.URL.Path!="/identity/restricted/accounts" || r.Header.Get("Idempotency-Key")!="test-idempotency" {t.Error("incorrect request")}
		raw,_:=io.ReadAll(r.Body)
		if string(raw)!=input {t.Error("account body changed")}
		w.Header().Set("Content-Type","application/json")
		w.WriteHeader(409)
		_,_=io.WriteString(w,response)
	})
	result,err:=client.Business(context.Background(),"identity.accounts.create","",json.RawMessage(input),"proof","test-idempotency")
	if err!=nil || result.Status!=409 || string(result.Body)!=response {t.Fatalf("unexpected response: %v",err)}
}

func TestGeneratedBusinessRejectsInvalidDispatch(t *testing.T) {
	client:=newTestClient(t,func(w http.ResponseWriter,r *http.Request){t.Error("invalid operation reached Platform")})
	for _,tc:=range []struct{scope,id,body,key string}{
		{"identity.unknown","","",""},
		{"identity.accounts.get","../invalid","",""},
		{"identity.accounts.create","","null",""},
		{"identity.accounts.search","","[]",""},
		{"identity.roles.get","e93696cc-15f0-4e06-9577-0d5cfc2a83c8","{}",""},
		{"identity.accounts.create","",`{"person":{"firstName":"A","lastName":"B"},"email":"a@example.com"}`,"bad\nkey"},
	}{
		_,err:=client.Business(context.Background(),tc.scope,tc.id,json.RawMessage(tc.body),"proof",tc.key)
		if !errors.Is(err,ErrInvalidBusinessRequest){t.Errorf("%s: expected invalid request, got %v",tc.scope,err)}
	}
}
