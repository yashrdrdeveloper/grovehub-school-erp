package impl

import (
 "net/http"
 httprequest "grovehub/libs/core/http/request"
 httpresponse "grovehub/libs/core/http/response"
 integrationerrors "grovehub/services/integration/constants/error"
 integrationrequest "grovehub/services/integration/dto/request"
 integrationresponse "grovehub/services/integration/dto/response"
)

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
