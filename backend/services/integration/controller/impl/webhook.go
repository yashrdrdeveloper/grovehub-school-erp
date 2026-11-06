package impl

import (
 "net/http"
 integrationerrors "grovehub/services/integration/constants/error"
 "grovehub/services/integration/providers/gmail"
)

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
