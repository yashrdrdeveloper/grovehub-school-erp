package workflow

import (
 "encoding/json"
 "testing"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

func TestTemplateValidation(t *testing.T){
 valid:=TemplateInput{Code:"account-get",Module:"identity",Method:"GET",Endpoint:"/identity/restricted/accounts/${accountId}",Scope:"SCHOOL",BodySchema:map[string]entity.Field{"accountId":{Type:"string",Required:true}}}
 if !validateTemplate(valid){t.Fatal("valid template rejected")}
 for _,path:=range []string{"https://example.com/identity/restricted/accounts","/identity/public/accounts","/identity/restricted/../accounts","/identity/restricted/accounts/${missing}","/identity/restricted/accounts?x=1","/foundation/restricted/persons"}{
  bad:=valid;bad.Endpoint=path;if validateTemplate(bad){t.Errorf("accepted %q",path)}
 }
 bad:=valid;bad.Headers=map[string]string{"Authorization":"secret"};if validateTemplate(bad){t.Fatal("authorization header accepted")}
}
func TestInputValidation(t *testing.T){
 schema:=map[string]entity.Field{"version":{Type:"integer",Required:true},"enabled":{Type:"boolean"}}
 if !validInput(json.RawMessage(`{"version":9007199254740993,"enabled":false}`),schema){t.Fatal("valid input rejected")}
 for _,raw:=range []string{`{}`,`null`,`{"version":null}`,`{"version":1.5}`,`{"version":1,"extra":"x"}`,`{"version":1,"enabled":"true"}`}{
  if validInput(json.RawMessage(raw),schema){t.Errorf("accepted %s",raw)}
 }
}
