package impl

import (
 "net/http"

 httpresponse "grovehub/libs/core/http/response"
)

func (c *ConnectionController) Providers(w http.ResponseWriter,r *http.Request){
 _=httpresponse.WriteSuccess(w,http.StatusOK,c.providerService.Providers(r.Context()))
}

func(c *ConnectionController)Provider(w http.ResponseWriter,r *http.Request){
 result,err:=c.providerService.Provider(r.Context(),r.PathValue("provider"))
 if err!=nil{c.writeServiceError(w,r,err);return}
 _=httpresponse.WriteSuccess(w,http.StatusOK,result)
}
