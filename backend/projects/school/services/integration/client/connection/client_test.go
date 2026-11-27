package connection

import (
 "context"
 "encoding/json"
 "io"
 "net/http"
 "net/http/httptest"
 "net/url"
 "strings"
 "testing"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/platformhttp"
)

const testID="e93696cc-15f0-4e06-9577-0d5cfc2a83c8"

func TestBusinessOperations(t *testing.T){
 cases:=[]struct{scope,method,path,body,id,provider,action string}{
  {"integration.create","POST","/connections",`{"provider":"GMAIL","ownerGid":"`+testID+`"}`,"","",""},
  {"integration.get","GET","/connections/"+testID,"",testID,"",""},
  {"integration.list","GET","/connections","","","",""},
  {"integration.disconnect","DELETE","/connections/"+testID,"",testID,"",""},
  {"integration.oauth.start","POST","/connections/"+testID+"/oauth/start","",testID,"",""},
  {"integration.test","POST","/connections/"+testID+"/test",`{"direction":"INBOUND"}`,testID,"",""},
  {"integration.tests","GET","/connections/"+testID+"/tests","",testID,"",""},
  {"integration.watch","POST","/connections/"+testID+"/watch","",testID,"",""},
  {"integration.systemmail","POST","/connections/system-mail/sender",`{"connectionId":"`+testID+`"}`,"","",""},
  {"integration.action","POST","/connections/"+testID+"/actions/mail.send",`{"input":{"empty":[],"nullable":null,"historyId":9007199254740993}}`,testID,"","mail.send"},
  {"integration.providers","GET","/providers","","","",""},
  {"integration.provider","GET","/providers/GMAIL","","","GMAIL",""},
 }
 const response=`{"success":false,"errorCode":"FORBIDDEN","errorMessage":"denied","httpErrorCode":403,"errorMetaData":{}}`
 for _,tc:=range cases{t.Run(tc.scope,func(t *testing.T){
  c:=testClient(t,func(w http.ResponseWriter,r *http.Request){
   if r.Method!=tc.method||r.URL.Path!="/integrations/restricted"+tc.path{t.Errorf("unexpected operation: %s %s",r.Method,r.URL.Path)}
   if r.Header.Get("X-GH-Identity-Proof")!="proof"||r.Header.Get("Cookie")!=""{t.Error("incorrect credential")}
   raw,err:=io.ReadAll(r.Body);if err!=nil||string(raw)!=tc.body{t.Errorf("body changed: %q",raw)}
   w.Header().Set("Content-Type","application/json");w.WriteHeader(403);_,_=io.WriteString(w,response)
  })
  result,err:=c.Business(context.Background(),tc.scope,tc.id,tc.provider,tc.action,json.RawMessage(tc.body),"proof")
  if err!=nil||result.Status!=403||string(result.Body)!=response{t.Fatalf("unexpected result: %v",err)}
 })}
}

func TestOAuthUsesStateWithoutBrowserProof(t *testing.T){
 state:=strings.Repeat("x",43)
 c:=testClient(t,func(w http.ResponseWriter,r *http.Request){
  if r.Method!="GET"||r.URL.Path!="/integrations/restricted/connections/oauth/gmail/callback"||r.URL.Query().Get("code")!="a+b/c="||r.URL.Query().Get("state")!=state{t.Error("incorrect callback")}
  if r.Header.Get("X-GH-Identity-Proof")!=""||r.Header.Get("Cookie")!=""{t.Error("unexpected browser credential")}
  w.Header().Set("Content-Type","application/json");_,_=io.WriteString(w,`{"success":true,"data":{"status":"ACTIVE"}}`)
 })
 result,err:=c.FinishOAuth(platformhttp.WithIdentityProof(context.Background(),"do-not-forward"),"gmail","a+b/c=",state)
 if err!=nil||result.Status!=200{t.Fatalf("callback failed: %v",err)}
 if _,err:=c.FinishOAuth(context.Background(),"gmail","code","invalid");err==nil{t.Fatal("invalid state accepted")}
}

func testClient(t *testing.T, handler http.HandlerFunc)*Client{
 t.Helper();server:=httptest.NewServer(handler);t.Cleanup(server.Close)
 origin,err:=url.Parse(server.URL);if err!=nil{t.Fatal(err)}
 client:=server.Client();client.Transport=platformhttp.ProofTransport{Base:client.Transport}
 return New(origin,client)
}
