package gateway

import (
 "context"
 "crypto/ed25519"
 "crypto/rand"
 "io"
 "net/http"
 "net/url"
 "strings"
 "testing"
)

type roundTripFunc func(*http.Request)(*http.Response,error)
func (f roundTripFunc) RoundTrip(r *http.Request)(*http.Response,error){return f(r)}

func TestWorkloadTransport(t *testing.T){
 pub,priv,err:=ed25519.GenerateKey(rand.Reader)
 if err!=nil {t.Fatal(err)}
 origin,_:=url.Parse("https://platform.example")
 verifier:=Ed25519Credentials{Keys:map[string]TrustedWorkloadKey{"key":{PublicKey:pub,Issuer:"school",Subject:"school",Scopes:[]string{"audit.search"},AllowDelegation:true}}}
 calls:=0
 transport:=&WorkloadTransport{Origin:origin,KeyID:"key",Issuer:"school",Subject:"school",Audience:"platform",PrivateKey:priv,Policy:func(r *http.Request)(WorkloadClaims,error){return WorkloadClaims{Scopes:[]string{"audit.search"},TenantID:"tenant",ActorGID:"actor"},nil},Base:roundTripFunc(func(r *http.Request)(*http.Response,error){
 calls++
 if r.Header.Get("Cookie")!="" || r.Header.Get("X-Actor")!="" {t.Fatal("untrusted header forwarded")}
 if r.Header.Get("Idempotency-Key")!="id" {t.Fatal("idempotency key lost")}
 claims,err:=verifier.VerifyCredential(context.Background(),strings.TrimPrefix(r.Header.Get("Authorization"),"Bearer "))
 if err!=nil || claims.ActorGID!="actor" || claims.Audience!="platform" {t.Fatalf("invalid signed claims: %v",err)}
 return &http.Response{StatusCode:302,Header:http.Header{"Location":[]string{"https://other.example/"}},Body:io.NopCloser(strings.NewReader("")),Request:r},nil
 })}
 client,err:=NewWorkloadHTTPClient(transport)
 if err!=nil {t.Fatal(err)}
 req,_:=http.NewRequest("POST","https://platform.example/audit/restricted/search",nil)
 req.Header.Set("Cookie","session=secret")
 req.Header.Set("Authorization","Bearer browser")
 req.Header.Set("X-Actor","forged")
 req.Header.Set("Idempotency-Key","id")
 resp,err:=client.Do(req)
 if err!=nil {t.Fatal(err)}
 defer resp.Body.Close()
 if resp.StatusCode!=302 || calls!=1 {t.Fatal("redirect followed")}
 if req.Header.Get("Authorization")!="Bearer browser" {t.Fatal("original request mutated")}
 for _, target:=range []string{"https://other.example/","http://platform.example/"} {
  r,_:=http.NewRequest("GET",target,nil)
  if _,err:=client.Do(r);err==nil {t.Fatalf("unexpectedly accepted %s",target)}
 }
 if calls!=1 {t.Fatal("denied request reached network")}
}
