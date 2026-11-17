package impl
import(
 "context"
 "net/http"
 request "grovehub/libs/core/http/request"
 rest "grovehub/libs/core/pipeline/rest"
 dto "grovehub/services/identity/internal/app/dto/authentication"
 definitions "grovehub/services/identity/internal/app/constants/error"
 helpers "grovehub/services/identity/internal/app/controller/helpers"
 "grovehub/services/identity/internal/app/middleware"
 public "grovehub/services/identity/internal/app/controller/authentication/public"
 auth "grovehub/services/identity/internal/app/services/authentication"
 secrets "grovehub/services/identity/internal/app/services/authentication/helpers"
)
type Controller struct{service auth.Service}
func New(s auth.Service)public.Controller{return &Controller{service:s}}
func(c *Controller)Login(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.LoginRequest,dto.LoginResult]{
 Decode:request.DecodeStrictJSON,InvalidRequestBody:definitions.InvalidRequest,Fallback:definitions.Internal,
 Service:func(ctx context.Context,in dto.LoginRequest)(rest.Result[dto.LoginResult],error){
  result,err:=c.service.Login(ctx,in.RequestBody,r.UserAgent());if err!=nil{return rest.Result[dto.LoginResult]{},helpers.Map(err)}
  csrf,err:=secrets.NewSecret(32);if err!=nil{return rest.Result[dto.LoginResult]{},helpers.Map(err)}
  middleware.IssueCookies(w,result.Cookies.Access,result.Cookies.Refresh,csrf)
  return rest.OK(result.Data),nil
 },
 })
}
func(c *Controller)Refresh(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.EmptyRequest,dto.RefreshResult]{
 Fallback:definitions.Internal,
 Service:func(ctx context.Context,_ dto.EmptyRequest)(rest.Result[dto.RefreshResult],error){
  result,err:=c.service.Refresh(ctx,middleware.ReadCookie(r,middleware.RefreshCookie));if err!=nil{return rest.Result[dto.RefreshResult]{},helpers.Map(err)}
  middleware.IssueCookies(w,result.Cookies.Access,result.Cookies.Refresh,"")
  return rest.OK(result.Data),nil
 },
 })
}
func(c *Controller)ForgotPassword(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.ForgotPasswordRequest,dto.AcceptedResult]{
 Decode:request.DecodeStrictJSON,InvalidRequestBody:definitions.InvalidRequest,Fallback:definitions.Internal,
 Service:func(ctx context.Context,in dto.ForgotPasswordRequest)(rest.Result[dto.AcceptedResult],error){
  // Always return the same success response; do not disclose eligibility.
  _=c.service.Forgot(ctx,in.RequestBody.Email)
  return rest.Result[dto.AcceptedResult]{StatusCode:http.StatusAccepted,Data:dto.AcceptedResult{Accepted:true,Delivery:"EMAIL"}},nil
 },
 })
}
func(c *Controller)exchange(w http.ResponseWriter,r *http.Request,prefix,purpose,cookie,path string){
 rest.RunStruct(w,r,rest.StructConfig[dto.EmptyRequest,dto.GrantResult]{
 Fallback:definitions.Internal,
 Service:func(ctx context.Context,_ dto.EmptyRequest)(rest.Result[dto.GrantResult],error){
  grant,err:=c.service.Exchange(ctx,middleware.Proof(r,prefix),purpose);if err!=nil{return rest.Result[dto.GrantResult]{},helpers.Map(err)}
  middleware.SetGrant(w,cookie,path,grant.Secret)
  return rest.OK(grant.Data),nil
 },
 })
}
func(c *Controller)RecoveryExchange(w http.ResponseWriter,r *http.Request){c.exchange(w,r,"Recovery","recovery",middleware.ResetCookie,"/identity/public/auth/reset-password")}
func(c *Controller)InvitationExchange(w http.ResponseWriter,r *http.Request){c.exchange(w,r,"Invitation","invitation",middleware.SetupCookie,"/identity/public/accounts/setup")}
func(c *Controller)ResetPassword(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.ResetPasswordRequest,dto.StatusResult]{
 Decode:request.DecodeStrictJSON,InvalidRequestBody:definitions.InvalidRequest,Fallback:definitions.Internal,
 Service:func(ctx context.Context,in dto.ResetPasswordRequest)(rest.Result[dto.StatusResult],error){
  err:=c.service.ResetPassword(ctx,middleware.ReadCookie(r,middleware.ResetCookie),in.RequestBody.NewPassword)
  if err!=nil{return rest.Result[dto.StatusResult]{},helpers.Map(err)}
  middleware.ClearGrant(w,middleware.ResetCookie,"/identity/public/auth/reset-password")
  return rest.OK(dto.StatusResult{Status:"PASSWORD_RESET"}),nil
 },
 })
}
