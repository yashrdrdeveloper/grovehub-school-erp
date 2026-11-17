package impl
import(
 "context"
 "net/http"
 "grovehub/libs/core/search"
 headercontext "grovehub/libs/core/headercontext"
 httprequest "grovehub/libs/core/http/request"
 rest "grovehub/libs/core/pipeline/rest"
 dto "grovehub/services/identity/internal/app/dto/account"
 definitions "grovehub/services/identity/internal/app/constants/error"
 helpers "grovehub/services/identity/internal/app/controller/helpers"
 protected "grovehub/services/identity/internal/app/controller/account/protected"
 service "grovehub/services/identity/internal/app/services/account"
)
type Controller struct{service service.Service}
func New(s service.Service)protected.AccountController{return &Controller{service:s}}
func(c *Controller)Create(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.CreateRequest,dto.CreateResult]{
  Decode:httprequest.DecodeStrictJSON,InvalidRequestBody:definitions.InvalidRequest,Fallback:definitions.Internal,
  Service:func(ctx context.Context,in dto.CreateRequest)(rest.Result[dto.CreateResult],error){
   out,err:=c.service.Create(ctx,in.RequestBody,headercontext.IdempotencyKey(ctx))
   return rest.Created(out),helpers.Map(err)
  },
 })
}
func(c *Controller)Search(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.SearchRequest,search.Result[dto.Account]]{
  Decode:httprequest.DecodeStrictJSON,InvalidRequestBody:definitions.InvalidRequest,Fallback:definitions.Internal,
  Service:func(ctx context.Context,in dto.SearchRequest)(rest.Result[search.Result[dto.Account]],error){
   out,err:=c.service.Search(ctx,in.RequestBody);return rest.OK(out),helpers.Map(err)
  },
 })
}
func(c *Controller)Get(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.GetRequest,dto.Account]{InvalidParams:definitions.InvalidRequest,Fallback:definitions.Internal,
  Service:func(ctx context.Context,in dto.GetRequest)(rest.Result[dto.Account],error){
   out,err:=c.service.Get(ctx,in.Params.AccountID);return rest.OK(out),helpers.Map(err)
  },
 })
}
func(c *Controller)Update(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.UpdateRequest,dto.Account]{
  Decode:httprequest.DecodeStrictJSON,InvalidParams:definitions.InvalidRequest,InvalidRequestBody:definitions.InvalidRequest,Fallback:definitions.Internal,
  Service:func(ctx context.Context,in dto.UpdateRequest)(rest.Result[dto.Account],error){
   out,err:=c.service.Update(ctx,in.Params.AccountID,in.RequestBody);return rest.OK(out),helpers.Map(err)
  },
 })
}
func(c *Controller)transition(w http.ResponseWriter,r *http.Request,verb string){
 rest.RunStruct(w,r,rest.StructConfig[dto.GetRequest,dto.Account]{InvalidParams:definitions.InvalidRequest,Fallback:definitions.Internal,
  Service:func(ctx context.Context,in dto.GetRequest)(rest.Result[dto.Account],error){
   out,err:=c.service.Transition(ctx,in.Params.AccountID,verb);return rest.OK(out),helpers.Map(err)
  },
 })
}
func(c *Controller)Suspend(w http.ResponseWriter,r *http.Request){c.transition(w,r,"suspend")}
func(c *Controller)Activate(w http.ResponseWriter,r *http.Request){c.transition(w,r,"activate")}
func(c *Controller)Deactivate(w http.ResponseWriter,r *http.Request){c.transition(w,r,"deactivate")}
func(c *Controller)Resend(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.GetRequest,map[string]any]{InvalidParams:definitions.InvalidRequest,Fallback:definitions.Internal,
  Service:func(ctx context.Context,in dto.GetRequest)(rest.Result[map[string]any],error){
   err:=c.service.Resend(ctx,in.Params.AccountID)
   return rest.Result[map[string]any]{StatusCode:http.StatusAccepted,Data:map[string]any{"accountId":in.Params.AccountID,"invitationStatus":"QUEUED"}},helpers.Map(err)
  },
 })
}
