package impl

import (
 "context"
 "encoding/json"

 schoolgateway "grovehub/launcher/session"
 dto "grovehub/services/integration/internal/app/dto/connection"
 "grovehub/services/integration/internal/app/externalService/mail/gmail"
 providerdto "grovehub/services/integration/internal/app/providers/dto"
 service "grovehub/services/integration/internal/app/services/connection"
 helpers "grovehub/services/integration/internal/app/services/connection/helpers"
)

// tenantService is the production boundary for direct clients, controllers,
// callbacks and scheduled workers. Every operation is explicitly forwarded.
type tenantService struct{
 base service.ConnectionService
 lookup schoolgateway.TenantStateLookup
}
var _ service.ConnectionService=(*tenantService)(nil)

func WithTenantState(base service.ConnectionService, lookup schoolgateway.TenantStateLookup) service.ConnectionService{
 if lookup==nil{return base}
 return &tenantService{base:base,lookup:lookup}
}
func(s *tenantService) context(ctx context.Context)context.Context{return helpers.WithTenantState(ctx,s.lookup)}
func(s *tenantService) scope(ctx context.Context,tenant string)(context.Context,error){
 ctx=s.context(ctx)
 _,err:=helpers.RequireTenantGeneration(ctx,tenant)
 return ctx,err
}
func(s *tenantService) Create(ctx context.Context,tenant string,input dto.CreateConnection)(dto.Connection,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return dto.Connection{},err};return s.base.Create(ctx,tenant,input)
}
func(s *tenantService) Get(ctx context.Context,tenant,id string)(dto.Connection,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return dto.Connection{},err};return s.base.Get(ctx,tenant,id)
}
func(s *tenantService) List(ctx context.Context,tenant,owner string)([]dto.Connection,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return nil,err};return s.base.List(ctx,tenant,owner)
}
func(s *tenantService) Disconnect(ctx context.Context,tenant,id string)(dto.Connection,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return dto.Connection{},err};return s.base.Disconnect(ctx,tenant,id)
}
func(s *tenantService) BeginOAuth(ctx context.Context,tenant,id string)(dto.OAuthStart,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return dto.OAuthStart{},err};return s.base.BeginOAuth(ctx,tenant,id)
}
func(s *tenantService) FinishOAuth(ctx context.Context,provider,code,state string)(dto.OAuthCompletion,error){
 return s.base.FinishOAuth(s.context(ctx),provider,code,state)
}
func(s *tenantService) RunTest(ctx context.Context,tenant,id,direction string)(dto.ConnectionTest,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return dto.ConnectionTest{},err};return s.base.RunTest(ctx,tenant,id,direction)
}
func(s *tenantService) ListTests(ctx context.Context,tenant,id string)([]dto.ConnectionTest,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return nil,err};return s.base.ListTests(ctx,tenant,id)
}
func(s *tenantService) Watch(ctx context.Context,tenant,id string)(dto.Watch,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return dto.Watch{},err};return s.base.Watch(ctx,tenant,id)
}
func(s *tenantService) ProcessPush(ctx context.Context,n gmail.PushNotification)error{
 return s.base.ProcessPush(s.context(ctx),n)
}
func(s *tenantService) RenewExpiring(ctx context.Context)error{
 return s.base.RenewExpiring(s.context(ctx))
}
func(s *tenantService) ExecuteAction(ctx context.Context,tenant,id string,action providerdto.Action,input json.RawMessage)(any,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return nil,err};return s.base.ExecuteAction(ctx,tenant,id,action,input)
}
func(s *tenantService) ExecuteSystemMail(ctx context.Context,tenant string,input json.RawMessage)(any,error){
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return nil,err};return s.base.ExecuteSystemMail(ctx,tenant,input)
}
func(s *tenantService) ConfigureSystemMailSender(ctx context.Context,tenant,id string)error{
 ctx,err:=s.scope(ctx,tenant);if err!=nil{return err};return s.base.ConfigureSystemMailSender(ctx,tenant,id)
}
