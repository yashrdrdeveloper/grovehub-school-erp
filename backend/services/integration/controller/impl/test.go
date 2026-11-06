package impl

import (
 "net/http"
 httprequest "grovehub/libs/core/http/request"
 httpresponse "grovehub/libs/core/http/response"
 integrationerrors "grovehub/services/integration/constants/error"
)

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
