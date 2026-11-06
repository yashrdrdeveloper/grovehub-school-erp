package impl

import (
 "context"
 "errors"
 "net/http"
 coreerrors "grovehub/libs/core/errors"
 httpresponse "grovehub/libs/core/http/response"
 integrationerrors "grovehub/services/integration/constants/error"
 integrationrepository "grovehub/services/integration/repository"
 integrationservices "grovehub/services/integration/services"
 serviceimpl "grovehub/services/integration/services/impl"
 "grovehub/services/integration/providers"
 "grovehub/services/integration/providers/gmail"
 integrationconfig "grovehub/services/integration/config"
)


type VerifiedPrincipal struct{TenantID,ActorGID string}
type principalContextKey struct{}
type ConnectionController struct{
 service integrationservices.ConnectionService
 verifier *gmail.PushVerifier
 pushEnabled bool
}
func NewConnectionController(s integrationservices.ConnectionService,push integrationconfig.GmailPush)*ConnectionController{
 return &ConnectionController{service:s,verifier:gmail.NewPushVerifier(push.Audience,push.ServiceAccount),pushEnabled:push.Enabled()}
}

func WithVerifiedPrincipal(ctx context.Context,principal VerifiedPrincipal) context.Context {
 return context.WithValue(ctx,principalContextKey{},principal)
}

func principalFrom(ctx context.Context)(VerifiedPrincipal,bool){
 p,ok:=ctx.Value(principalContextKey{}).(VerifiedPrincipal)
 return p,ok && p.TenantID!="" && p.ActorGID!=""
}


func (c *ConnectionController) writeServiceError(w http.ResponseWriter,r *http.Request,err error){
 if status:=providers.HTTPStatus(err);status!=0{
  switch status{
  case http.StatusBadRequest:writeError(w,r,status,integrationerrors.Invalid)
  case http.StatusUnauthorized,http.StatusForbidden:writeError(w,r,status,integrationerrors.Forbidden)
  case http.StatusNotFound:writeError(w,r,status,integrationerrors.NotFound)
  case http.StatusConflict:writeError(w,r,status,integrationerrors.NotReady)
  default:writeError(w,r,status,integrationerrors.Failed)
  }
  return
 }
 switch {
 case errors.Is(err,serviceimpl.ErrInvalidConnection),errors.Is(err,integrationrepository.ErrConnectionNotFound):
  if errors.Is(err,integrationrepository.ErrConnectionNotFound){writeError(w,r,http.StatusNotFound,integrationerrors.NotFound)}else{writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid)}
 case errors.Is(err,providers.ErrInvalidActionInput):
  writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid)
 case errors.Is(err,providers.ErrUnsupportedAction):
  writeError(w,r,http.StatusNotImplemented,integrationerrors.NotReady)
 case errors.Is(err,providers.ErrProviderUnavailable):
  writeError(w,r,http.StatusServiceUnavailable,integrationerrors.NotReady)
 case errors.Is(err,providers.ErrUnknownProvider):
  writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid)
 case errors.Is(err,serviceimpl.ErrOAuthInvalid):
  writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid)
 case errors.Is(err,serviceimpl.ErrPushUnavailable):
 writeError(w,r,http.StatusServiceUnavailable,integrationerrors.NotReady)
 case errors.Is(err,serviceimpl.ErrConnectionInactive):
  writeError(w,r,http.StatusConflict,integrationerrors.NotReady)
 case errors.Is(err,serviceimpl.ErrOAuthUnavailable),errors.Is(err,serviceimpl.ErrConnectionNotReady):
  writeError(w,r,http.StatusNotImplemented,integrationerrors.NotReady)
 default:
  writeError(w,r,http.StatusInternalServerError,integrationerrors.Failed)
 }
}

func writeError(w http.ResponseWriter,r *http.Request,status int,definition coreerrors.Definition){
 _=httpresponse.WriteErrorContext(r.Context(),w,status,definition)
}

func(c *ConnectionController)owned(w http.ResponseWriter,r *http.Request)(VerifiedPrincipal,bool){
 p,ok:=principalFrom(r.Context())
 if !ok{writeError(w,r,http.StatusForbidden,integrationerrors.Forbidden);return p,false}
 current,err:=c.service.Get(r.Context(),p.TenantID,r.PathValue("id"))
 if err!=nil{c.writeServiceError(w,r,err);return p,false}
 if current.OwnerGID!=p.ActorGID{writeError(w,r,http.StatusForbidden,integrationerrors.Forbidden);return p,false}
 return p,true
}
