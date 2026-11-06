package impl

import (
 "context"
 "errors"
 "net/http"

 httprequest "grovehub/libs/core/http/request"
 httpresponse "grovehub/libs/core/http/response"
 integrationerrors "grovehub/services/integration/constants/error"
 coreerrors "grovehub/libs/core/errors"
 "grovehub/services/integration/providers"
 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
 integrationrepository "grovehub/services/integration/repository"
 integrationservices "grovehub/services/integration/services"
 serviceimpl "grovehub/services/integration/services/impl"
)

// VerifiedPrincipal must only be populated by a trusted authorization middleware.
// Callers cannot grant access using arbitrary headers or JSON tenant IDs.
type VerifiedPrincipal struct{TenantID,ActorGID string}
type principalContextKey struct{}
func WithVerifiedPrincipal(ctx context.Context,principal VerifiedPrincipal) context.Context {
 return context.WithValue(ctx,principalContextKey{},principal)
}
func principalFrom(ctx context.Context)(VerifiedPrincipal,bool){
 p,ok:=ctx.Value(principalContextKey{}).(VerifiedPrincipal)
 return p,ok && p.TenantID!="" && p.ActorGID!=""
}

type ConnectionController struct{service integrationservices.ConnectionService}
func NewConnectionController(s integrationservices.ConnectionService)*ConnectionController{return &ConnectionController{service:s}}
func (c *ConnectionController) Providers(w http.ResponseWriter,r *http.Request){
 _=httpresponse.WriteSuccess(w,http.StatusOK,c.service.Providers(r.Context()))
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
 rows,err:=c.service.List(r.Context(),p.TenantID)
 if err!=nil{c.writeServiceError(w,r,err);return}
 own:=make([]integrationresponse.Connection,0,len(rows))
 for _,row:=range rows{if row.OwnerGID==p.ActorGID{own=append(own,row)}}
 _=httpresponse.WriteSuccess(w,http.StatusOK,own)
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
func (c *ConnectionController) writeServiceError(w http.ResponseWriter,r *http.Request,err error){
 switch {
 case errors.Is(err,serviceimpl.ErrInvalidConnection),errors.Is(err,integrationrepository.ErrConnectionNotFound):
  if errors.Is(err,integrationrepository.ErrConnectionNotFound){writeError(w,r,http.StatusNotFound,integrationerrors.NotFound)}else{writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid)}
 case errors.Is(err,providers.ErrUnknownProvider):
  writeError(w,r,http.StatusBadRequest,integrationerrors.Invalid)
 case errors.Is(err,serviceimpl.ErrConnectionNotReady):
  writeError(w,r,http.StatusNotImplemented,integrationerrors.NotReady)
 default:
  writeError(w,r,http.StatusInternalServerError,integrationerrors.Failed)
 }
}
func writeError(w http.ResponseWriter,r *http.Request,status int,definition coreerrors.Definition){
 _=httpresponse.WriteErrorContext(r.Context(),w,status,definition)
}
