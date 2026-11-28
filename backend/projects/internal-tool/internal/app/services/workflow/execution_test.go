package workflow

import (
 "encoding/json"
 "testing"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
func TestResolveApprovedInvocation(t *testing.T){
 tpl:=entity.Template{Method:"PATCH",Endpoint:"/identity/restricted/accounts/${accountId}",Headers:map[string]string{"Idempotency-Key":"${key}"},BodySchema:map[string]entity.Field{"accountId":{Type:"string",Required:true},"key":{Type:"string",Required:true}}}
 a:=entity.Approval{ID:"request",Input:json.RawMessage(`{"accountId":"safe-id","key":"stable-key"}`)}
 call,err:=resolve(tpl,a)
 if err!=nil||call.Path!="/identity/restricted/accounts/safe-id"||call.Headers["Idempotency-Key"]!="stable-key"{t.Fatal("resolution failed",call,err)}
 for _,id:=range []string{"..","a/b","%2f","a?b","a#b","a\nb"}{
  raw,_:=json.Marshal(map[string]string{"accountId":id,"key":"stable-key"});a.Input=raw
  if _,err=resolve(tpl,a);err==nil{t.Errorf("unsafe path parameter accepted %q",id)}
 }
}
