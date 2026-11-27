package client
import("net/http/httptest";"testing")
func TestOperationAllowlist(t *testing.T){
 seen:=map[string]bool{}
 for _,op:=range Operations(){key:=op.Method+" "+op.Path;if seen[key]{t.Fatal("duplicate",key)};seen[key]=true}
 cases:=[]struct{method,path,scope string}{
 {"POST","/identity/restricted/accounts/search","identity.accounts.search"},
 {"GET","/integrations/restricted/providers/gmail","integration.provider"},
 {"POST","/integrations/restricted/connections/00000000-0000-4000-8000-000000000000/actions/mail.send","integration.action"},
 {"GET","/foundation/restricted/persons",""},
 {"GET","/identity/restricted/accounts/not-a-uuid",""},
 {"GET","/integrations/restricted/providers/unknown",""},
 {"POST","/identity/public/accounts/search",""},
 }
 for _,c:=range cases{r:=httptest.NewRequest(c.method,c.path,nil);got,ok:=Scope(r);if got!=c.scope||ok!=(c.scope!=""){t.Errorf("%s %s: got %q %v",c.method,c.path,got,ok)}}
}
