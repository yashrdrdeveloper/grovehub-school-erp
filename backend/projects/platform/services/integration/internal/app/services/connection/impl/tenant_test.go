package impl

import (
 "context"
 "encoding/json"
 "errors"
 "testing"

 schoolgateway "grovehub/launcher/services/gateway/session"
 dto "grovehub/services/integration/internal/app/dto/connection"
 service "grovehub/services/integration/internal/app/services/connection"
 helpers "grovehub/services/integration/internal/app/services/connection/helpers"
)

type boundaryFake struct{
 service.ConnectionService
 calls int
}
func(f *boundaryFake) ExecuteSystemMail(ctx context.Context,tenant string,_ json.RawMessage)(any,error){
 f.calls++
 if _,err:=helpers.RequireTenantGeneration(ctx,tenant);err!=nil{return nil,err}
 return "sent",nil
}
func(f *boundaryFake) Get(context.Context,string,string)(dto.Connection,error){f.calls++;return dto.Connection{},nil}

func TestDirectIntegrationBoundary(t *testing.T){
 base:=&boundaryFake{}
 active,available:=false,true
 guarded:=WithTenantState(base,func(_ context.Context,tenant string)(schoolgateway.TenantState,error){
  if !available{return schoolgateway.TenantState{},errors.New("registry unavailable")}
  return schoolgateway.TenantState{Generation:2,Active:active},nil
 })
 ctx:=context.Background()
 if _,err:=guarded.ExecuteSystemMail(ctx,"school",nil);!errors.Is(err,helpers.ErrConnectionInactive){t.Fatal("disabled school accepted",err)}
 if _,err:=guarded.Get(ctx,"school","connection");!errors.Is(err,helpers.ErrConnectionInactive){t.Fatal("disabled read accepted",err)}
 active=true;available=false
 if _,err:=guarded.ExecuteSystemMail(ctx,"school",nil);!errors.Is(err,helpers.ErrConnectionInactive){t.Fatal("unavailable registry accepted",err)}
 available=true
 stale:=schoolgateway.WithPrincipal(ctx,schoolgateway.Principal{AccountID:"account",SessionID:"session",TenantID:"school",SessionGeneration:1})
 if _,err:=guarded.ExecuteSystemMail(stale,"school",nil);!errors.Is(err,helpers.ErrConnectionInactive){t.Fatal("stale actor accepted",err)}
 if base.calls!=0{t.Fatal("denied calls reached service")}
 if _,err:=guarded.ExecuteSystemMail(ctx,"school",nil);err!=nil{t.Fatal(err)}
 if base.calls!=1{t.Fatal("active direct call not forwarded")}
}
