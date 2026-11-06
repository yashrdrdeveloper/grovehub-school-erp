package impl

import (
 "context"
 "encoding/json"
 "errors"
 "net/http"
 "strings"

 coreerrors "grovehub/libs/core/errors"
 httprequest "grovehub/libs/core/http/request"
 httpresponse "grovehub/libs/core/http/response"
 integrationconfig "grovehub/services/integration/config"
 integrationerrors "grovehub/services/integration/constants/error"
 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/providers"
 "grovehub/services/integration/providers/gmail"
 integrationrepository "grovehub/services/integration/repository"
 integrationservices "grovehub/services/integration/services"
 serviceimpl "grovehub/services/integration/services/impl"
)

type ConnectionController struct{
 service integrationservices.ConnectionService
 providerService integrationservices.ProviderService
 verifier *gmail.PushVerifier
 pushEnabled bool
}

type VerifiedPrincipal struct{TenantID,ActorGID string}
type principalContextKey struct{}
func NewConnectionController(s integrationservices.ConnectionService,p integrationservices.ProviderService,push integrationconfig.GmailPush)*ConnectionController{
 return &ConnectionController{service:s,providerService:p,verifier:gmail.NewPushVerifier(push.Audience,push.ServiceAccount),pushEnabled:push.Enabled()}
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

func (c *ConnectionController) Create(w http.ResponseWriter,r *http.Request){
 p,ok:=principalFrom(r.Context());if !ok {writeError(w,r,http.StatusForbidden,integrationerrors.Forbidden);return}
 var req integrationrequest.CreateConnection
 if err:=httprequest.DecodeStrictJSON(r,&req);err!=nil{writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid);return}
 if req.OwnerGID!=p.ActorGID {writeError(w,r,http.StatusForbidden,integrationerrors.Forbidden);return}
 result,err:=c.service.Create(r.Context(),p.TenantID,req)
 if err!=nil {c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusCreated,result)
}

func (c *ConnectionController) Get(w http.ResponseWriter,r *http.Request){
 p,ok:=principalFrom(r.Context());if !ok {writeError(w,r,http.StatusForbidden,integrationerrors.Forbidden);return}
 result,err:=c.service.Get(r.Context(),p.TenantID,r.PathValue("id"))
 if err!=nil{c.writeServiceError(w,r,err);return}
 if result.OwnerGID!=p.ActorGID{writeError(w,r,http.StatusForbidden,integrationerrors.Forbidden);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}

func (c *ConnectionController) List(w http.ResponseWriter,r *http.Request){
 p,ok:=principalFrom(r.Context());if !ok {writeError(w,r,http.StatusForbidden,integrationerrors.Forbidden);return}
 rows,err:=c.service.List(r.Context(),p.TenantID,p.ActorGID)
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,rows)
}

func (c *ConnectionController) Disconnect(w http.ResponseWriter,r *http.Request){
 p,ok:=principalFrom(r.Context());if !ok {writeError(w,r,http.StatusForbidden,integrationerrors.Forbidden);return}
 id:=r.PathValue("id")
 current,err:=c.service.Get(r.Context(),p.TenantID,id)
 if err!=nil{c.writeServiceError(w,r,err);return}
 if current.OwnerGID!=p.ActorGID{writeError(w,r,http.StatusForbidden,integrationerrors.Forbidden);return}
 result,err:=c.service.Disconnect(r.Context(),p.TenantID,id)
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}


func(c *ConnectionController)BeginOAuth(w http.ResponseWriter,r *http.Request){
 p,ok:=c.owned(w,r);if !ok{return}
 result,err:=c.service.BeginOAuth(r.Context(),p.TenantID,r.PathValue("id"))
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}

func(c *ConnectionController)FinishOAuth(w http.ResponseWriter,r *http.Request){
 if r.URL.Query().Get("error")!=""{writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid);return}
 result,err:=c.service.FinishOAuth(r.Context(),strings.ToUpper(r.PathValue("provider")),r.URL.Query().Get("code"),r.URL.Query().Get("state"))
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}


func(c *ConnectionController)RunTest(w http.ResponseWriter,r *http.Request){
 p,ok:=c.owned(w,r);if !ok{return}
 var payload struct{Direction string `json:"direction"`}
 if err:=httprequest.DecodeStrictJSON(r,&payload);err!=nil{writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid);return}
 result,err:=c.service.RunTest(r.Context(),p.TenantID,r.PathValue("id"),payload.Direction)
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}

func(c *ConnectionController)ListTests(w http.ResponseWriter,r *http.Request){
 p,ok:=c.owned(w,r);if !ok{return}
 result,err:=c.service.ListTests(r.Context(),p.TenantID,r.PathValue("id"))
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}


func(c *ConnectionController)Watch(w http.ResponseWriter,r *http.Request){
 p,ok:=c.owned(w,r);if !ok{return}
 result,err:=c.service.Watch(r.Context(),p.TenantID,r.PathValue("id"))
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}


func(c *ConnectionController)GmailPush(w http.ResponseWriter,r *http.Request){
 if !c.pushEnabled{
  writeError(w,r,http.StatusServiceUnavailable,integrationerrors.NotReady);return
 }
 if err:=c.verifier.Verify(r.Context(),r.Header.Get("Authorization"));err!=nil{
  writeError(w,r,http.StatusUnauthorized,integrationerrors.Forbidden);return
 }
 data,err:=gmail.DecodePush(r.Body)
 if err!=nil{writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid);return}
 if err=c.service.ProcessPush(r.Context(),data);err!=nil{
  // Pub/Sub retries on publish failure; Communication will ingest later.
  writeError(w,r,http.StatusServiceUnavailable,integrationerrors.Failed);return
 }
 w.WriteHeader(http.StatusNoContent)
}


func(c *ConnectionController)ExecuteAction(w http.ResponseWriter,r *http.Request){
 p,ok:=c.owned(w,r)
 if !ok{return}
 action:=providers.Action(strings.TrimSpace(r.PathValue("action")))
 if action==""||len(action)>100 {
  writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid)
  return
 }
 var payload struct{Input json.RawMessage `json:"input"`}
 if err:=httprequest.DecodeStrictJSON(r,&payload);err!=nil{
  writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid)
  return
 }
 result,err:=c.service.ExecuteAction(r.Context(),p.TenantID,r.PathValue("id"),action,payload.Input)
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}
