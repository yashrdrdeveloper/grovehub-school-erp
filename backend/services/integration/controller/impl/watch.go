package impl

import (
 "net/http"
 httpresponse "grovehub/libs/core/http/response"
)

func(c *ConnectionController)Watch(w http.ResponseWriter,r *http.Request){
 p,ok:=c.owned(w,r);if !ok{return}
 result,err:=c.service.Watch(r.Context(),p.TenantID,r.PathValue("id"))
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}
