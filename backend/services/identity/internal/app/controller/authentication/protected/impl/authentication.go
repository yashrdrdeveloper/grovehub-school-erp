package impl
import(
 "context"
 "net/http"
 request "grovehub/libs/core/http/request"
 rest "grovehub/libs/core/pipeline/rest"
 definitions "grovehub/services/identity/internal/app/constants/error"
 dto "grovehub/services/identity/internal/app/dto/authentication"
 accountdto "grovehub/services/identity/internal/app/dto/account"
 helpers "grovehub/services/identity/internal/app/controller/helpers"
 "grovehub/services/identity/internal/app/middleware"
 protected "grovehub/services/identity/internal/app/controller/authentication/protected"
 auth "grovehub/services/identity/internal/app/services/authentication"
)
type Controller struct{service auth.Service}
func New(s auth.Service)protected.AuthenticationController{return &Controller{service:s}}
func(c *Controller)Logout(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.EmptyRequest,dto.StatusResult]{Fallback:definitions.Internal,
 Service:func(ctx context.Context,_ dto.EmptyRequest)(rest.Result[dto.StatusResult],error){
  if err:=c.service.Logout(ctx);err!=nil{return rest.Result[dto.StatusResult]{},helpers.Map(err)}
  middleware.ClearSessionCookies(w);return rest.OK(dto.StatusResult{Status:"LOGGED_OUT"}),nil
 }})
}
func(c *Controller)LogoutAll(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.EmptyRequest,map[string]any]{Fallback:definitions.Internal,
 Service:func(ctx context.Context,_ dto.EmptyRequest)(rest.Result[map[string]any],error){
  count,err:=c.service.LogoutAll(ctx);if err!=nil{return rest.Result[map[string]any]{},helpers.Map(err)}
  middleware.ClearSessionCookies(w);return rest.OK(map[string]any{"status":"ALL_SESSIONS_REVOKED","revokedCount":count}),nil
 }})
}
func(c *Controller)ChangePassword(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.ChangePasswordRequest,map[string]any]{Decode:request.DecodeStrictJSON,InvalidRequestBody:definitions.InvalidRequest,Fallback:definitions.Internal,
 Service:func(ctx context.Context,in dto.ChangePasswordRequest)(rest.Result[map[string]any],error){
  if err:=c.service.ChangePassword(ctx,in.RequestBody.NewPassword);err!=nil{return rest.Result[map[string]any]{},helpers.Map(err)}
  middleware.ClearSessionCookies(w);return rest.OK(map[string]any{"status":"PASSWORD_CHANGED","sessionsRevoked":true}),nil
 }})
}
func(c *Controller)Reauthenticate(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.ReauthenticateRequest,dto.StatusResult]{Decode:request.DecodeStrictJSON,InvalidRequestBody:definitions.InvalidRequest,Fallback:definitions.Internal,
 Service:func(ctx context.Context,in dto.ReauthenticateRequest)(rest.Result[dto.StatusResult],error){
  if err:=c.service.Reauthenticate(ctx,in.RequestBody.Password);err!=nil{return rest.Result[dto.StatusResult]{},helpers.Map(err)}
  return rest.OK(dto.StatusResult{Status:"REAUTHENTICATED"}),nil
 }})
}
func(c *Controller)Sessions(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.EmptyRequest,map[string]any]{Fallback:definitions.Internal,
 Service:func(ctx context.Context,_ dto.EmptyRequest)(rest.Result[map[string]any],error){
  items,err:=c.service.ListSessions(ctx);if err!=nil{return rest.Result[map[string]any]{},helpers.Map(err)}
  return rest.OK(map[string]any{"items":items}),nil
 }})
}
func(c *Controller)RevokeSession(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.EmptyRequest,map[string]any]{Fallback:definitions.Internal,
 Service:func(ctx context.Context,_ dto.EmptyRequest)(rest.Result[map[string]any],error){
  id:=r.PathValue("sessionId");if err:=c.service.RevokeSession(ctx,id);err!=nil{return rest.Result[map[string]any]{},helpers.Map(err)}
  return rest.OK(map[string]any{"sessionId":id,"status":"REVOKED"}),nil
 }})
}
func(c *Controller)Me(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.EmptyRequest,accountdto.Account]{Fallback:definitions.Internal,
 Service:func(ctx context.Context,_ dto.EmptyRequest)(rest.Result[accountdto.Account],error){
  v,err:=c.service.Me(ctx);return rest.OK(v),helpers.Map(err)
 }})
}
