package impl

import (
 "net/http"
 "strings"
 httpresponse "grovehub/libs/core/http/response"
 integrationerrors "grovehub/services/integration/constants/error"
)

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
