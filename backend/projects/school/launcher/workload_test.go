package launcher

import (
 "net/http"
 "net/http/httptest"
 "testing"
)

func TestServiceWorkloadIsolation(t *testing.T){
 const id="e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
 cases:=[]struct{method,path,scope string;owner int}{
 {"POST","/identity/restricted/auth/login","identity.auth.login",0},
 {"DELETE","/identity/restricted/auth/sessions/"+id,"identity.auth.sessions.revoke",0},
 {"POST","/identity/restricted/accounts/search","identity.accounts.search",0},
 {"POST","/audit/restricted/search","audit.search",1},
 {"GET","/integrations/restricted/providers/gmail","integration.provider",2},
 {"GET","/integrations/restricted/connections/oauth/gmail/callback?code=a%2Bb&state=x","integration.oauth.callback",2},
 {"POST","/foundation/restricted/persons","",-1},
 {"POST","/integrations/restricted/webhooks/gmail","",-1},
 {"GET","/identity/restricted/auth/me?extra=x","",-1},
 {"GET","/integrations/restricted/connections/oauth/gmail/callback","",-1},
 {"GET","/identity/protected/auth/me","",-1},
 }
 for _,tc:=range cases{
  for i,policy:=range []workloadPolicy{identityPolicy,auditPolicy,integrationPolicy}{
   r:=httptest.NewRequest(tc.method,tc.path,nil)
   claims,err:=policy(r)
   if i==tc.owner {if err!=nil||len(claims.Scopes)!=1||claims.Scopes[0]!=tc.scope{t.Errorf("%s: expected %s, got %v %v",tc.path,tc.scope,claims,err)}} else if err==nil{t.Errorf("%s: unexpected permission from service %d",tc.path,i)}
  }
 }
 for _,policy:=range []workloadPolicy{identityPolicy,auditPolicy,integrationPolicy}{if _,err:=policy((*http.Request)(nil));err==nil{t.Fatal("nil request accepted")}}
}
