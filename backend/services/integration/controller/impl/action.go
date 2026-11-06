package impl

import (
 "net/http"
 "encoding/json"
 "strings"
 httprequest "grovehub/libs/core/http/request"
 httpresponse "grovehub/libs/core/http/response"
 integrationerrors "grovehub/services/integration/constants/error"
 "grovehub/services/integration/providers"
)

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
