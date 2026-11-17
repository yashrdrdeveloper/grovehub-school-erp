package impl
import(
 "context"
 "net/http"
 httprequest "grovehub/libs/core/http/request"
 rest "grovehub/libs/core/pipeline/rest"
 dto "grovehub/services/identity/internal/app/dto/account"
 definitions "grovehub/services/identity/internal/app/constants/error"
 helpers "grovehub/services/identity/internal/app/controller/helpers"
 "grovehub/services/identity/internal/app/middleware"
 public "grovehub/services/identity/internal/app/controller/account/public"
 service "grovehub/services/identity/internal/app/services/account"
)
type Controller struct{service service.Service}
func New(s service.Service)public.Controller{return &Controller{service:s}}
func(c *Controller)Setup(w http.ResponseWriter,r *http.Request){
 rest.RunStruct(w,r,rest.StructConfig[dto.SetupRequest,dto.SetupResult]{
 Decode:httprequest.DecodeStrictJSON,InvalidRequestBody:definitions.InvalidRequest,Fallback:definitions.Internal,
 Service:func(ctx context.Context,in dto.SetupRequest)(rest.Result[dto.SetupResult],error){
  out,err:=c.service.Setup(ctx,middleware.ReadCookie(r,middleware.SetupCookie),in.RequestBody)
  if err!=nil{return rest.Result[dto.SetupResult]{},helpers.Map(err)}
  middleware.ClearGrant(w,middleware.SetupCookie,"/identity/public/accounts/setup")
  return rest.OK(out),nil
 },
 })
}
